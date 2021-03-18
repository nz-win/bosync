package internal

import (
	"backorder_updater/internal/pkg/logging"
	"backorder_updater/internal/pkg/sync/mysql"
	"backorder_updater/internal/pkg/sync/sqlite"
	"fmt"
)

type AppContext struct {
	RichLogger *logging.SqliteLogger
	SqliteCQR  *sqlite.CommandQueryRepository
	MysqlCQR   *mysql.CommandQueryRepository
}

func (c *AppContext) Close() error {
	er1 := c.SqliteCQR.Close()
	er2 := c.MysqlCQR.Close()

	if er1 != nil && er2 != nil {
		return fmt.Errorf("multiple errors occured closing app context: %v %v", er1, er2)
	}

	if er1 != nil {
		return er1
	}

	if er2 != nil {
		return er2
	}

	return nil
}
