package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rii-3112/diary-report/backend/db"
)

func main() {
	e := echo.New()

	// ログを表示してくれる便利な設定
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// 1. DB接続設定
	ctx := context.Background()
	connStr := os.Getenv("DB_SOURCE")
	if connStr == "" {
		connStr = "postgresql://user:password@localhost:5432/diary_db?sslmode=disable"
	}

	dbConn, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatalf("DB接続失敗: %v", err)
	}
	defer dbConn.Close()

	queries := db.New(dbConn)

	// 2. エンドポイント（疎通確認用）
	// ブラウザで http://localhost:8080/hello を開くと見れます
	e.GET("/hello", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"message": "Hello! Diary Backend is working!",
		})
	})

	// 3. テスト用のユーザー作成API
	// POSTリクエストを送ると、DBにユーザーを保存します
	e.POST("/users", func(c echo.Context) error {
		user, err := queries.CreateUser(ctx, db.CreateUserParams{
			GoogleID: "test-id-" + fmt.Sprint(os.Getpid()),
			Email:    "test@example.com",
			Name:     sql.NullString{String: "Gopher", Valid: true},
		})
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusCreated, user)
	})

	// サーバー起動 (8080ポートで待ち受け)
	e.Logger.Fatal(e.Start(":8080"))
}