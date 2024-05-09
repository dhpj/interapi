package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
	"net/http"

	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
)

var Stdlog *log.Logger

var GoClient *http.Client = &http.Client{
	Timeout: time.Second * 30,
	Transport: &http.Transport{
		TLSHandshakeTimeout: 10 * time.Second,
	},
}

func InitConfig() {
	realpath, _ := os.Executable()
	dir := filepath.Dir(realpath)
	logDir := filepath.Join(dir, "logs")
	err := createDir(logDir)
	if err != nil {
		log.Fatalf("Failed to ensure log directory: %s", err)
	}
	path := filepath.Join(logDir, "DHNServer")
	loc, _ := time.LoadLocation("Asia/Seoul")
	writer, err := rotatelogs.New(
		fmt.Sprintf("%s-%s.log", path, "%Y-%m-%d"),
		rotatelogs.WithLocation(loc),
		rotatelogs.WithMaxAge(-1),
		rotatelogs.WithRotationCount(7),
	)

	if err != nil {
		log.Fatalf("Failed to Initialize Log File %s", err)
	}

	log.SetOutput(writer)
	stdlog := log.New(os.Stdout, "INFO -> ", log.Ldate|log.Ltime)
	stdlog.SetOutput(writer)
	Stdlog = stdlog

}

func createDir(dirName string) error {
	err := os.MkdirAll(dirName, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	return nil
}
