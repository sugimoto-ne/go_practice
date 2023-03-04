package infrastracture

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sugimoto-ne/go_practice.git/config"
)

type MyServer struct {
	Srv *http.Server
}

func NewServer(cfg *config.Config, mux http.Handler) *MyServer {
	server := &MyServer{
		Srv: &http.Server{
			Addr:              fmt.Sprintf(":%d", cfg.Port),
			Handler:           mux,
			ReadHeaderTimeout: time.Second * time.Duration(cfg.ReadHeaderTimeoutSecond),
		},
	}

	return server
}

func (ms *MyServer) Run(ctx context.Context) error {
	closed := make(chan error)
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM, os.Kill)
	defer stop()

	go func(ctx context.Context) {
		<-ctx.Done()
		println("キャンセル通知！")
		if err := ms.Srv.Shutdown(context.Background()); err != nil {
			// シャットダウンでエラーが発生した場合にチャネルエラー文言を送信する
			closed <- fmt.Errorf("server shutdown error: %v", err)
		} else {
			println("エラーなし")
			close(closed)
		}
	}(ctx)

	err := ms.Srv.ListenAndServe()
	shutdownLog := <-closed
	if err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("server close error: %v", err)
	} else {
		return shutdownLog
	}
}
