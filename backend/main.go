package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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

	// main.go の適当な場所（e.POST("/users", ...) の下あたり）に追記

// 日報作成の構造体（リクエストを受け取る用）
type CreateReportRequest struct {
	UserID        string `json:"user_id"`
	Title         string `json:"title"`
	Content       string `json:"content"`
	LearningNotes string `json:"learning_notes"`
	IsHabitDone   bool   `json:"is_habit_done"`
	IsPublic      bool   `json:"is_public"`
	SubmittedDate string `json:"submitted_date"` // "2026-01-31" 形式
}

// 日報作成エンドポイント
e.POST("/reports", func(c echo.Context) error {
	var req CreateReportRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// 文字列のIDをUUIDに変換（google/uuidライブラリを使うと便利です）
	userID, _ := uuid.Parse(req.UserID)
	// 日付のパース
	subDate, _ := time.Parse("2006-01-02", req.SubmittedDate)
	// 共有用URLのトークンを適当に生成（本来はもっと複雑に）
	publicToken := uuid.New().String()

	report, err := queries.CreateReport(ctx, db.CreateReportParams{
		UserID:        userID,
		Title:         req.Title,
		Content:       sql.NullString{String: req.Content, Valid: true},
		LearningNotes: sql.NullString{String: req.LearningNotes, Valid: true},
		IsHabitDone:   sql.NullBool{Bool: req.IsHabitDone, Valid: true},
		IsPublic:      sql.NullBool{Bool: req.IsPublic, Valid: true},
		PublicToken:   publicToken,
		SubmittedDate: subDate,
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, report)
})
	// サーバー起動 (8080ポートで待ち受け)
	e.Logger.Fatal(e.Start(":8080"))
}