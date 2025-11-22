package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"example.com/backend/internal/application/user/service"
)

// UserHandler ユーザーエンドポイントのHTTPリクエストを処理
type UserHandler struct {
	service *service.Service
}

// NewUserHandler 新しいユーザーハンドラーを作成
func NewUserHandler(svc *service.Service) *UserHandler {
	return &UserHandler{
		service: svc,
	}
}

// Register POST /users を処理
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "メソッドが許可されていません", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// TODO: パスワードをハッシュ化
	// user, err := h.service.RegisterUser(r.Context(), req.Email, hashedPassword, req.Name)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "ユーザーを作成しました"})
}

// GetUser GET /users/:id を処理
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "メソッドが許可されていません", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "無効なID", http.StatusBadRequest)
		return
	}

	// TODO: サービスからユーザーを取得
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int64{"id": id})
}

