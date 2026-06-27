package service

import (
    "context"
    "errors"
    "lambda_tracker/internal/domain/entity"
    "lambda_tracker/internal/repository"
)

// erros de negócio tipados — permite tratar cada caso no handler
var (
    ErrInvalidCredentials = errors.New("credenciais inválidas")
    ErrAPIKeyNotFound     = errors.New("API Key não encontrada")
    ErrAPIKeyInactive     = errors.New("API Key inativa")
)

// AuthService define o CONTRATO do serviço
// o handler depende desta interface — não da implementação
type AuthService interface {
    RegistrarCliente(ctx context.Context, nome, email, senha string) (*entity.Cliente, error)
    GerarAPIKey(ctx context.Context, clienteID, descricao string) (*entity.APIKey, error)
    ValidarAPIKey(ctx context.Context, key string) (*entity.Cliente, error)
    RevogarAPIKey(ctx context.Context, apiKeyID string) error
    ListarAPIKeys(ctx context.Context, clienteID string) ([]entity.APIKey, error)
}

// AuthServiceImpl implementa AuthService
// depende das INTERFACES dos repositórios — não das implementações PostgreSQL
// isso permite testar o service com mocks
type AuthServiceImpl struct {
    clienteRepo repository.ClienteRepository // interface — não PostgresClienteRepo
    apiKeyRepo  repository.APIKeyRepository  // interface — não PostgresAPIKeyRepo
}

// construtor com injeção de dependência
// recebe as interfaces — não sabe se é postgres, mysql ou mock
func NewAuthService(
    clienteRepo repository.ClienteRepository,
    apiKeyRepo  repository.APIKeyRepository,
) *AuthServiceImpl {
    return &AuthServiceImpl{
        clienteRepo: clienteRepo,
        apiKeyRepo:  apiKeyRepo,
    }
}

func (s *AuthServiceImpl) RegistrarCliente(ctx context.Context, nome, email, senha string) (*entity.Cliente, error) {
    // regra de negócio: email único
    existing, _ := s.clienteRepo.BuscarPorEmail(ctx, email)
    if existing != nil {
        return nil, errors.New("email já registrado")
    }

    // ⚠️ simplificado — em produção usa bcrypt
    senhaHash := "hash_" + senha

    // usa o construtor da entidade — garante ID e timestamp
    cliente := entity.NovoCliente(nome, email, senhaHash)

    // delega o salvamento ao repositório
    if err := s.clienteRepo.Salvar(ctx, cliente); err != nil {
        return nil, err
    }

    return cliente, nil
}

func (s *AuthServiceImpl) ValidarAPIKey(ctx context.Context, key string) (*entity.Cliente, error) {
    // busca pelo hash da key recebida
    apiKey, err := s.apiKeyRepo.BuscarPorHash(ctx, key)
    if err != nil {
        return nil, ErrAPIKeyNotFound
    }

    // regra de negócio: key deve estar ativa
    if !apiKey.Ativo {
        return nil, ErrAPIKeyInactive
    }

    // go → atualiza o último uso em segundo plano
    // não bloqueia a resposta para o cliente
    // ← goroutine: fire and forget
    go s.apiKeyRepo.AtualizarUltimoUso(ctx, apiKey.ID)

    // busca o cliente dono da key
    return s.clienteRepo.BuscarPorID(ctx, apiKey.ClienteID)
}

func (s *AuthServiceImpl) GerarAPIKey(ctx context.Context, clienteID, descricao string) (*entity.APIKey, error) {
    // verifica se cliente existe
    cliente, err := s.clienteRepo.BuscarPorID(ctx, clienteID)
    if err != nil {
        return nil, err
    }
    if cliente == nil {
        return nil, errors.New("cliente não encontrado")
    }

    // cria nova API Key
    apiKey, err := entity.NovaAPIKey(clienteID, descricao)
    if err != nil {
        return nil, err
    }

    // salva no repositório (salva somente o hash, o campo Key é retornado ao cliente)
    if err := s.apiKeyRepo.Salvar(ctx, apiKey); err != nil {
        return nil, err
    }

    return apiKey, nil
}

func (s *AuthServiceImpl) RevogarAPIKey(ctx context.Context, apiKeyID string) error {
    return s.apiKeyRepo.Revogar(ctx, apiKeyID)
}

func (s *AuthServiceImpl) ListarAPIKeys(ctx context.Context, clienteID string) ([]entity.APIKey, error) {
    return s.apiKeyRepo.BuscarPorClienteID(ctx, clienteID)
}