package sync

import (
	"backorder_updater/internal/pkg/sync/sqlite"
	_ "embed"
	"github.com/jmoiron/sqlx"
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

func getDefaultConn() (*sqlx.DB, error) {
	conn, err := sqlx.Open("sqlite3", os.Getenv("BOSYNC_SQLITE_DB_PATH"))
	return conn, err
}
