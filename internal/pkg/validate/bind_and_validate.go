package validate

import (
	"errors"
	"net/http"
	"reflect"

	"github.com/example/internal/i18n"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func Bind(c echo.Context, data interface{}) error {
	if err := c.Bind(data); err != nil {
		return err
	}
	if err := c.Validate(data); err != nil {
		return err
	}
	return nil
}

func BindFormData(c echo.Context, data interface{}) error {
	val := reflect.ValueOf(data)
	typ := reflect.TypeOf(data)

	multipartForm, err := c.MultipartForm()
	if err != nil {
		return err
	}

	for i := 0; i < val.Elem().NumField(); i++ {
		typElem := typ.Elem().Field(i)
		tag := typElem.Tag
		required := tag.Get("validate") == "required"
		fieldName := tag.Get("name")
		isFile := tag.Get("type") == "file"
		isFiles := tag.Get("type") == "files"

		dataSetter := val.Elem().Field(i).Set
		requiredErr := errors.New(i18n.Trsl(c, "{field_name} must not be empty", fieldName))

		if isFile {
			file, err := c.FormFile(fieldName)
			if err != nil {
				if required && err == http.ErrMissingFile {
					return requiredErr
				} else if err != http.ErrMissingFile {
					return err
				}
			}
			if file != nil {
				dataSetter(reflect.ValueOf(&file))
			}
		} else if isFiles {
			files := multipartForm.File[fieldName]
			if required && len(files) == 0 {
				return requiredErr
			}
			dataSetter(reflect.ValueOf(&files))
		} else {
			value := c.FormValue(fieldName)
			if required && value == "" {
				return requiredErr
			}
			if value != "" {
				dataSetter(reflect.ValueOf(&value))
			}
		}
	}

	return nil
}

// To initiate custom validation in main.go
type CustomValidator struct {
	Validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.Validator.Struct(i)
}
