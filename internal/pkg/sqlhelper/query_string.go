package sqlhelper

import (
	"fmt"
	"strings"

	"github.com/lib/pq"
)

func SelectFrom(sel []string, table string, from string) string {
	selString := ""
	for i := 0; i < len(sel); i++ {
		if i == len(sel)-1 {
			selString += sel[i]
		} else {
			selString += ", " + sel[i]
		}
	}

	return "SELECT " + selString + " FROM" + table
}

func ArrayContain(column string, parameter *[]any, filter *[]string, operator string, fallback []string) string {
	var f []string
	if filter != nil && len(*filter) > 0 {
		f = *filter
	} else if len(fallback) > 0 {
		f = fallback
	}

	if len(f) == 0 {
		return ""
	}

	*parameter = append(*parameter, pq.Array(f))
	paramNumber := len(*parameter)
	return fmt.Sprintf("%s %s = ANY($%d)\n", operator, column, paramNumber)
}

func Pagination(page *int, limit *int) string {
	paginationQuery := ""

	if limit != nil && *limit > 0 {
		paginationQuery += fmt.Sprintf(" LIMIT %d \n", *limit)
		if page != nil && *page > 0 {
			offset := (*page - 1) * *limit
			paginationQuery += fmt.Sprintf(" OFFSET %d \n", offset)
		}
	}

	return paginationQuery
}

func DefaultOrder(entiryName string) string {
	if entiryName != "" {
		return fmt.Sprintf("ORDER BY %screate_date DESC\n", entiryName)
	}
	return "ORDER BY create_date DESC\n"
}

type SearchQueryItem struct {
	Value    string
	CastType string
	Weight   int
}

func Search(values []SearchQueryItem, parameter *[]any, search *string, searchLike bool, operator string) (matchScoreQuery string, searchQuery string) {
	if search == nil || *search == "" {
		return "0", ""
	}

	*parameter = append(*parameter, search)
	paramNumber := len(*parameter)

	if searchLike {
		*parameter = append(*parameter, fmt.Sprintf("%%%s%%", *search))
	}
	ilikeParamNumber := len(*parameter)

	similarityScore := []string{}
	searchItems := []string{}

	for _, e := range values {
		value := e.Value
		castType := ""
		if e.CastType != "" {
			castType = fmt.Sprintf(`::%s`, e.CastType)
		}
		weight := e.Weight

		/// Process similarity score
		similarityScore = append(similarityScore, fmt.Sprintf(
			`(COALESCE(similarity(%s%s, $%d), 0) * %d)`,
			value, castType, paramNumber, weight,
		))

		/// Process fussy search query statement
		searchItems = append(searchItems, fmt.Sprintf(
			`%s%s %% $%d`,
			value, castType, paramNumber,
		))

		/// Process fussy search query statement
		if searchLike {
			searchItems = append(searchItems, fmt.Sprintf(
				`%s%s ILIKE $%d`,
				value, castType, ilikeParamNumber,
			))
		}

	}

	matchScoreQuery = strings.Join(similarityScore, " +\n")
	searchQuery = fmt.Sprintf("%s (\n %s \n )\n",
		operator,
		strings.Join(searchItems, "\nOR "),
	)
	return matchScoreQuery, searchQuery
}
