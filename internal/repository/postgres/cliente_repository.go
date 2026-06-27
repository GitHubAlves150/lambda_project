package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"lambda_tracker/internal/domain/entity"
)

// ClienteRepositoryPostgres é a implementação concreta para PostgreSQL
type ClienteRepositoryPostgres struct {
	db *sql.DB
}

// NewClienteRepositoryPostgres é a FUNÇÃO CONSTRUTORA!
// Ela recebe a conexão com o banco e retorna uma instância do repositório.
func NewClienteRepositoryPostgres(db *sql.DB) *ClienteRepositoryPostgres {
	return &ClienteRepositoryPostgres{
		db: db,
	}
}

// Salvar insere um novo cliente no banco
func (r *ClienteRepositoryPostgres) Salvar(ctx context.Context, cliente *entity.Cliente) error {
	query := `INSERT INTO clientes (id, nome, email, senha_hash, criado_em) 
              VALUES ($1, $2, $3, $4, $5)`

	_, err := r.db.ExecContext(ctx, query,
		cliente.ID,
		cliente.Nome,
		cliente.Email,
		cliente.SenhaHash,
		cliente.CriadoEm,
	)
	if err != nil {
		return fmt.Errorf("erro ao salvar cliente: %w", err)
	}
	return nil
}

// BuscarPorEmail busca um cliente pelo email
func (r *ClienteRepositoryPostgres) BuscarPorEmail(ctx context.Context, email string) (*entity.Cliente, error) {
	query := `SELECT id, nome, email, senha_hash, criado_em 
              FROM clientes 
              WHERE email = $1`

	var cliente entity.Cliente
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&cliente.ID,
		&cliente.Nome,
		&cliente.Email,
		&cliente.SenhaHash,
		&cliente.CriadoEm,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Cliente não encontrado
		}
		return nil, fmt.Errorf("erro ao buscar cliente: %w", err)
	}
	return &cliente, nil
}

// BuscarPorID busca um cliente pelo ID
func (r *ClienteRepositoryPostgres) BuscarPorID(ctx context.Context, id string) (*entity.Cliente, error) {
	query := `SELECT id, nome, email, senha_hash, criado_em 
              FROM clientes 
              WHERE id = $1`

	var cliente entity.Cliente
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&cliente.ID,
		&cliente.Nome,
		&cliente.Email,
		&cliente.SenhaHash,
		&cliente.CriadoEm,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar cliente: %w", err)
	}
	return &cliente, nil
}

// AtualizarCliente atualiza os dados de um cliente
func (r *ClienteRepositoryPostgres) AtualizarCliente(ctx context.Context, cliente *entity.Cliente) error {
	query := `UPDATE clientes 
              SET nome = $1, email = $2, senha_hash = $3 
              WHERE id = $4`

	_, err := r.db.ExecContext(ctx, query,
		cliente.Nome,
		cliente.Email,
		cliente.SenhaHash,
		cliente.ID,
	)
	if err != nil {
		return fmt.Errorf("erro ao atualizar cliente: %w", err)
	}
	return nil
}
