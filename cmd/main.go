package main

import (
	"github.com/QuangTung97/go-dblib"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

const primaryDB = 1
const secondaryDB = 2

var simpleQuery string = dblib.NewQuery(primaryDB, `
INSERT INTO counter(id, value) VALU (?, ?)
`)

var selectQuery string = dblib.NewQuery(primaryDB, `
SELECT id, value FROM counter
`)

var copyQuery string = dblib.NewQuery(primaryDB, `
INSERT INTO counter(tung, valu)
SELECT id, value FROM counter
`)

var someNamedQuery string = dblib.NewNamedQuery(primaryDB, `
INSERT INTO counter(id, value) VALUES(:id, :value)
`)

var anotherQuery string = dblib.NewQuery(secondaryDB, `
INSERT INTO tung(id, name) VALUES (1, 2)
`)

func main() {
	db := sqlx.MustConnect("mysql", "root:1@/bench")
	dblib.CheckRegisteredQueries(db, primaryDB)

	v := struct {
		ID    int `db:"id"`
		Value int `db:"value"`
	}{
		ID:    3,
		Value: 100,
	}

	db.MustExec(copyQuery)

	_, err := db.NamedExec(someNamedQuery, v)
	if err != nil {
		panic(err)
	}
}
