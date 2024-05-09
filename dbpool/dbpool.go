package dbpool

import(
	"fmt"
	"database/sql"
	"log"

	_ "github.com/denisenkom/go-mssqldb"
)

var DB *sql.DB


func InitDB() {
	// MSSQL 연결 정보
    server := "localhost"
    port := 1433
    user := "sa"
    password := "dhn7985!"
    database := "test"

    connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s", server, user, password, port, database)

	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
    if err != nil {
        panic(err.Error())
    }

	db.SetMaxIdleConns(50)
	db.SetMaxOpenConns(50)

	DB = db

}