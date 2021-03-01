package sqlutil

import sq "github.com/Masterminds/squirrel"

// LimitOffsetPagination defines pagination by limit and offset
type LimitOffsetPagination struct {
	Limit  int
	Offset int
}

// UseLimitOffsetPagination adds the limit and offset defined in LimitOffsetPagination to a sql query
func UseLimitOffsetPagination(q sq.SelectBuilder, p LimitOffsetPagination) sq.SelectBuilder {
	switch {
	case p.Limit > 0:
		q = q.Limit(uint64(p.Limit))
		break
	case p.Limit < 0:
		break
	default:
		q = q.Limit(25)
	}

	if p.Offset > 0 {
		q = q.Offset(uint64(p.Offset))
	}

	return q
}
