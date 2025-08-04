package main

import (
	"context"
	"github.com/AndreySirin/04.08/internal/server"
	"github.com/AndreySirin/04.08/internal/zip"
	"log/slog"
	"os"
	"os/signal"
	"sync"
)

const port = ":8080"

func main() {
	lg := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	err := zip.InitDir()
	if err != nil {
		lg.Error("error creating archives directory")
	}
	srv := server.New(port, lg)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	wg := new(sync.WaitGroup)
	wg.Add(2)

	go func() {
		srv.Run()
		wg.Done()
	}()

	go func() {
		<-ctx.Done()
		srv.ShutDown()
		wg.Done()
	}()
	wg.Wait()
}
