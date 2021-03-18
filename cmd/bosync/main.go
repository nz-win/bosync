package main

import (
	"backorder_updater/cmd/bosync/commands/cmd_log"
	"backorder_updater/cmd/bosync/commands/cmd_sync"
	"backorder_updater/cmd/bosync/internal"
	"backorder_updater/internal/pkg"
	"backorder_updater/internal/pkg/logging"
	"backorder_updater/internal/pkg/sync"
	"log"

	"backorder_updater/internal/pkg/sync/mysql"
	"database/sql"
	_ "embed"
	"github.com/urfave/cli/v2"
	"os"
)

//go:embed resources/sqlite_init.sql
//goland:noinspection GoUnusedGlobalVariable
var sqliteInitSql string

func main() {
	setupEnv()
	ctx := buildContext()
	defer func() {
		pkg.CheckAndLogFatal(ctx.Close())
	}()

	app := &cli.App{
		Name: "Backorder Sync - Tools for managing backorder sync",
		Commands: []*cli.Command{
			cmd_log.NewCommand(ctx),
			cmd_sync.NewCommand(ctx),
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func buildContext() *internal.AppContext {
	sqliteCqr, err := sync.Initialise(sqliteInitSql)

	pkg.CheckAndPanic(err)

	mysqlConn, err := sql.Open("mysql", os.Getenv("BOSYNC_MYSQL_CONN_STR"))

	pkg.CheckAndPanic(err)
	mysqlCqr := mysql.NewCommandQueryRepository(mysqlConn)
	logger := logging.NewSqliteLogger(sqliteCqr)
	pkg.CheckAndLogFatal(err)

	return &internal.AppContext{
		RichLogger: logger,
		SqliteCQR:  sqliteCqr,
		MysqlCQR:   mysqlCqr,
	}
}

func setupEnv() {
	homeRoot, homeIsSet := os.LookupEnv("BOSYNC_HOME")

	if homeIsSet {
		err := pkg.LoadEnv(homeRoot, ".env")
		pkg.CheckAndPanic(err)
	} else {
		envRoot, err := os.Getwd()
		pkg.CheckAndPanic(err)

		err = pkg.LoadEnv(envRoot, ".env")
		pkg.CheckAndPanic(err)
	}

}
