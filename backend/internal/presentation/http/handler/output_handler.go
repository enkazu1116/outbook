package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"example.com/backend/internal/application/output/service"
	"example.com/backend/internal/shared/dto"
	"example.com/backend/internal/shared/middleware"
)

// OutputHandler アウトプットエンドポイントのHTTPリクエストを処理
type OutputHandler struct {
	service *service.Service
}

// NewOutputHandler 新しいアウトプットハンドラーを作成
func NewOutputHandler(svc *service.Service) *OutputHandler {
	return &OutputHandler{
		service: svc,
	}
}

// CreateOutput POST /api/outputs を処理
func (h *OutputHandler) CreateOutput(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "メソッドが許可されていません", http.StatusMethodNotAllowed)
		return
	}

	// コンテキストからユーザーIDを取得
	userID := middleware.GetUserID(r.Context())
	if userID == nil {
		h.sendError(w, "認証が必要です", http.StatusUnauthorized)
		return
	}

	var req dto.CreateOutputRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// DTOからエンティティを作成（認証されたユーザーIDを使用）
	output := req.ToEntity(*userID)

	if err := h.service.CreateOutput(r.Context(), output); err != nil {
		h.sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.sendSuccess(w, output, http.StatusCreated)
}

// GetOutput GET /api/outputs/:id を処理
func (h *OutputHandler) GetOutput(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "メソッドが許可されていません", http.StatusMethodNotAllowed)
		return
	}

	// ルーターが/apiプレフィックスを削除するので、/outputs/123の形式
	idStr := r.URL.Path[len("/outputs/"):]
	if idStr == "" {
		h.sendError(w, "無効なID", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.sendError(w, "無効なID", http.StatusBadRequest)
		return
	}

	// コンテキストからユーザーIDを取得（認証が任意の場合はnilの可能性がある）
	userID := middleware.GetUserID(r.Context())

	output, err := h.service.GetOutput(r.Context(), id, userID)
	if err != nil {
		if err == service.ErrOutputNotFound || err == service.ErrOutputUnauthorized {
			h.sendError(w, err.Error(), http.StatusNotFound)
		} else {
			h.sendError(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	h.sendSuccess(w, output, http.StatusOK)
}

// GetOutputsByUserID GET /api/users/:user_id/outputs を処理
func (h *OutputHandler) GetOutputsByUserID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "メソッドが許可されていません", http.StatusMethodNotAllowed)
		return
	}

	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		h.sendError(w, "user_idが必要です", http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		h.sendError(w, "無効なuser_id", http.StatusBadRequest)
		return
	}

	outputs, err := h.service.GetOutputsByUserID(r.Context(), userID, nil)
	if err != nil {
		h.sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.sendSuccess(w, outputs, http.StatusOK)
}

// GetOutputsByBookID GET /api/books/:book_id/outputs を処理
func (h *OutputHandler) GetOutputsByBookID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "メソッドが許可されていません", http.StatusMethodNotAllowed)
		return
	}

	bookIDStr := r.URL.Query().Get("book_id")
	if bookIDStr == "" {
		h.sendError(w, "book_idが必要です", http.StatusBadRequest)
		return
	}

	bookID, err := strconv.ParseInt(bookIDStr, 10, 64)
	if err != nil {
		h.sendError(w, "無効なbook_id", http.StatusBadRequest)
		return
	}

	outputs, err := h.service.GetOutputsByBookID(r.Context(), bookID, nil)
	if err != nil {
		h.sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.sendSuccess(w, outputs, http.StatusOK)
}

// UpdateOutput PUT /api/outputs/:id を処理
func (h *OutputHandler) UpdateOutput(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "メソッドが許可されていません", http.StatusMethodNotAllowed)
		return
	}

	// ルーターが/apiプレフィックスを削除するので、/outputs/123の形式
	idStr := r.URL.Path[len("/outputs/"):]
	if idStr == "" {
		h.sendError(w, "無効なID", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.sendError(w, "無効なID", http.StatusBadRequest)
		return
	}

	var req dto.UpdateOutputRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// コンテキストからユーザーIDを取得
	userID := middleware.GetUserID(r.Context())
	if userID == nil {
		h.sendError(w, "認証が必要です", http.StatusUnauthorized)
		return
	}

	// 既存のアウトプットを取得
	existing, err := h.service.GetOutput(r.Context(), id, userID)
	if err != nil {
		if err == service.ErrOutputNotFound || err == service.ErrOutputUnauthorized {
			h.sendError(w, err.Error(), http.StatusNotFound)
		} else {
			h.sendError(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// DTOのフィールドを既存のエンティティに適用
	req.ApplyToEntity(existing)

	if err := h.service.UpdateOutput(r.Context(), existing, *userID); err != nil {
		if err == service.ErrOutputNotFound || err == service.ErrOutputUnauthorized {
			h.sendError(w, err.Error(), http.StatusNotFound)
		} else {
			h.sendError(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	h.sendSuccess(w, existing, http.StatusOK)
}

// DeleteOutput DELETE /api/outputs/:id を処理
func (h *OutputHandler) DeleteOutput(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "メソッドが許可されていません", http.StatusMethodNotAllowed)
		return
	}

	// ルーターが/apiプレフィックスを削除するので、/outputs/123の形式
	idStr := r.URL.Path[len("/outputs/"):]
	if idStr == "" {
		h.sendError(w, "無効なID", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.sendError(w, "無効なID", http.StatusBadRequest)
		return
	}

	// コンテキストからユーザーIDを取得
	userID := middleware.GetUserID(r.Context())
	if userID == nil {
		h.sendError(w, "認証が必要です", http.StatusUnauthorized)
		return
	}

	if err := h.service.DeleteOutput(r.Context(), id, *userID); err != nil {
		if err == service.ErrOutputNotFound || err == service.ErrOutputUnauthorized {
			h.sendError(w, err.Error(), http.StatusNotFound)
		} else {
			h.sendError(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	h.sendSuccess(w, map[string]string{"message": "アウトプットを削除しました"}, http.StatusOK)
}

// sendSuccess 成功レスポンスを送信
func (h *OutputHandler) sendSuccess(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(dto.StandardResponse{
		Success: true,
		Data:    data,
	})
}

// sendError エラーレスポンスを送信
func (h *OutputHandler) sendError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(dto.StandardResponse{
		Success: false,
		Error:   message,
	})
}
