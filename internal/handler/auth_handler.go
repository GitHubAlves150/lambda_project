package handler

import (
	"encoding/json"
	service "lambda_tracker/internal/domain/services"
	"net/http"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// ============================================
// 1. RegistrarCliente - POST /auth/register
// ============================================
func (h *AuthHandler) RegistrarCliente(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Nome  string `json:"nome"`
		Email string `json:"email"`
		Senha string `json:"senha"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Requisição inválida", http.StatusBadRequest)
		return
	}

	cliente, err := h.authService.RegistrarCliente(r.Context(), req.Nome, req.Email, req.Senha)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(cliente)
}

// ============================================
// 2. GerarAPIKey - POST /auth/keys
// ============================================
func (h *AuthHandler) GerarAPIKey(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ClienteID string `json:"cliente_id"`
		Descricao string `json:"descricao"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Requisição inválida", http.StatusBadRequest)
		return
	}

	apiKey, err := h.authService.GerarAPIKey(r.Context(), req.ClienteID, req.Descricao)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(apiKey)
}

// ============================================
// 3. ValidarAPIKey - POST /auth/validate
// ============================================
func (h *AuthHandler) ValidarAPIKey(w http.ResponseWriter, r *http.Request) {
	var req struct {
		APIKey string `json:"api_key"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Requisição inválida", http.StatusBadRequest)
		return
	}

	cliente, err := h.authService.ValidarAPIKey(r.Context(), req.APIKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(cliente)
}

// ============================================
// 4. RevogarAPIKey - DELETE /auth/keys/{id}
// ============================================
func (h *AuthHandler) RevogarAPIKey(w http.ResponseWriter, r *http.Request) {
	// TODO: Implementar
	http.Error(w, "Não implementado", http.StatusNotImplemented)
}

// ============================================
// 5. ListarAPIKeys - GET /auth/keys
// ============================================
func (h *AuthHandler) ListarAPIKeys(w http.ResponseWriter, r *http.Request) {
	// TODO: Implementar
	http.Error(w, "Não implementado", http.StatusNotImplemented)
}
