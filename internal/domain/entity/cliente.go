package entity

import (
    "time"
    "github.com/google/uuid" // biblioteca para gerar IDs únicos universais
)

// Cliente representa a entidade de negócio
// é o molde dos dados de um cliente no sistema
// não sabe nada de banco, HTTP ou framework
type Cliente struct {
    ID        string    `json:"id"`        // exposto no JSON
    Nome      string    `json:"nome"`      // exposto no JSON
    Email     string    `json:"email"`     // exposto no JSON
    SenhaHash string    `json:"-"`         // ← o "-" ESCONDE este campo do JSON
                                           // nunca vai aparecer na resposta HTTP
                                           // segurança: ninguém vê o hash da senha
    CriadoEm  time.Time `json:"criado_em"` // exposto no JSON
}

// NovoCliente é o construtor da entidade
// garante que todo cliente nasce com ID gerado e dados completos
// ninguém cria um Cliente sem passar por aqui
func NovoCliente(nome, email, senhaHash string) *Cliente {
    return &Cliente{
        ID:        uuid.New().String(), // gera UUID único ex: "a1b2c3-..."
        Nome:      nome,
        Email:     email,
        SenhaHash: senhaHash,
        CriadoEm:  time.Now(),         // momento exato da criação
    }
}