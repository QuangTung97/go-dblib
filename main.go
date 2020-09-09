package main

import (
	"fmt"
	"runtime"
	"strings"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var queryMut sync.Mutex
var queries = make(map[string]string)

var namedQueryMut sync.Mutex
var namedQueries = make(map[string]string)

// NewQuery register a query for DB checking
func NewQuery(query string) string {
	query = strings.TrimSpace(query)

	queryMut.Lock()
	defer queryMut.Unlock()

	_, existed := queries[query]
	if existed {
		panic("must use NewQuery for global variables")
	}
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		panic("runtime.Caller failed")
	}
	queries[query] = fmt.Sprintf("%s:%d", file, line)

	return query
}

// NewNamedQuery register a named query (SQLX) for DB checking
func NewNamedQuery(query string) string {
	query = strings.TrimSpace(query)

	namedQueryMut.Lock()
	defer namedQueryMut.Unlock()

	_, existed := namedQueries[query]
	if existed {
		panic("must use NewNamedQuery for global variables")
	}
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		panic("runtime.Caller failed")
	}
	namedQueries[query] = fmt.Sprintf("%s:%d", file, line)

	return query
}

func checkNormalQueries(db *sqlx.DB, msg string) string {
	queryMut.Lock()
	defer queryMut.Unlock()

	for query, line := range queries {
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

func checkNamedQueries(db *sqlx.DB, msg string) string {
	for query, line := range namedQueries {
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

// CheckRegisteredQueries uses perpared statements to check
// syntax and semantic
func CheckRegisteredQueries(db *sqlx.DB) {
	msg := ""
	msg = checkNormalQueries(db, msg)
	msg = checkNamedQueries(db, msg)
	if msg != "" {
		msg += `
=============================================================`
		panic(msg)
	}
}

var simpleQuery string = NewQuery(`
INSERT INTO counter(id, value) VALUES (?, ?)
`)

var copyQuery string = NewQuery(`
INSERT INTO counter(id, value)
SELECT id, value FROM counter
`)

var someNamedQuery string = NewNamedQuery(`
INSERT INTO counter(id, value) VALUES(:id, :value)
`)

func main() {
	db := sqlx.MustConnect("mysql", "root:1@/bench")
	CheckRegisteredQueries(db)

	v := struct {
		ID    int `db:"id"`
		Value int `db:"value"`
	}{
		ID:    3,
		Value: 100,
	}

	_, err := db.NamedExec(someNamedQuery, v)
	if err != nil {
		panic(err)
	}
}
