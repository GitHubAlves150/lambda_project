package repository

import (
	"context"
	"lambda_tracker/internal/domain/entity"
)

// ClienteRepository define o CONTRATO do repositório de clientes
// o service depende desta interface — não da implementação PostgreSQL
// isso permite trocar o banco sem mudar o service
type ClienteRepository interface {
	Salvar(ctx context.Context, cliente *entity.Cliente) error
	BuscarPorEmail(ctx context.Context, email string) (*entity.Cliente, error)
	BuscarPorID(ctx context.Context, id string) (*entity.Cliente, error)
}

// APIKeyRepository define o CONTRATO do repositório de api keys
type APIKeyRepository interface {
	Salvar(ctx context.Context, apiKey *entity.APIKey) error
	BuscarPorHash(ctx context.Context, keyHash string) (*entity.APIKey, error)
	BuscarPorClienteID(ctx context.Context, clienteID string) ([]entity.APIKey, error)
	AtualizarUltimoUso(ctx context.Context, apiKeyID string) error
	Revogar(ctx context.Context, apiKeyID string) error
}
