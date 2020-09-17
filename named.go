package dblib

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/jmoiron/sqlx"
)

// NewNamedQuery register a named query (SQLX) for DB checking
func NewNamedQuery(dbNum int, query string) string {
	query = strings.TrimSpace(query)

	namedQueryMut.Lock()
	defer namedQueryMut.Unlock()

	_, existed := namedQueries[dbNum]
	if !existed {
		namedQueries[dbNum] = make(database)
	}

	_, existed = namedQueries[dbNum][query]
	if existed {
		panic("must use NewNamedQuery for global variables")
	}
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		panic("runtime.Caller failed")
	}
	namedQueries[dbNum][query] = fmt.Sprintf("%s:%d", file, line)

	return query
}

func checkNamedQueries(db *sqlx.DB, dbNum int, msg string) string {
	namedQueryMut.Lock()
	defer namedQueryMut.Unlock()

	if namedQueries[dbNum] == nil {
		return msg
	}

	for query, line := range namedQueries[dbNum] {
		stmt, err := db.PrepareNamed(query)
		if err != nil {
			msg += fmt.Sprintf(`
=============================================================
%s
-------------------------------------------------------------
%s
-------------------------------------------------------------
%v
`, line, query, err)
			continue
		}
		_ = stmt.Close()
	}
	return msg
}
