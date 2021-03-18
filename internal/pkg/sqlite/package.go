package sqlite

import (
	"database/sql"
	_ "embed"
	"os"
)

func Initialise(initSql string) (*CommandQueryRepository, error) {
	db, err := getDefaultConn()

	if err != nil {
		return nil, err
	}

	_, err = db.Exec(initSql)

	if err != nil {
		return nil, err
	}

	return &CommandQueryRepository{conn: db}, nil
}

func getDefaultConn() (*sql.DB, error) {
	conn, err := sql.Open("sqlite3", os.Getenv("BOSYNC_SQLITE_DB_PATH"))
	return conn, err
}
