package sqlutil

import (
	"fmt"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/iconmobile-dev/go-core/errors"
)

// OneColumnSort defines what column should be sorted in what order
type OneColumnSort struct {
	Column string
	Order  string
}

// UseOneColumnSort adds the sort column and sort order defined in OneColumnSort to a sql query
func UseOneColumnSort(q sq.SelectBuilder, p OneColumnSort, columnMapping map[string]string) (sq.SelectBuilder, error) {
	var column string
	switch {
	case strings.TrimSpace(p.Column) != "":
		column = strings.TrimSpace(p.Column)
	default:
		column = "id"
	}

	var order string
	switch {
	case strings.ToLower(strings.TrimSpace(p.Order)) == "desc":
		order = "DESC"
	default:
		order = "ASC"
	}

	mColumn, ok := columnMapping[column]
	if !ok {
		return q, errors.E(fmt.Errorf("can't sort by column %v", p.Column), errors.Unprocessable, fmt.Sprintf("can't sort by column %v", p.Column))
	}

	q = q.OrderBy(fmt.Sprintf("%v %v", mColumn, order))

	return q, nil
}
