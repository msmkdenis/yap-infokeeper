package main

import (
	"os"
	"os/signal"
	"syscall"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"

	"github.com/msmkdenis/yap-infokeeper/internal/app/infokeeper"
)

func main() {
	quitSignal := make(chan os.Signal, 1)
	signal.Notify(quitSignal, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	infokeeper.Run(quitSignal)
}
