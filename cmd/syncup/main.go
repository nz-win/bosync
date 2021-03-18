package syncup

import (
	"backorder_updater/internal/pkg"
	"backorder_updater/internal/pkg/logging"
	"backorder_updater/internal/pkg/sync/mysql"
	"backorder_updater/internal/pkg/sync/sqlite"
	"backorder_updater/internal/pkg/types"
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

func Run(sqliteCqr *sqlite.CommandQueryRepository, mysqlCqr *mysql.CommandQueryRepository, logger *logging.SqliteLogger) {
	var previousHash = ""
	var currentHash = ""
	var isNewData = false
	var response = types.ApiResponse{}

	resp, err := http.Get(os.Getenv("BOSYNC_API_ENDPOINT_URL"))

	if err != nil {
		err = logger.LogAndNotify(types.Fatal,
			fmt.Sprintf("Unable to connect to API @ %s",
				os.Getenv("BOSYNC_API_ENDPOINT_URL")),
			err)
		return
	}

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

	isNewData = strings.Compare(currentHash, previousHash) != 0

	err = json.Unmarshal(bodyBytes, &response)
	pkg.CheckAndPanic(err)

	if isNewData {
		if err = mysqlCqr.UpdateBackOrders(response.Data); err != nil {
			logErr := sqliteCqr.InsertLog(err, types.Fatal)
			log.Println(logErr)
			log.Fatal(err)
		}
		sqliteCqr.RecordNewDataLoad(currentHash)
	}

}
