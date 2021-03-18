package types

import (
	"fmt"
	_ "github.com/jmoiron/sqlx"
	"time"
)

type LogLevel string

const (
	Info  = "INF"
	Debug = "DBG"
	Warn  = "WAR"
	Error = "ERR"
	Fatal = "FTL"
)

type Log struct {
	CreatedAt time.Time `db:"created_at"`
	Level     LogLevel  `db:"log_level"`
	Message   string    `db:"message"`
}

func (l *Log) String() string {
	return fmt.Sprintf("%v %s %s", l.CreatedAt, l.Level, l.Message)
}

type ApiResponse struct {
	Data    []BackOrder `json:"data"`
	Query   string      `json:"query"`
	Records int64       `json:"records"`
	Status  string      `json:"status"`
}

type Date struct {
	Date     string `json:"date"`
	Timezone string `json:"timezone"`
}

type BackOrder struct {
	AdmNo          string `json:"adm_no"`
	AreaNo         string `json:"area_no"`
	BackorderQty   int64  `json:"backorder_qty"`
	BusinessAreaNo string `json:"business_area_no"`
	MatAvailDate   Date   `json:"matAvailDate"`
	Material       string `json:"material"`
	Name           string `json:"name"`
	SalesDate      Date   `json:"salesDate"`
	SalesDoc       string `json:"salesDoc"`
	SoldToParty    string `json:"sold_to_party"`
}
