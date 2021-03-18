package main

import (
	"backorder_updater/cmd/syncup/internal"
	"backorder_updater/internal/pkg"
	"backorder_updater/internal/pkg/sqlite"
	"crypto/sha256"
	"database/sql"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

//go:embed .env_defaults
//goland:noinspection GoUnusedGlobalVariable
var dotEnvDefault string

//go:embed sqlite_init.sql
//goland:noinspection GoUnusedGlobalVariable
var sqliteInitSql string

func main() {
	var previousHash = ""
	var currentHash = ""
	var isNewData = false
	var response = internal.ApiResponse{}

	if os.Getenv("BOSYNC_ENV") != "DEVELOPMENT" {
		logfile, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		pkg.CheckAndLogFatal(err)
		log.SetOutput(logfile)
		defer func() {
			pkg.CheckAndLogFatal(logfile.Close())
		}()
	}

	envRoot, err := os.Getwd()
	pkg.CheckAndPanic(err)

	err = pkg.LoadEnv(envRoot, ".env")
	pkg.CheckAndPanic(err)

	sqliteCqr, err := sqlite.Initialise(sqliteInitSql)

	defer func() {
		pkg.CheckAndLogFatal(sqliteCqr.Close())
	}()

	mysql, err := sql.Open("mysql", os.Getenv("BOSYNC_MYSQL_CONN_STR"))
	pkg.CheckAndPanic(err)
	defer func() {
		pkg.CheckAndLogFatal(mysql.Close())
	}()

	resp, err := http.Get(os.Getenv("BOSYNC_API_ENDPOINT_URL"))
	pkg.CheckAndPanic(err)

	defer func() {
		pkg.CheckAndLogFatal(resp.Body.Close())
	}()

	bodyBytes, err := io.ReadAll(resp.Body)
	pkg.CheckAndPanic(err)

	currentHash = fmt.Sprintf("%x", sha256.Sum256(bodyBytes))

	if err = sqliteCqr.GetPreviousDataLoadHash(&previousHash); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			isNewData = true
		} else {
			pkg.CheckAndPanic(err)
		}
	}

	if strings.Compare(currentHash, previousHash) != 0 {
		isNewData = true
	}

	err = json.Unmarshal(bodyBytes, &response)
	pkg.CheckAndPanic(err)

	if isNewData {
		internal.UpdateMysql(mysql, response.Data)
		sqliteCqr.RecordNewDataLoad(currentHash)
	}

}
