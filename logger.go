package dailylogrus

import (
	"github.com/sirupsen/logrus"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"time"
)

var Logger = logrus.New()

var done = make(chan bool)

func InitLogger(dir string) {

	log := logrus.New()
	// log format config
	log.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: time.RFC3339,
		FullTimestamp:   true,
	})

	LogFileWriter := getFileWriter(getLogFilePath(dir))

	mw := io.MultiWriter(os.Stdout, LogFileWriter)

	log.SetOutput(mw)

	ticker := time.NewTicker(time.Second)

	go func() {
		for {
			select {
			case <-done:
				ticker.Stop()
				return
			case <-ticker.C:
				filePath := getLogFilePath(dir)
				_, err := os.Stat(filePath)
				if err != nil && LogFileWriter.Filename != filePath {
					// Create a new log file with the current timestamp
					NewLogFileWritter := getFileWriter(filePath)
					// Switch the logger to the new file
					LogFileWriter.Close()
					LogFileWriter = NewLogFileWritter
					mw := io.MultiWriter(os.Stdout, NewLogFileWritter)
					log.SetOutput(mw)
				}
			}
		}
	}()

	Logger = log
}

func LogStop() {
	done <- true
}

func getFileWriter(filePath string) *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   filePath,
		MaxSize:    50,
		MaxBackups: 7,
		MaxAge:     30,
		Compress:   true,
	}
}

func getLogFilePath(dir string) string {
	fileDir := os.Getenv("HOME") + "/log/" + dir
	fileName := time.Now().Format("2006-01-02") + ".log"
	return fileDir + "/" + fileName
}
