package middleware

import (
	"context"
	"net/http"
	"strconv"
	"strings"
)

// contextKey はコンテキストキーの型
type contextKey string

const (
	// UserIDKey はコンテキストにユーザーIDを保存するためのキー
	UserIDKey contextKey = "user_id"
)

// GetUserID コンテキストからユーザーIDを取得
func GetUserID(ctx context.Context) *int64 {
	userID, ok := ctx.Value(UserIDKey).(int64)
	if !ok {
		return nil
	}
	return &userID
}

// AuthMiddleware アクセストークンを検証し、ユーザーIDをコンテキストに設定
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "認証が必要です", http.StatusUnauthorized)
			return
		}

		// Bearerトークンを抽出
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			http.Error(w, "無効な認証トークンです", http.StatusUnauthorized)
			return
		}

		// TODO: JWTトークンを検証してユーザーIDを取得
		// 現時点では、トークンからユーザーIDを取得するロジックを実装する必要があります
		// 例: JWTをデコードしてユーザーIDを取得
		userID, err := extractUserIDFromToken(token)
		if err != nil {
			http.Error(w, "認証トークンが無効です", http.StatusUnauthorized)
			return
		}

		// コンテキストにユーザーIDを設定
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

// OptionalAuthMiddleware 認証が任意のミドルウェア（トークンがある場合のみユーザーIDを設定）
func OptionalAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			token := strings.TrimPrefix(authHeader, "Bearer ")
			if token != "" {
				userID, err := extractUserIDFromToken(token)
				if err == nil {
					// コンテキストにユーザーIDを設定
					ctx := context.WithValue(r.Context(), UserIDKey, userID)
					r = r.WithContext(ctx)
				}
			}
		}

		next.ServeHTTP(w, r)
	})
}

// extractUserIDFromToken トークンからユーザーIDを抽出
// TODO: 実際のJWT検証ロジックに置き換える必要があります
func extractUserIDFromToken(token string) (int64, error) {
	// TODO: JWTトークンを検証してユーザーIDを取得
	// 暫定的に、トークンが数値の場合はそのままユーザーIDとして扱う
	// 実際の実装では、JWTをデコードしてユーザーIDを取得する必要があります
	userID, err := strconv.ParseInt(token, 10, 64)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

// CORSMiddleware CORSヘッダーを追加
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// LoggingMiddleware HTTPリクエストをログ記録
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: ロギングを実装
		next.ServeHTTP(w, r)
	})
}
