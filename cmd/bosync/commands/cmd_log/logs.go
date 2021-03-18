package cmd_log

import (
	"backorder_updater/cmd/bosync/internal"
	"backorder_updater/internal/pkg/types"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/urfave/cli/v2"
	"os"
)

func NewCommand(ctx *internal.AppContext) *cli.Command {
	return &cli.Command{
		Name:   "logs",
		Usage:  "show logs",
		Action: buildAction(ctx),
	}
}

func buildAction(appCtx *internal.AppContext) func(ctx *cli.Context) error {
	return func(c *cli.Context) error {
		var logs []*types.Log
		if success := tryOpenLogs(&logs); success == false {
			return errors.New("failed to open logs")
		}

		//goland:noinspection GoNilness
		for _, log := range logs {
			fmt.Println(log.String())
		}

		return nil
	}
}

func tryOpenLogs(out *[]*types.Log) bool {

	filePath, isSet := os.LookupEnv("BOSYNC_SQLITE_DB_PATH")

	if isSet == false || filePath == "" {
		panic("BOSYNC_SQLITE_DB_PATH env var missing")
	}

	db, err := sqlx.Open("sqlite3", os.Getenv("BOSYNC_SQLITE_DB_PATH"))

	if err != nil {
		return false
	}

	err = db.Select(out, "SELECT * FROM logs ORDER BY created_at DESC")

	if err != nil {
		return false
	}

	return true

}
