package dblib

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/jmoiron/sqlx"
)

// NewQuery register a query for DB checking
func NewQuery(dbNum int, query string) string {
	query = strings.TrimSpace(query)

	normalQueryMut.Lock()
	defer normalQueryMut.Unlock()

	_, existed := normalQueries[dbNum]
	if !existed {
		normalQueries[dbNum] = make(database)
	}

	_, existed = normalQueries[dbNum][query]
	if existed {
		panic("must use NewQuery for global variables")
	}
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		panic("runtime.Caller failed")
	}
	normalQueries[dbNum][query] = fmt.Sprintf("%s:%d", file, line)

	return query
}

func checkNormalQueries(db *sqlx.DB, dbNum int, msg string) string {
	normalQueryMut.Lock()
	defer normalQueryMut.Unlock()

	if normalQueries[dbNum] == nil {
		return msg
	}

	for query, line := range normalQueries[dbNum] {
		stmt, err := db.Preparex(query)
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
