## Estrutura de pastas em um projeto Go 

```bash
meu-projeto/
│
├── cmd/                        ← Ponto de entrada da aplicação
│   └── api/
│       └── main.go             ← Só inicializa e liga tudo
│
├── internal/                   ← Código privado do projeto
│   ├── domain/                 ← Regras de negócio puras (Clean Architecture)
│   │   ├── usuario.go          ← Entidade
│   │   ├── repositorio.go      ← Interface (contrato)
│   │   └── servico.go          ← Casos de uso / regras
│   │
│   ├── handler/                ← Camada HTTP (Gin fica aqui)
│   │   └── usuario_handler.go  ← Recebe requisição, chama serviço
│   │
│   ├── repository/             ← Implementações concretas do banco
│   │   └── postgres_usuario.go ← Implementa a interface do domain
│   │
│   └── middleware/             ← Autenticação, logs, CORS
│       └── auth.go
│
├── pkg/                        ← Código reutilizável por outros projetos
│   ├── logger/                 ← Seu wrapper de logs (zap, zerolog)
│   └── validator/              ← Validações genéricas
│
├── config/                     ← Configurações da aplicação
│   └── config.go               ← Lê variáveis de ambiente (.env)
│
├── migrations/                 ← Scripts SQL de banco de dados
│   ├── 001_create_usuarios.sql
│   └── 002_create_pedidos.sql
│
├── docker/                     ← Arquivos Docker
│   └── Dockerfile
│
├── docs/                       ← Documentação gerada (Swagger)
│   └── swagger.json
│
├── .env                        ← Variáveis de ambiente (nunca sobe pro git)
├── .env.example                ← Modelo do .env (sobe pro git)
├── docker-compose.yml          ← Sobe banco, redis, rabbit localmente
├── Makefile                    ← Atalhos de comandos
└── go.mod                      ← Dependências do projeto
```