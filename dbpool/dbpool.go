package dbpool

import(
	"fmt"
	"database/sql"
	"log"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB
var DB2 *sql.DB


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

	db.SetMaxIdleConns(50)
	db.SetMaxOpenConns(50)

	DB = db

	//mariaDB 연결 정보
	server = "localhost"
	port = 3306
	user = "root"
	password = "sjk4556!!22"
	database = "song"

	connString = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", user, password, server, port, database)
	db, err = sql.Open("mysql", connString)
	if err != nil {
		log.Fatal(err)
	}

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(50)

	DB2 = db
}