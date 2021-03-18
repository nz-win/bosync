package sqlite

import (
	"backorder_updater/internal/pkg"
	"database/sql"
	"errors"
	sqlite "github.com/mattn/go-sqlite3"
)

type CommandQueryRepository struct {
	conn *sql.DB
}

type RecordDataLoadResult struct {
	duplicate bool
}

func (cq *CommandQueryRepository) Close() error {
	return cq.conn.Close()
}

func (cq *CommandQueryRepository) RecordNewDataLoad(sha256hash string) {
	query, err := cq.getSql("RECORD_NEW_DATALOAD")
	pkg.CheckAndLogFatal(err)

	tx, err := cq.conn.Begin()
	pkg.CheckAndLogFatal(err)

	stmt, err := tx.Prepare(query)
	pkg.CheckAndLogFatal(err)
	defer func() {
		err = stmt.Close()
	}()

	var e = &sqlite.Error{}

	if _, err = stmt.Exec(sha256hash); err != nil && errors.As(err, &e) {
		switch e.ExtendedCode {
		case sqlite.ErrConstraintForeignKey:

		}
	} else {

	}
	err = tx.Commit()
	pkg.CheckAndLogFatal(err)

}

func (cq *CommandQueryRepository) GetPreviousDataLoadHash(out *string) error {
	query, err := cq.getSql("RETRIEVE_LAST_DATALOAD_HASH")

	if err != nil {
		return err
	}

	err = cq.conn.QueryRow(query).Scan(&out)

	if err != nil {
		return err
	}

	return nil
}

func (cq *CommandQueryRepository) getSql(queryName string) (string, error) {
	stmt, err := cq.conn.Prepare(`SELECT query FROM queries WHERE query_name = ? LIMIT 1`)
	if err != nil {
		return "", err
	}

	defer func() {
		err = stmt.Close()
	}()

	var sqlQuery string
	err = stmt.QueryRow(queryName).Scan(&sqlQuery)

	if err != nil {
		return "", err
	}

	return sqlQuery, nil
}