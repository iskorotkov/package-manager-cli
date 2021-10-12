package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/iskorotkov/package-manager-cli/internal/commands"
	"github.com/iskorotkov/package-manager-cli/internal/keys"
	"github.com/iskorotkov/package-manager-cli/pkg/xlog"
)

const version = "dev"

func main() {
	flush := setupLogger()
	defer flush()

	xlog.Push(version)
	defer xlog.Pop()

	commands.Execute()
}

func setupLogger() func() {
	if err := os.MkdirAll(keys.LogsPath, keys.LogsPermissions); err != nil {
		log.Fatalf("error creating logs directory: %v", err)
	}

	filename := logFilename(keys.LogsPath)

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, keys.LogsPermissions)
	if err != nil {
		log.Fatalf("error creating/opening log file: %v", err)
	}

	log.SetOutput(file)
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile | log.Lmsgprefix)

	return func() {
		if err := file.Close(); err != nil {
			log.Fatalf("error closing log file: %v", err)
		}
	}
}

func logFilename(dir string) string {
	logFile := filepath.Base(os.Args[0])
	logFile = fmt.Sprintf("%s.log", logFile)
	logFile = filepath.Join(dir, logFile)

	return logFile
}
