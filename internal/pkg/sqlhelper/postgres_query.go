package sqlhelper

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"reflect"
	"strings"

	"github.com/example/internal/i18n"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresModel interface {
	TableName() string
}

func Query[T PostgresModel](DB *pgxpool.Pool, sql string, args ...any) (*[]T, error) {
	row, err := DB.Query(context.Background(), sql, args...)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	data := []T{}

	for i := 0; row.Next(); i++ {
		ptrs := []any{}
		var d T
		data = append(data, d)
		fieldMap, _, _ := fieldMap(&data[len(data)-1], false)
		fields := row.FieldDescriptions()
		for _, fd := range fields {
			ptr := fieldMap[fd.Name]
			if ptr != nil {
				ptrs = append(ptrs, ptr)
			}
		}

		err := row.Scan(ptrs...)
		if err != nil {
			return nil, err
		}

	}

	return &data, nil
}

func Create[T PostgresModel](DB *pgxpool.Pool, models ...*T) (*[]T, error) {
	data := []T{}
	if len(models) == 0 {
		return &data, nil
	}

	queryArgs := []any{}
	tableName := (*models[0]).TableName()
	query := fmt.Sprintf("INSERT INTO %s (\n", tableName)
	_, columns, _ := fieldMap(models[0], true)
	query += strings.Join(columns, ", ")
	query += ") VALUES "

	valuesToUpdate := []string{}
	for _, model := range models {
		_, _, values := fieldMap(model, true)
		value := []string{}
		for _, v := range values {
			queryArgs = append(queryArgs, v)
			value = append(value, fmt.Sprintf("$%d", len(queryArgs)))
		}
		valuesToUpdate = append(valuesToUpdate, fmt.Sprintf("(%s)", strings.Join(value, ", ")))
	}

	query += strings.Join(valuesToUpdate, ",\n")
	query += " RETURNING  *;"

	row, err := DB.Query(context.Background(), query, queryArgs...)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	for i := 0; row.Next(); i++ {
		ptrs := []any{}
		var d T
		data = append(data, d)
		fieldMap, _, _ := fieldMap(&data[len(data)-1], false)
		fields := row.FieldDescriptions()
		for _, fd := range fields {
			ptr := fieldMap[fd.Name]
			if ptr != nil {
				ptrs = append(ptrs, ptr)
			}
		}

		err := row.Scan(ptrs...)
		if err != nil {
			return nil, err
		}

	}

	return &data, nil
}

func Update[T PostgresModel](DB *pgxpool.Pool, model *T, Where string, args ...any) error {
	queryArgs := []any{}
	fieldMap, columns, values := fieldMap(model, true)
	tableName := (*model).TableName()

	query := fmt.Sprintf("UPDATE %s SET \n", tableName)

	updateStr := []string{}
	for i, columnName := range columns {
		queryArgs = append(queryArgs, values[i])
		updateStr = append(updateStr, fmt.Sprintf("%s = $%d", columnName, len(queryArgs)))
	}
	query += strings.Join(updateStr, ",\n")

	Where = processWhereQuery(Where, len(queryArgs))
	queryArgs = append(queryArgs, args...)
	query += fmt.Sprintf(" WHERE %s \n", Where)

	query += " RETURNING  *;"

	row, err := DB.Query(context.Background(), query, queryArgs...)
	if err != nil {
		return err
	}
	defer row.Close()

	if !row.Next() {
		return errors.New(i18n.KEY_INTERNAL_ERROR_UPDATE_FAILED)
	}

	ptrs := []any{}
	for _, fd := range row.FieldDescriptions() {
		ptr := fieldMap[fd.Name]
		if ptr != nil {
			ptrs = append(ptrs, ptr)
		}
	}

	err = row.Scan(ptrs...)
	if err != nil {
		return err
	}

	return nil
}

func RemoveField[T PostgresModel](DB *pgxpool.Pool, columns *[]string, Where string, args ...any) error {
	if len(*columns) == 0 {
		return nil
	}
	var model T
	tableName := model.TableName()
	query := fmt.Sprintf("UPDATE %s SET \n", tableName)

	updateStr := []string{}
	for _, columnName := range *columns {
		updateStr = append(updateStr, fmt.Sprintf("%s = DEFAULT ", columnName))
	}
	query += strings.Join(updateStr, ",\n")

	Where = processWhereQuery(Where, 0)
	query += fmt.Sprintf(" WHERE %s ;", Where)

	row, err := DB.Query(context.Background(), query)
	if err != nil {
		return err
	}
	defer row.Close()

	return nil
}

func Delete[T PostgresModel](DB *pgxpool.Pool, Where string, args ...any) error {
	queryArgs := []any{}
	var model T
	tableName := model.TableName()
	Where = processWhereQuery(Where, len(queryArgs))
	queryArgs = append(queryArgs, args...)
	query := fmt.Sprintf("DELETE FROM %s WHERE %s", tableName, Where)

	_, err := DB.Exec(context.Background(), query, queryArgs...)
	if err != nil {
		return err
	}

	return nil
}

func Check(DB *pgxpool.Pool, scan any, query string, args ...any) error {
	rows, err := DB.Query(context.Background(), query, args...)
	if err != nil {
		return err
	}

	if rows.Next() {
		if err := rows.Scan(scan); err != nil {
			return err
		}
	}

	return nil
}

func fieldMap(value any, hideNil bool) (fieldMapData map[string]any, columns []string, values []any) {
	val := reflect.ValueOf(value)
	typ := reflect.TypeOf(value)
	fieldMapData = make(map[string]any)
	for i := 0; i < val.Elem().NumField(); i++ {
		valElem := val.Elem().Field(i)
		typElem := typ.Elem().Field(i)
		tag := typElem.Tag
		fieldName := tag.Get("column")
		isEmbedded := tag.Get("data") == "embedded"

		if fieldName != "" {
			fieldMapData[fieldName] = valElem.Addr().Interface()

			if !hideNil || valElem.Kind() != reflect.Ptr || !valElem.IsNil() {
				columns = append(columns, fieldName)
				values = append(values, valElem.Interface())
			}
		} else if isEmbedded {
			subVal := valElem
			subTyp := typElem.Type
			subMap, subColumn, subValue := processEmbeddedElem(subVal, subTyp, hideNil)
			columns = append(columns, subColumn...)
			values = append(values, subValue...)
			maps.Copy(fieldMapData, subMap)
		}
	}
	return fieldMapData, columns, values
}

func processEmbeddedElem(subVal reflect.Value, subTyp reflect.Type, hideNil bool) (fieldMapData map[string]any, columns []string, values []any) {
	fieldMapData = make(map[string]any)

	for i := 0; i < subVal.NumField(); i++ {
		val := subVal.Field(i)
		typ := subTyp.Field(i)
		tag := typ.Tag
		fieldName := tag.Get("column")
		isEmbedded := tag.Get("data") == "embedded"

		if fieldName != "" {
			fieldMapData[fieldName] = val.Addr().Interface()

			if !hideNil || val.Kind() != reflect.Ptr || !val.IsNil() {
				columns = append(columns, fieldName)
				values = append(values, val.Interface())
			}
		} else if isEmbedded {
			subVal := val
			subTyp := typ.Type
			subMap, subColumn, subValue := processEmbeddedElem(subVal, subTyp, hideNil)
			columns = append(columns, subColumn...)
			values = append(values, subValue...)
			maps.Copy(fieldMapData, subMap)
		}
	}
	return fieldMapData, columns, values
}

func processWhereQuery(s string, offset int) string {
	var sb strings.Builder
	counter := offset
	for _, ch := range s {
		if ch == '?' {
			counter++
			sb.WriteString(fmt.Sprintf("$%d", counter))
		} else {
			sb.WriteRune(ch)
		}
	}

	return sb.String()
}
