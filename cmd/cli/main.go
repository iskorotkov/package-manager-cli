package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/iskorotkov/package-manager-cli/internal/commands"
	"github.com/iskorotkov/package-manager-cli/internal/keys"
)

const version = "dev"

func main() {
	flush := setupLogger()
	defer flush()

	commands.Execute()
}

func setupLogger() func() {
	if err := os.MkdirAll(keys.LogsPath, keys.LogsPermissions); err != nil {
		log.Fatalf("error creating logs directory: %v", err)
	}

	logFile := filepath.Join(keys.LogsPath, "cli.log")

	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, keys.LogsPermissions)
	if err != nil {
		log.Fatalf("error creating/opening log file: %v", err)
	}

	log.SetOutput(file)
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile | log.Lmsgprefix)
	log.SetPrefix(fmt.Sprintf("[%s] ", version))

	return func() {
		if err := file.Close(); err != nil {
			log.Fatalf("error closing log file: %v", err)
		}
	}
}
