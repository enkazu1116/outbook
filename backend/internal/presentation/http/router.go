package http

import (
	"database/sql"
	"net/http"

	"example.com/backend/internal/application/output/service"
	database "example.com/backend/internal/infrastructure/database/turso"
	outputHandler "example.com/backend/internal/presentation/http/handler"
	"example.com/backend/internal/shared/middleware"
)

// NewRouter HTTPルートを設定
func NewRouter(db *sql.DB) http.Handler {
	// リポジトリの初期化
	outputRepo := database.NewOutputRepository(db)

	// サービスの初期化
	outputService := service.NewService(outputRepo)

	// ハンドラーの初期化
	outputHandler := outputHandler.NewOutputHandler(outputService)

	mux := http.NewServeMux()

	// ヘルスチェック
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// API ルート
	api := http.NewServeMux()

	// Outputエンドポイント
	// POST /api/outputs - アウトプットを作成（認証必須）
	api.HandleFunc("/outputs", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			outputHandler.CreateOutput(w, r)
		} else {
			http.Error(w, "メソッドが許可されていません", http.StatusMethodNotAllowed)
		}
	})

	// GET/PUT/DELETE /api/outputs/:id - アウトプットを取得/更新/削除
	// GETは認証任意、PUT/DELETEは認証必須
	api.HandleFunc("/outputs/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			outputHandler.GetOutput(w, r)
		case http.MethodPut:
			outputHandler.UpdateOutput(w, r)
		case http.MethodDelete:
			outputHandler.DeleteOutput(w, r)
		default:
			http.Error(w, "メソッドが許可されていません", http.StatusMethodNotAllowed)
		}
	})

	// GET /api/users/outputs?user_id=:id - ユーザーのアウトプット一覧（認証任意）
	api.HandleFunc("/users/outputs", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			outputHandler.GetOutputsByUserID(w, r)
		} else {
			http.Error(w, "メソッドが許可されていません", http.StatusMethodNotAllowed)
		}
	})

	// GET /api/books/outputs?book_id=:id - 書籍のアウトプット一覧（認証任意）
	api.HandleFunc("/books/outputs", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			outputHandler.GetOutputsByBookID(w, r)
		} else {
			http.Error(w, "メソッドが許可されていません", http.StatusMethodNotAllowed)
		}
	})

	// すべてのAPIエンドポイントにOptionalAuthMiddlewareを適用
	// 認証が必要なエンドポイントはハンドラー内でチェック
	apiWithAuth := middleware.OptionalAuthMiddleware(api)
	mux.Handle("/api/", http.StripPrefix("/api", apiWithAuth))

	// ミドルウェアを適用
	var handler http.Handler = mux
	handler = middleware.LoggingMiddleware(handler)
	handler = middleware.CORSMiddleware(handler)

	return handler
}
