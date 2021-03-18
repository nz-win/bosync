package sync

import (
	"backorder_updater/internal/pkg/sync/sqlite"
	"database/sql"
	_ "embed"
	"os"
)

func Initialise(initSql string) (*sqlite.CommandQueryRepository, error) {
	db, err := getDefaultConn()

	if err != nil {
		return nil, err
	}

	_, err = db.Exec(initSql)

	if err != nil {
		return nil, err
	}

	return sqlite.NewCommandQueryRepository(db), nil
}

func getDefaultConn() (*sql.DB, error) {
	conn, err := sql.Open("sqlite3", os.Getenv("BOSYNC_SQLITE_DB_PATH"))
	return conn, err
}
