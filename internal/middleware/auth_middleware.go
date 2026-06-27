package middleware

import (
    "context"
    "net/http"
    "strings"
    "lambda_tracker/internal/domain/services"
)

// ContextKey tipo próprio para evitar conflito de chaves no contexto
type ContextKey string

// constante tipada — chave para guardar o userID no contexto
const UserIDKey ContextKey = "user_id"

// AuthMiddleware retorna uma função que envolve o handler
// é executado ANTES de qualquer rota protegida
func AuthMiddleware(authService service.AuthService) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

            // 1. pega a key do header X-API-Key
            apiKey := r.Header.Get("X-API-Key")
            if apiKey == "" {
                // interrompe aqui — não chega na rota
                http.Error(w, "API Key não fornecida", http.StatusUnauthorized)
                return
            }

            // 2. remove "Bearer " se vier com prefixo
            apiKey = strings.TrimPrefix(apiKey, "Bearer ")

            // 3. valida a key no service
            cliente, err := authService.ValidarAPIKey(r.Context(), apiKey)
            if err != nil {
                http.Error(w, "API Key inválida", http.StatusUnauthorized)
                return
            }

            // 4. injeta o ID do cliente no contexto da requisição
            // próximos handlers podem ler com r.Context().Value(UserIDKey)
            ctx := context.WithValue(r.Context(), UserIDKey, cliente.ID)

            // 5. passa para o próximo handler com o contexto enriquecido
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}