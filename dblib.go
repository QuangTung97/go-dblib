package dblib

import (
	"sync"

	"github.com/jmoiron/sqlx"
)

type database map[string]string

var normalQueryMut sync.Mutex
var normalQueries = make(map[int]database)

var namedQueryMut sync.Mutex
var namedQueries = make(map[int]database)

// CheckRegisteredQueries uses perpared statements to check
// syntax and semantic
func CheckRegisteredQueries(db *sqlx.DB, dbNum int) {
	msg := ""
	msg = checkNormalQueries(db, dbNum, msg)
	msg = checkNamedQueries(db, dbNum, msg)
	if msg != "" {
		msg += `
=============================================================`
		panic(msg)
	}
}
