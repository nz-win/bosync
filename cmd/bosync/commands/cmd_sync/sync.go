package cmd_sync

import (
	"backorder_updater/cmd/bosync/internal"
	"backorder_updater/internal/pkg"
	"backorder_updater/internal/pkg/types"
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/briandowns/spinner"
	colors "github.com/logrusorgru/aurora/v3"
	"github.com/urfave/cli/v2"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func NewCommand(ctx *internal.AppContext) *cli.Command {
	return &cli.Command{
		Name:   "sync",
		Usage:  "run backorder sync",
		Action: buildAction(ctx),
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "verbose",
				Usage:       "show verbose logs on stdout",
				Value:       false,
				DefaultText: "false",
			},
			&cli.BoolFlag{
				Name:        "skip-stale",
				Usage:       "set whether stale data loads should be skipped",
				Value:       false,
				DefaultText: "false",
			},
		},
	}
}

func buildAction(appCtx *internal.AppContext) func(ctx *cli.Context) error {
	return func(c *cli.Context) error {

		s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
		s.Suffix = " Retrieving latest API data"

		if c.Bool("verbose") {
			fmt.Println(colors.Bold(colors.Cyan("🔍 Checking For Fresh Data")))
			s.Start()
		}
		response, currentHash, err := getApiData(appCtx)

		if err != nil {
			appCtx.RichLogger.Log(fmt.Sprintf("unable to retrieve api data: %v", err), types.Error)
			if c.Bool("verbose") {
				fmt.Println("")
				fmt.Println(colors.Bold(colors.Red("✗ Data Sync Failed")))
			}
			return err
		}

		if c.Bool("verbose") {
			s.Stop()
		}

		dataIsFresh, err := isDataFresh(appCtx, currentHash)
		if err != nil {
			appCtx.RichLogger.Log(err, types.Error)
		}

		dataIsFresh = dataIsFresh || !c.Bool("skip-stale")

		if c.Bool("verbose") {
			switch dataIsFresh {
			case true:
				fmt.Println(colors.Bold(colors.Green("✓ Fresh Data Available")))
			default:
				fmt.Println(colors.Yellow("✗ No Fresh Data Available"))
			}

		}

		if dataIsFresh {
			s = spinner.New(spinner.CharSets[11], 100*time.Millisecond)
			s.Suffix = " Updating MySql Records"

			if c.Bool("verbose") {
				s.Start()
			}

			if err = appCtx.MysqlCQR.UpdateBackOrders(response.Data); err != nil {
				err = appCtx.RichLogger.LogAndNotify(types.Fatal, "An Error occurred inserting new back order records", err)

				if err != nil {
					log.Println(err)
				}

				return nil

			}
			if c.Bool("verbose") {
				s.Stop()
			}

			if c.Bool("verbose") {
				fmt.Println(colors.Cyan("✎ Recording Data Hash"))
			}
			appCtx.SqliteCQR.RecordNewDataLoad(currentHash)
		}

		if c.Bool("verbose") {
			fmt.Println(colors.Bold(colors.Green("★ Data Sync Complete")))
		}
		return nil
	}
}

func isDataFresh(ctx *internal.AppContext, currentHash string) (bool, error) {
	var lastSeenHash string

	if err := ctx.SqliteCQR.GetPreviousDataLoadHash(&lastSeenHash); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return true, nil
		} else {
			return false, err
		}
	}

	return currentHash == lastSeenHash, nil
}

func getApiData(ctx *internal.AppContext) (*types.ApiResponse, string, error) {

	endpoint, isSet := os.LookupEnv("BOSYNC_API_ENDPOINT_URL")

	if isSet == false || endpoint == "" {
		return nil, "", errors.New("BOSYNC_API_ENDPOINT_URL env var is not set")
	}

	resp, err := http.Get(endpoint)

	if err != nil {
		_ = ctx.RichLogger.LogAndNotify(types.Fatal,
			fmt.Sprintf("Unable to connect to API @ %s",
				os.Getenv("BOSYNC_API_ENDPOINT_URL")),
			err)
		return nil, "", err
	}
	defer func() {
		pkg.CheckAndLogFatal(resp.Body.Close())
	}()

	bodyBytes, err := io.ReadAll(resp.Body)
	pkg.CheckAndPanic(err)

	currentHash := fmt.Sprintf("%x", sha256.Sum256(bodyBytes))
	var response *types.ApiResponse
	err = json.Unmarshal(bodyBytes, &response)
	pkg.CheckAndPanic(err)

	return response, currentHash, nil

}
