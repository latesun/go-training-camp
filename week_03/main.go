package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

func launchHTTPServer(ctx context.Context, exit chan struct{}) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	// 通过接口来关闭程序
	mux.HandleFunc("/shutdown", func(w http.ResponseWriter, r *http.Request) {
		exit <- struct{}{}
	})

	svr := http.Server{Addr: ":8081", Handler: mux}

	go func() {
		select {
		case <-ctx.Done():
			svr.Shutdown(ctx)
		}
	}()

	log.Println("listen 8081...")
	err := svr.ListenAndServe()
	if err != nil {
		exit <- struct{}{}
	}

	return err
}

func main() {
	// stop 用于捕获外部终止信号
	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)

	// exit 用户捕获服务端退出信号
	exit := make(chan struct{}, 1)
	ctx, cancel := context.WithCancel(context.Background())

	group, _ := errgroup.WithContext(ctx)
	group.Go(func() error {
		return launchHTTPServer(ctx, exit)
	})

	group.Go(func() error {
		select {
		case sig := <-stop:
			cancel()
			return fmt.Errorf("catch stop signal: %v", sig)
		case <-exit:
			cancel()
			return errors.New("internel server error")
		}
	})

	log.Fatal("exit reason: ", group.Wait())
}
