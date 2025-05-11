package shared

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type DBHandler struct {
	Conn *sql.DB
}

func NewDBHandler(user, pass, host string) *DBHandler {
	dsn := user + ":" + pass + "@tcp(" + host + ")/"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Connection error:", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal("Ping error (can't connect to DB):", err)
	}
	return &DBHandler{Conn: db}
}

func (db *DBHandler) ExecQuery(query string) (int64, error) {
	res, err := db.Conn.Exec(query)
	if err != nil {
		return 0, err
	}
	affected, _ := res.RowsAffected()
	return affected, nil
}

func (db *DBHandler) QueryRows(query string) (*sql.Rows, error) {
	return db.Conn.Query(query)
}
