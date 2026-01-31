package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/yourname/diary-backend/db" // 自分のmodule名に合わせてね
)

func main() {
	ctx := context.Background()
	// Docker Composeで設定した接続情報
	connStr := "postgresql://user:password@localhost:5432/diary_db?sslmode=disable"

	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		log.Fatalf("DB接続失敗: %v", err)
	}
	defer conn.Close(ctx)

	// sqlcが生成したクエリ実行用の構造体
	queries := db.New(conn)

	// テスト：ユーザーを作ってみる
	user, err := queries.CreateUser(ctx, db.CreateUserParams{
		GoogleId: "google-test-id-123",
		Email:    "test@example.com",
		Name:     "テストユーザー",
	})

	if err != nil {
		log.Fatalf("ユーザー作成失敗: %v", err)
	}

	fmt.Printf("ユーザー作成成功！ ID: %s, Name: %s\n", user.ID, user.Name)
}