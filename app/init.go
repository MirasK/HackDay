package app

import (
	"hackday/db"
	"os"
)

// InitProg initialise programm data
func InitProg() error {
	timeLayout = "2006-01-02 15:04:05"
	logFile, e = os.Create("logs.txt")
	codes = map[string]string{}
	if e != nil {
		return e
	}

	e = db.Conn()
	if e != nil {
		return e
	}
	return nil
}
