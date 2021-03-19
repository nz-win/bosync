package cmd_info

import (
	"backorder_updater/cmd/bosync/internal"
	"backorder_updater/internal/pkg/types"
	_ "encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	colors "github.com/logrusorgru/aurora/v3"
	"github.com/urfave/cli/v2"
	"os"
	"time"
)

type InfoLog struct {
	LastFreshDataLoadAt time.Time    `json:"last_fresh_data_load" db:"created_at"`
	LastDataLoadAttempt time.Time    `json:"last_data_load_attempt" db:"last_seen_at"`
	LatestLogs          *[]types.Log `json:"latest_logs"`
}

func NewCommand(ctx *internal.AppContext) *cli.Command {
	return &cli.Command{
		Name:   "info",
		Usage:  "print various pieces of information about bosync",
		Action: buildAction(ctx),
	}
}

func buildAction(ctx *internal.AppContext) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		loc, _ := time.LoadLocation("Local")

		db, err := sqlx.Open("sqlite3", os.Getenv("BOSYNC_SQLITE_DB_PATH"))
		if err != nil {
			panic(err)
		}

		var infoResults []InfoLog
		var recentLogs []types.Log

		err = db.Select(&infoResults, `SELECT created_at,last_seen_at FROM dataloads ORDER BY created_at DESC LIMIT 1`)
		if err != nil {
			panic(err)
		}

		i := infoResults[0]

		err = db.Select(&recentLogs, `SELECT * FROM logs WHERE created_at >= ? ORDER BY created_at DESC LIMIT 10`, i.LastDataLoadAttempt)
		if err != nil {
			panic(err)
		}

		fmt.Println(
			fmt.Sprintf("%s\t-\t%s",
				colors.Bold(colors.Cyan("Last Fresh Data Load")),
				i.LastFreshDataLoadAt.In(loc).Format("Mon Jan _2 15:04:05 2006")),
		)

		fmt.Println(colors.Bold(colors.Cyan("\nRecent Logs:")))

		for _, l := range recentLogs {
			createdAt := l.CreatedAt.In(loc).Format("Mon Jan _2 15:04:05 2006")
			format := "\t[%s - %s]\t%s\n"
			switch l.Level {
			case types.Fatal:
				fmt.Printf(format, createdAt, colors.Bold(colors.Red(l.Level)), l.Message)
			case types.Error:
				fmt.Printf(format, createdAt, colors.Red(l.Level), l.Message)
			case types.Warn:
				fmt.Printf(format, createdAt, colors.Yellow(l.Level), l.Message)
			default:
				fmt.Printf(format, createdAt, colors.Cyan(l.Level), l.Message)
			}
		}

		return nil
	}
}
