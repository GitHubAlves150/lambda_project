package entity

import (
	"crypto/rand"  // geração de bytes aleatórios criptograficamente seguros
	"encoding/hex" // converte bytes para string hexadecimal legível
	"time"         // Adicione esta linha no bloco de imports

	"github.com/google/uuid"
)

type APIKey struct {
	ID        string     `json:"id"`
	ClienteID string     `json:"cliente_id"`
	KeyHash   string     `json:"-"`   // hash salvo no banco — nunca exposto
	Key       string     `json:"key"` // chave real — só aparece UMA VEZ no retorno
	Descricao string     `json:"descricao"`
	Ativo     bool       `json:"ativo"`
	CriadoEm  time.Time  `json:"criado_em"`
	UltimoUso *time.Time `json:"ultimo_uso,omitempty"` // ponteiro pois pode ser nulo
	// omitempty = omite do JSON se for nil
}

// NovaAPIKey gera uma API Key criptograficamente segura
func NovaAPIKey(clienteID, descricao string) (*APIKey, error) {
	// cria um slice de 32 bytes vazios
	bytes := make([]byte, 32)

	// preenche com bytes aleatórios seguros do sistema operacional
	// é diferente de math/rand que é previsível
	if _, err := rand.Read(bytes); err != nil {
		return nil, err
	}

	// converte os 32 bytes em 64 caracteres hexadecimais
	// ex: "a3f9b2c1d4e5..." — a chave que o cliente vai usar
	key := hex.EncodeToString(bytes)

	return &APIKey{
		ID:        uuid.New().String(),
		ClienteID: clienteID,
		Key:       key,          // chave real — retornada só aqui, nunca mais
		KeyHash:   hashKey(key), // hash que fica salvo no banco
		Descricao: descricao,
		Ativo:     true,
		CriadoEm:  time.Now(),
	}, nil
}

// hashKey gera o hash da chave
// ⚠️ no código está simplificado
// em produção usa bcrypt.GenerateFromPassword
func hashKey(key string) string {
	return key // ← SIMPLIFICADO — em produção NUNCA faça isso
}
