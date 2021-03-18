package main

import (
	"backorder_updater/cmd/syncup"
	"backorder_updater/internal/pkg"
	"backorder_updater/internal/pkg/logging"
	"backorder_updater/internal/pkg/sync"
	"backorder_updater/internal/pkg/sync/mysql"
	"database/sql"
	_ "embed"
	"log"
	"os"
)

//go:embed skeleton.env
//goland:noinspection GoUnusedGlobalVariable
var dotEnvDefault string

//go:embed sqlite_init.sql
//goland:noinspection GoUnusedGlobalVariable
var sqliteInitSql string

func main() {
	if os.Getenv("BOSYNC_ENVIRONMENT") != "DEVELOPMENT" {
		logfile, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		pkg.CheckAndLogFatal(err)
		log.SetOutput(logfile)
		defer func() {
			pkg.CheckAndLog(logfile.Close())
		}()
	}

	envRoot, err := os.Getwd()
	pkg.CheckAndPanic(err)

	err = pkg.LoadEnv(envRoot, ".env")
	pkg.CheckAndPanic(err)

	sqliteCqr, err := sync.Initialise(sqliteInitSql)

	defer func() {
		pkg.CheckAndLogFatal(sqliteCqr.Close())
	}()

	mysqlConn, err := sql.Open("mysql", os.Getenv("BOSYNC_MYSQL_CONN_STR"))
	pkg.CheckAndPanic(err)
	defer func() {
		pkg.CheckAndLogFatal(mysqlConn.Close())
	}()

	mysqlCqr := mysql.NewCommandQueryRepository(mysqlConn)
	defer func() {
		pkg.CheckAndLogFatal(mysqlCqr.Close())
	}()

	logger := logging.NewSqliteLogger(sqliteCqr)

	pkg.CheckAndLogFatal(err)

	syncup.Run(sqliteCqr, mysqlCqr, logger)

}
