package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/sugimoto-ne/go_practice.git/config"
	"github.com/sugimoto-ne/go_practice.git/testutil"
)

func TestRun(t *testing.T) {
	t.Run("testOverTimeout", func(t *testing.T) {
		resultShutdownError := make(chan error, 1)
		port := "3333"
		//1秒でtimeout
		readHeaderTimeoutSecond := "1"
		os.Setenv("PORT", port)
		os.Setenv("READ_HEADER_TIMEOUT_SECOND", readHeaderTimeoutSecond)
		cfg, err := config.NewConfig()
		if err != nil {
			t.Errorf("want no error, but got %v", err)
		}

		ctx, cancel := context.WithCancel(context.Background())

		mux := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			io.WriteString(w, "<html>foo</html>")
		})
		server := NewServer(cfg, mux)

		errShutdownCh := make(chan error, 1)

		go func() {
			// サーバが停止したら返却されるerrorをチャネルに送信
			errShutdownCh <- server.Run(ctx)
		}()

		// ReadHeaderTimeout 超過
		wantErr := "EOF"
		refusedErr := fmt.Sprintf("dial tcp 127.0.0.1:%s: connect: connection refused", port)
		retry := 0
		url := fmt.Sprintf("127.0.0.1:%s", port)
		client := testutil.NewClient(url)
		_, conErr := client.CreateGetRequestByTCPConn(2)

		if conErr != nil && strings.Contains(conErr.Error(), refusedErr) && retry <= 10 {
			for retry < 3 {
				_, conErr = client.CreateGetRequestByTCPConn(2)
				time.Sleep(1 * time.Second)
				retry++
			}
		}

		if conErr == nil {
			t.Errorf("want error but no error")
		}

		if conErr != nil && strings.Contains(conErr.Error(), refusedErr) {
			t.Fatalf("want %s, but got %s", wantErr, conErr.Error())
		}

		cancel()
		resultShutdownError <- <-errShutdownCh

		result := <-resultShutdownError
		if result != nil {
			text := result.Error()
			log.Fatal(text)
		}
		t.Log("pass ReadHeaderTimeout over test")
	})

	t.Run("testSafeTimeout", func(t *testing.T) {
		// config関連
		resultShutdownError := make(chan error, 1)
		port := "3334"
		readHeaderTimeoutSecond := "2"
		os.Setenv("PORT", port)
		os.Setenv("READ_HEADER_TIMEOUT_SECOND", readHeaderTimeoutSecond)
		cfg, err := config.NewConfig()
		if err != nil {
			t.Errorf("want no error, but got %v", err)
		}
		//キャンセル可能なcontext作成
		ctx, cancel := context.WithCancel(context.Background())

		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			println("receive")
			io.WriteString(w, "<html>foo</html>")
		})

		//サーバ起動
		server := NewServer(cfg, nil)
		errShutdownCh := make(chan error, 1)
		go func() {
			// サーバが停止したら返却されるerrorをチャネルに送信
			errShutdownCh <- server.Run(ctx)
		}()

		// ReadHeaderTimeout safe
		refusedErr := fmt.Sprintf("dial tcp 127.0.0.1:%s: connect: connection refused", port)
		retry := 0
		url := fmt.Sprintf("127.0.0.1:%s", port)
		client := testutil.NewClient(url)

		rsp, conErr := client.CreateGetRequestByTCPConn(1)

		if conErr != nil && strings.Contains(conErr.Error(), refusedErr) && retry <= 10 {
			for retry < 3 {
				rsp, conErr = client.CreateGetRequestByTCPConn(1)
				time.Sleep(1 * time.Second)
				retry++
			}
		}

		if conErr != nil {
			t.Errorf("want no error, but got %v", conErr)
		}

		if !strings.Contains(rsp, "<html>foo</html>") {
			t.Errorf("want contents <html>foo</html>, but got %s", rsp)
		}

		cancel()

		resultShutdownError <- <-errShutdownCh

		result := <-resultShutdownError
		if result != nil {
			text := result.Error()
			log.Fatal(text)
		}
		t.Log("pass ReadHeaderTimeout safe test")
	})

	t.Run("testGraceful", func(t *testing.T) {
		// -------------config関連---------------
		port := "3335"
		readHeaderTimeoutSecond := "2"
		os.Setenv("PORT", port)
		os.Setenv("READ_HEADER_TIMEOUT_SECOND", readHeaderTimeoutSecond)
		cfg, err := config.NewConfig()
		if err != nil {
			t.Errorf("want no error, but got %v", err)
		}
		//---------------------------------------

		//キャンセル可能なcontext作成
		ctx, cancel := context.WithCancel(context.Background())

		// リクエストを受け取るまで処理をブロックするためのチャネル
		recieveReqGracefulTest := make(chan struct{})

		refusedErr := fmt.Sprintf("dial tcp 127.0.0.1:%s: connect: connection refused", port)
		resultShutdownError := make(chan error, 1)

		// -------------サーバ関連---------------
		http.HandleFunc("/graceful/test", func(w http.ResponseWriter, r *http.Request) {
			//リクエストを受け取ったらチャネルに通知しsleep中にキャンセルを入れる
			close(recieveReqGracefulTest)
			time.Sleep(3 * time.Second)
			io.WriteString(w, "<html>graceful</html>")
		})

		//サーバ起動
		server := NewServer(cfg, nil)
		errShutdownCh := make(chan error, 1)

		go func() {
			// サーバが停止したら返却されるerrorをチャネルに送信
			errShutdownCh <- server.Run(ctx)
		}()
		//---------------------------------------

		// httpレスポンスの結果
		isErrRsp := make(chan error, 1)

		// gracefulシャットダウンの動作を確認するため通信処理は別チャネル
		go func() {
			retry := 0
			_, err := http.Get(fmt.Sprintf("http://localhost:%s/graceful/test", port))
			if err != nil && strings.Contains(err.Error(), refusedErr) && retry <= 10 {
				for retry < 3 {
					_, err = http.Get(fmt.Sprintf("http://localhost:%s/graceful/test", port))
				}
			}
			if err != nil {
				isErrRsp <- err
			} else {
				close(isErrRsp)
			}
		}()

		// /graceful/testにリクエストが届くまでブロック
		<-recieveReqGracefulTest
		// sleep中にシャットダウンする
		cancel()

		// レスポンス内容が入るまでブロック
		resultRsp2 := <-isErrRsp
		if resultRsp2 != nil {
			t.Errorf("want no error response, bu got: %v", resultRsp2)
		}

		//シャットダウン待機
		resultShutdownError <- <-errShutdownCh
		result := <-resultShutdownError
		if result != nil {
			text := result.Error()
			log.Fatal(text)
		}
		t.Log("pass testGraceful safe test")
	})
}
