package postgres

import (
    "context"
    "database/sql"
    "lambda_tracker/internal/domain/entity"
    "time"
    "fmt"
)

// APIKeyRepositoryPostgres implementa a interface APIKeyRepository
// é a única parte do sistema que conhece SQL e PostgreSQL
type APIKeyRepositoryPostgres struct {
    db *sql.DB // campo privado — ninguém acessa diretamente
}

// construtor — recebe o banco e retorna o repositório pronto
func NewAPIKeyRepositoryPostgres(db *sql.DB) *APIKeyRepositoryPostgres {
    return &APIKeyRepositoryPostgres{db: db}
}

// Salvar implementa APIKeyRepository.Salvar
// ExecContext para INSERT pois não retorna linhas
func (r *APIKeyRepositoryPostgres) Salvar(ctx context.Context, apiKey *entity.APIKey) error {
    query := `INSERT INTO api_keys 
              (id, cliente_id, key_hash, descricao, ativo, criado_em) 
              VALUES ($1, $2, $3, $4, $5, $6)`

    // ExecContext — executa sem retornar linhas
    // ctx — permite cancelar se demorar demais
    // $1,$2... — placeholders evitam SQL Injection
    _, err := r.db.ExecContext(ctx, query,
        apiKey.ID,
        apiKey.ClienteID,
        apiKey.KeyHash,   // ← salva o HASH, nunca a chave real
        apiKey.Descricao,
        apiKey.Ativo,
        apiKey.CriadoEm,
    )
    return err
}

// BuscarPorHash — QueryRowContext pois retorna UMA linha
func (r *APIKeyRepositoryPostgres) BuscarPorHash(ctx context.Context, keyHash string) (*entity.APIKey, error) {
    query := `SELECT id, cliente_id, key_hash, descricao, ativo, criado_em, ultimo_uso 
              FROM api_keys 
              WHERE key_hash = $1 
              AND ativo = true` // ← já filtra keys revogadas aqui no SQL

    var apiKey entity.APIKey

    // sql.NullTime — usado para colunas que podem ser NULL no banco
    // *time.Time no Go seria nil, mas o Scan precisa de NullTime
    var ultimoUso sql.NullTime

    // QueryRowContext — retorna exatamente uma linha
    // Scan — mapeia cada coluna do resultado para cada variável
    //        a ordem deve ser EXATAMENTE igual ao SELECT acima
    err := r.db.QueryRowContext(ctx, query, keyHash).Scan(
        &apiKey.ID,        // $1 → id
        &apiKey.ClienteID, // $2 → cliente_id
        &apiKey.KeyHash,   // $3 → key_hash
        &apiKey.Descricao, // $4 → descricao
        &apiKey.Ativo,     // $5 → ativo
        &apiKey.CriadoEm,  // $6 → criado_em
        &ultimoUso,        // $7 → ultimo_uso (pode ser NULL)
    )
    if err != nil {
        return nil, err // retorna nil se não encontrou
    }

    // converte NullTime para *time.Time apenas se tiver valor
    if ultimoUso.Valid {
        apiKey.UltimoUso = &ultimoUso.Time
    }

    return &apiKey, nil
}

// BuscarPorClienteID retorna todas as chaves de um cliente
func (r *APIKeyRepositoryPostgres) BuscarPorClienteID(ctx context.Context, clienteID string) ([]entity.APIKey, error) {
    query := `SELECT id, cliente_id, key_hash, descricao, ativo, criado_em, ultimo_uso
              FROM api_keys
              WHERE cliente_id = $1`

    rows, err := r.db.QueryContext(ctx, query, clienteID)
    if err != nil {
        return nil, fmt.Errorf("erro ao buscar api keys: %w", err)
    }
    defer rows.Close()

    var result []entity.APIKey
    for rows.Next() {
        var apiKey entity.APIKey
        var ultimoUso sql.NullTime
        if err := rows.Scan(&apiKey.ID, &apiKey.ClienteID, &apiKey.KeyHash, &apiKey.Descricao, &apiKey.Ativo, &apiKey.CriadoEm, &ultimoUso); err != nil {
            return nil, err
        }
        if ultimoUso.Valid {
            apiKey.UltimoUso = &ultimoUso.Time
        }
        result = append(result, apiKey)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

    return result, nil
}

// AtualizarUltimoUso atualiza o timestamp do último uso
func (r *APIKeyRepositoryPostgres) AtualizarUltimoUso(ctx context.Context, apiKeyID string) error {
    query := `UPDATE api_keys SET ultimo_uso = $1 WHERE id = $2`
    _, err := r.db.ExecContext(ctx, query, time.Now(), apiKeyID)
    return err
}

// Revogar marca a API key como inativa
func (r *APIKeyRepositoryPostgres) Revogar(ctx context.Context, apiKeyID string) error {
    query := `UPDATE api_keys SET ativo = false WHERE id = $1`
    _, err := r.db.ExecContext(ctx, query, apiKeyID)
    return err
}