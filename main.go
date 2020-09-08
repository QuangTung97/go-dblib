package main

import (
	"fmt"
	"strings"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var queryMut sync.Mutex
var querySet = make(map[string]struct{})

var namedQueryMut sync.Mutex
var namedQuerySet = make(map[string]struct{})

func NewQuery(query string) string {
	query = strings.TrimSpace(query)

	queryMut.Lock()
	defer queryMut.Unlock()

	_, existed := querySet[query]
	if existed {
		panic("must use NewQuery for global variables")
	}
	querySet[query] = struct{}{}

	return query
}

func NewNamedQuery(query string) string {
	query = strings.TrimSpace(query)

	namedQueryMut.Lock()
	defer namedQueryMut.Unlock()

	_, existed := namedQuerySet[query]
	if existed {
		panic("must use NewNamedQuery for global variables")
	}
	namedQuerySet[query] = struct{}{}

	return query
}

func checkNormalQueries(db *sqlx.DB) {
	queryMut.Lock()
	defer queryMut.Unlock()

	for query := range querySet {
		stmt, err := db.Preparex(query)
		if err != nil {
			msg := fmt.Sprintf(`
=============================================================
%s
=============================================================
%v`, query, err)
			panic(msg)
		}
		_ = stmt.Close()
	}
}

func checkNamedQueries(db *sqlx.DB) {
	for query := range namedQuerySet {
		stmt, err := db.PrepareNamed(query)
		if err != nil {
			msg := fmt.Sprintf(`
=============================================================
%s
=============================================================
%v`, query, err)
			panic(msg)
		}
		_ = stmt.Close()
	}
}

func CheckRegisteredQueries(db *sqlx.DB) {
	checkNormalQueries(db)
	checkNamedQueries(db)
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
	db := sqlx.MustConnect("mysql", "root:tung@/bench")
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
