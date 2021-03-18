package logging

import (
	"backorder_updater/internal/pkg/sync/sqlite"
	"backorder_updater/internal/pkg/types"
	"fmt"
	"gopkg.in/gomail.v2"
	"log"
	"os"
	"strconv"
)

type SqliteLogger struct {
	cqr *sqlite.CommandQueryRepository
}

func NewSqliteLogger(cqr *sqlite.CommandQueryRepository) *SqliteLogger {
	return &SqliteLogger{cqr: cqr}
}

func (l *SqliteLogger) Log(m interface{}, level types.LogLevel) error {
	return l.cqr.InsertLog(m, level)
}

func (l *SqliteLogger) LogAndNotify(level types.LogLevel, notification string, m interface{}) error {

	port, err := strconv.Atoi(os.Getenv("BOSYNC_SMTP_PORT"))

	if err != nil {
		log.Println("error retrieving smtp port. defaulting to 25 ", err)
		port = 25
	}

	d := gomail.NewDialer(
		os.Getenv("BOSYNC_SMTP_HOST"),
		port,
		os.Getenv("BOSYNC_SMTP_USER"),
		os.Getenv("BOSYNC_SMTP_PASS"))

	to := os.Getenv("BOSYNC_MAILTO")
	from := os.Getenv("BOSYNC_MAILFROM")

	msg := gomail.NewMessage()
	msg.SetHeader("From", from)
	msg.SetHeader("To", to)

	switch level {
	case types.Fatal:
		msg.SetHeader("Subject", "Fatal Error Notification From Backorder Sync")
	case types.Error:
		msg.SetHeader("Subject", "Error Notification From Backorder Sync")
	case types.Warn:
		msg.SetHeader("Subject", "Warning Notification From Backorder Sync")
	case types.Debug:
		msg.SetHeader("Subject", "Debug Notification From Backorder Sync")
	default:
		msg.SetHeader("Subject", "Notification From Backorder Sync")
	}

	msg.SetBody("text/plain", fmt.Sprintf("%s\nDetails:\n%v", notification, m))

	//d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	err = l.Log(m, level)

	if err != nil {
		return err
	}

	return d.DialAndSend(msg)

}
