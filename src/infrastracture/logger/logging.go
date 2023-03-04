package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"golang.org/x/exp/slog"
)

type MyLogger struct {
	LogFilePath string
	Logger      *slog.Logger
}

func NewLogger(filepath, env, filename, ext string) (*MyLogger, error) {
	path := fmt.Sprintf("%s/%s", filepath, env)

	pathString := "."
	pathArr := strings.Split(path, "/")
	for _, pn := range pathArr {
		if pn == "." {
			continue
		} else {
			pathString = fmt.Sprintf("%s/%s", pathString, pn)
		}

		if _, err := os.Stat(pathString); os.IsNotExist(err) {
			os.Mkdir(pathString, 0777)
		}
	}

	file := "/" + filename + "." + ext
	logfile, err := os.OpenFile(pathString+file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("file=logfile err=%s", err.Error())

		return nil, err
	}
	multiLogFile := io.MultiWriter(os.Stdout, logfile)
	logger := slog.New(slog.HandlerOptions{
		AddSource: true,
	}.NewJSONHandler(multiLogFile))

	slog.SetDefault(logger)
	myLogger := &MyLogger{
		LogFilePath: path,
		Logger:      logger,
	}

	logger.Info("start logging")

	return myLogger, nil
}
