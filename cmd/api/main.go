// 2026-06-08
package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/google/uuid"
)

// =================================================
// 1. Estrutura de dados
// =================================================
// Cliente representa um usuário do sisitema
// Em um caso real, isso seria salvo no banco
type Cliente struct {
	ID    string `json:"id"`
	Nome  string `json:"nome"`
	Email string `json:"email"`
}

// APIkey representa a chava de acesso
// Cada cliente pode ter várias chaves(ex: uma por dispositivos)
type APIkey struct {
	ID        string `json:"id"`
	ClienteID string `json:"cliente_id"`
	Key       string `json:"key"`       //a chave em texto puro
	Descricao string `json:"descricao"` //ex: "ESP32-001"
	Ativo     bool   `json:"ativo"`
}

//=================================================
// 2. Banco de dados em memória para simulação
//=================================================

// em produção, isso seria um banco de dados real(PostgreSQL, Mysql, etc...)
var (
	clientes = make(map[string]Cliente, 10) //ID - cliente
	apikeys  = make(map[string]APIkey)      //Key - APIKey (para busca rápida)
)

//=================================================
// 3. Função de autenticação
//=================================================

// GerarAPIKey cria uma chave aleatória - usa crypto/rand para segurança (Não use èmath/rand)
func GerarAPIKey() string {
	//16 bytes = 128 bits de entrologia(seguro)
	bytes := make([]byte, 16)
	//crypto/rand é criptograficamente seguro
	//math/rand NÂO é seguro para chaves por que é previsível
	_, err := rand.Read(bytes)
	if err != nil { //=================================================
		panic("😱 Erro ao gerar a chave" + err.Error())
	}
	//hex.EncodeToString() converte bytes para texto string leve em formato hexadecimal
	//Ex: [0xAB, 0xAbe] para "absnfd433kkff09"
	return hex.EncodeToString(bytes)
}

// --Registra cliente cria um novo cliente
func RegistraCliente(nome, email string) Cliente {
	cliente := Cliente{
		ID:    uuid.New().String(), //uuid único
		Nome:  nome,
		Email: email,
	}

	clientes[cliente.ID] = cliente
	return cliente
}

// --CriarAPIKey cria uma nova chave para um cliente
func CriarAPIKey(clienteID, descricao string) APIkey {
	key := APIkey{
		ID:        uuid.New().String(),
		ClienteID: clienteID,
		Key:       GerarAPIKey(), //gera chave aleatória
		Descricao: descricao,
		Ativo:     true,
	}

	//garda a chave
	apikeys[key.Key] = key
	return key
}

// --ValidarAPIkey verifica se uma chave é válida
func ValidarAPIKey(key string) (*Cliente, error) {
	//1. Busca a chave no bunco
	apikey, existe := apikeys[key]
	if !existe {
		return nil, fmt.Errorf("Chave não encontrada")
	}
	//2. Verifica se está ativa
	if !apikey.Ativo {
		return nil, fmt.Errorf("Chave inativa")
	}
	//3. Busca o cliente dono da chave
	cliente, existe := clientes[apikey.ClienteID]
	if !existe {
		return nil, fmt.Errorf("Cliente não encontrado")
	}
	return &cliente, nil
}

//=================================================
// 4. handlers (o que processa as requisições)
//=================================================

// Handler para registrar cliente
// POST /auth/register
func HandlerResgistrarCliente(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Nome  string `json:"nome"`
		Email string `json:"email"`
	}

	//Decodifica JSON
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON Inválido", http.StatusBadRequest)
		return
	}

	//Registra o cliente
	cliente := RegistraCliente(req.Nome, req.Email)

	//Responde com o cliente criado
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(cliente)
}

// Handler para criar APIKEY
// POST /auth/keys
func handleCriarAPIKey(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ClienteID string `json:"cliente_id"`
		Descricao string `json:"descricao"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "JSON Inválido", http.StatusBadRequest)
	}

	//Verifica se o cliente existe
	if _, existe := clientes[req.ClienteID]; !existe {
		http.Error(w, "Cliente não encontrado", http.StatusBadRequest)
		return
	}

	//Cria a chave
	apikey := CriarAPIKey(req.ClienteID, req.Descricao)

	//Responde com a chave(GUARDE BEM ESTA CHAVE!)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(apikey)

}

// handler para validar a chave
func HandlerValidarAPIKey(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Key string `json:"key"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "JSON Inválido", http.StatusBadRequest)
		return
	}

	cliente, err := ValidarAPIKey(req.Key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cliente)

}

//=================================================
// 5. Rota que precisa de autenticação
//=================================================

//Handler para rtoa protegida
//GET /api/dados-secretos

func HandlerDadosSecretos(w http.ResponseWriter, r *http.Request) {
	//o middleware colocou o cliente no contexto!
	cliente, ok := r.Context().Value("cliente").(*Cliente)
	if !ok {
		http.Error(w, "Não autorizado", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"menssagem": "🎉 Você acessou dados secretos!!!",
		"cliente":   cliente.Nome,
	})
}

// =================================================
// 6. Middler de Autenticação (o coração)
// =================================================
// 🔒 AuthMiddleWare - O segurança da aplicação
func AuthMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// ==========================================
		// PASSO 1: Pega a chave do header
		// ==========================================

		//O esp32 vai enviar a chave no header X-API-key
		//EX: X-API-KEY : abc234rhdh67
		apikey := r.Header.Get("X-API-KEY")

		//Se não tiver chave, bloqueia imediatamente
		if apikey == "" {
			http.Error(w, "API key não fornecida", http.StatusUnauthorized)
			return
		}

		// ==========================================
		// PASSO 2: Remove prefixos (opcional)
		// ==========================================

		//Alguns sisitemas usam "Bearer " antes da chave
		//Ex; "Bearer abc345" para "absdh432"
		apikey = strings.TrimPrefix(apikey, "Bearer ")

		// ==========================================
		// PASSO 3: Valida a chave
		// ==========================================
		cliente, err := ValidarAPIKey(apikey)
		if err != nil {
			http.Error(w, "API key inválida"+err.Error(), http.StatusUnauthorized)
			return
		}

		// ==========================================
		// PASSO 4: Adiciona o cliente no contexto
		// ==========================================
		//O contexto permite passar informações entre handlers
		//Agora qualquer handler depois pode pegar o cliente
		ctx := context.WithValue(r.Context(), "cliente", cliente)

		// ==========================================
		// PASSO 5: Deixa a requisição passar
		// ==========================================
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func main() {

	//Cria o roteador chi
	r := chi.NewRouter()

	// ==========================================
	// Middlewares GLOBAIS (executam para TODAS as rotas)
	// ==========================================

	//Logger: Registra todas as requisições
	r.Use(middleware.Logger)

	//Recover: Recupera de panics (evita que o servidor caia)
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", HandlerResgistrarCliente)
		r.Post("/keys", handleCriarAPIKey)
		r.Post("/validate", HandlerValidarAPIKey)
		r.Get("/api/dados-secretos", HandlerDadosSecretos)
			  
	})

	// ==========================================
	// Rotas PROTEGIDAS (com autenticação)
	// ==========================================

	r.Group(func(r chi.Router) {
		//🔒Aplica middlerware de autenticação
		//Todas as rotas dentro deste grupo serão protegidasAPI key não fornecida
		//r.Post("/api/telemetria, HandlerTelemetria")

	})

	// ==========================================
	// Inicia o servidor
	// ==========================================
	fmt.Println("🚀 Servidor rodando em http://localhost:8080")
	fmt.Println("📚 Endpoints:")
	fmt.Println("  POST /auth/register   - Registrar cliente")
	fmt.Println("  POST /auth/keys       - Gerar API Key")
	fmt.Println("  POST /auth/validate   - Validar API Key")
	fmt.Println("  GET  /api/dados-secretos - Rota protegida")

	http.ListenAndServe(":8080", r)
}
