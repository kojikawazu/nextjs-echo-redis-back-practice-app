package main

import (
	"backend/handlers"
	"backend/supabase"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// 環境変数の読み込み
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	// Supabaseクライアントの初期化
	err = supabase.InitSupabase()
	if err != nil {
		log.Fatalf("Supabase initialization failed: %v", err)
	}

	// テストクエリの実行
	err = supabase.TestQuery()
	if err != nil {
		log.Fatalf("Test query failed: %v", err)
	}

	e := echo.New()

	// ミドルウェアの設定
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// APIエンドポイントの設定
	e.GET("/api/todos", handlers.GetTodos)
	e.POST("/api/todos", handlers.CreateTodo)
	e.PUT("/api/todos/:id", handlers.UpdateTodo)
	e.DELETE("/api/todos/:id", handlers.DeleteTodo)

	// シグナルハンドラーの設定
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		log.Println("Shutting down server...")

		// Echoサーバーのシャットダウン
		if err := e.Close(); err != nil {
			log.Fatalf("Echo shutdown failed: %v", err)
		}

		// Supabaseコネクションプールのクローズ
		supabase.ClosePool()
	}()

	// サーバーの起動
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := e.Start(":" + port); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Echo server failed: %v", err)
	}
}
