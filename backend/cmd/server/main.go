package main

import (
	"context"
	"log"
	stdhttp "net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	database "example.com/backend/internal/infrastructure/database/turso"
	"example.com/backend/internal/presentation/http"
)

func main() {
	// データベースの初期化
	db, err := database.NewClient()
	if err != nil {
		log.Fatalf("データベースの初期化に失敗しました: %v", err)
	}
	defer db.Close()

	// HTTPサーバーの初期化
	router := http.NewRouter(db.DB())
	
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &stdhttp.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// サーバーをgoroutineで起動
	go func() {
		log.Printf("サーバーを %s で起動中", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != stdhttp.ErrServerClosed {
			log.Fatalf("サーバーエラー: %v", err)
		}
	}()

	// グレースフルシャットダウン
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("サーバーをシャットダウン中...")
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("サーバーの強制終了: %v", err)
	}

	log.Println("サーバーを終了しました")
}
