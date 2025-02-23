package logger

import (
	"log"

	"github.com/natefinch/lumberjack"
)

var Log *log.Logger

func init() {
	/// Log file configuration
	logpath := "log/logfile.log"

	// Set up lumberjack logger
	logFile := &lumberjack.Logger{
		Filename:   logpath,
		MaxSize:    25,   // MB before rotation
		MaxBackups: 6,    // Number of backups to retain
		MaxAge:     30,   // Days before old logs are deleted
		Compress:   true, // Compress old log files
	}

	Log = log.New(logFile, "", log.LstdFlags|log.Lshortfile)
}
