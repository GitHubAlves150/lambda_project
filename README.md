
## 💻 lambda_project

&nbsp;&nbsp;&nbsp;&nbsp;Este projeto Backend em Golang é usado no projeto lilygo-vehicle-tracking disponívle em https://github.com/GitHubAlves150/lilygo-vehicle-tracking.git. Para cada branch deste repositório é um avanço gradual da minha dedicação ao Golang. A branch master é um código simples para validar um MVP e as demais branchs são estudos e práticas profissionais que abordam conhecimentos de um profissional pleno, trabalhando com banco de dados, API, microserviços, golang como linguagem principal e cleancode.

&nbsp;&nbsp;&nbsp;&nbsp;O projeto lilygo-tracking tem como alicerce em c++ para gerar os dados de GPS para então ser implementado a linguagem Go + as principais arquiteturas de software, framework e cleancode assim como S.O.L.I.D, microserviços e API's, context, gin, DB postgres, autenticação, echo entre outros framworks que envolve um sisitema backend robusto e escalável. Por escolha própria decidi gerar dados GPS de um dispositivo físico.

 &nbsp;&nbsp;&nbsp;&nbsp;Havia um professor meu na graduação de Engenharia Eletrônica no IFSC-FLN que, disse algo que me ajudou muito a destrinxar minha forma de entender assuntos complexos e construir sistemas eficases - "Pense como os romanos, Divida para conquistar". Todo sistema, seja ele eletrônico, firmware e/ou web é preciso dividir em partes para unir em algo maior.

 🔍&nbsp;&nbsp;&nbsp;&nbsp;Atenção! Todo este trabalho foi somente possível graças a ajuda e consulta da I.A Deepseek, perplexity, gemini, sites e alguns vídeos no youtube. A maioria dos estudos e laboratório de testes foi baseados em códigos de exemplo gerados e explicados pela I.A e, no final juntado os exercícios para criar o sistema backend completo, Os demais detalhes como organização, comentários preparação de estrutura e aplicação dos conceitos e técnicasde programação é dedicação do autor deste projeto.



 Toda a pesquisa pode ser conferida neste link https://chat.deepseek.com/share/gcw3kbzsvbnwdrrmfe .

### Indíce de Branchs:
 - main: Simples MVP para rodar no lambda AWS   
 - 1.Fase1/APIs_DataBase: Migração para postrgre - Go + Post
 - 2.Fase1.1/APIs_AP_KEY: APIKEY - Autenticação


### 🔐 Fase 1 - Tarefa 1: Conexão com o banco de dados postgresSQL [✅]
- git branch Fase1/APIs_DataBase
### 🔐 Fase 1 - TAREFA 2: Autenticação (API Key) - Explicação SIMPLES [você está aqui]
- git branch Fase1.1/APIs_AP_KEY

### Introdução
Qual problema que ela resolve !
```bash
❌ QUALQUER PESSOA pode enviar dados para sua API
❌ Um hacker pode enviar dados falsos
❌ Você não sabe quem está enviando
Com API Key:
✅ Só quem tem a chave pode enviar dados
✅ Você sabe qual veículo está enviando
✅ Dados são autenticados
``` 
## Instalar dependências

- go get github.com/go-chi/chi/v5
- go get github.com/google/uuid
- go get github.com/go-chi/chi
- go get github.com/go-chi/chi/middleware
- go get github.com/go-chi/chi/v5
- go get github.com/go-chi/chi/v5/middleware

## O código desta fase foi feito inteiramente na main.go para fins de bons entendimentos. Mas, após validar conseitos as práticas de cleancode e S.O.L.I.D foram estavelicidas, divididas em cada pasta com sua responsabilidades.
 Segue o código completo na main.go

 ```bash
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

 ```

 ### 🧪 Como Testar
 - 1. Registrar um cliente 
 ```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"nome":"João","email":"joao@email.com"}'

 # Resposta:
 # {"id":"abc123","nome":"João","email":"joao@email.com"}

 ```

2. Criar uma API Key para o cliente
```bash
curl -X POST http://localhost:8080/auth/keys \
  -H "Content-Type: application/json" \
  -d '{"cliente_id":"abc123","descricao":"ESP32-001"}'

# Resposta:
# {"id":"xyz456","cliente_id":"abc123","key":"f7e8a9b0c1d2e3f4","descricao":"ESP32-001","ativo":true}
#                        ↑
#                    GUARDE ESTA CHAVE!

``` 
3. Validar a chave (teste)

```bash
curl -X POST http://localhost:8080/auth/validate \
  -H "Content-Type: application/json" \
  -d '{"key":"f7e8a9b0c1d2e3f4"}'

# Resposta:
# {"id":"abc123","nome":"João","email":"joao@email.com"}
``` 
4. Acessar rota protegida COM a chave 

```bash
curl -X GET http://localhost:8080/api/dados-secretos \
  -H "X-API-Key: f7e8a9b0c1d2e3f4"

# Resposta:
# {"mensagem":"🎉 Você acessou dados secretos!","cliente":"João"}

``` 
5. Acessar rota protegida SEM a chave (bloqueado)

``` bash
curl -X GET http://localhost:8080/api/dados-secretos

# Resposta:
# API Key não fornecida (401 Unauthorized)

```

### 📊 Resumo da Arquitetura

``` bash
┌─────────────────────────────────────────────────────────────────────┐
│                         CLIENTE (ESP32)                             │
│                                                                     │
│  1. Obtém a chave: "f7e8a9b0c1d2e3f4"                               │
│  2. Envia requisição com Header: X-API-Key: f7e8a9b0c1d2e3f4        │
└─────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────┐
│                         MIDDLEWARE (Segurança)                      │
│                                                                     │
│  1. Pega a chave do header                                          │
│  2. Valida no "banco"                                               │
│  3. Se válido: passa o cliente no contexto                          │
│  4. Se inválido: retorna 401 Unauthorized                           │
└─────────────────────────────────────────────────────────────────────┘
                                    │
                    ┌───────────────┴───────────────┐
                    ▼                               ▼
          ┌─────────────────┐             ┌─────────────────┐
          │  ROTA PROTEGIDA │             │ 401 UNAUTHORIZED│
          │ (dados secretos)│             │   (bloqueado)   │
          └─────────────────┘             └─────────────────┘
``` 

```bash
✅ O que você aprendeu
Conceito	            Explicação
Autenticação      	Verificar se o cliente é quem diz ser
API Key	            Uma "senha" para acessar a API
Middleware	         Código que executa ANTES das rotas
Header	            Onde a chave é enviada (X-API-Key)
Contexto	            Como passar informações entre handlers
Rota protegida	      Rota que só pode ser acessada com autenticação

``` 


### Como funciona na prática.
```bash
1. Você (administrador) cria uma chave para o veículo
   ─────────────────────────────────────────────────────
   POST /auth/keys
   Body: {"cliente_id": "123", "descricao": "ESP32-001"}
   
   Resposta: {"key": "abc123xyz456"}

2. Você coloca essa chave no código do ESP32
   ─────────────────────────────────────────────────────
   // No código do ESP32
   http.setHeader("X-API-Key", "abc123xyz456")

3. O ESP32 envia dados com a chave
   ─────────────────────────────────────────────────────
   POST /api/telemetria
   Header: X-API-Key: abc123xyz456
   Body: {"latitude": -23.55, "longitude": -46.63}

4. O servidor valida a chave antes de aceitar os dados
   ─────────────────────────────────────────────────────
   ✅ Chave válida → aceita os dados
   ❌ Chave inválida → retorna 401 Unauthorized

```


