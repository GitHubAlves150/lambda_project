
## 💻 lambda_project

Este projeto Backend em Golang é usado no projeto lilygo-vehicle-tracking disponívle em https://github.com/GitHubAlves150/lilygo-vehicle-tracking.git. Para cada branc deste repositório é um avanço gradual da minha dedicação ao Golang. O projeto lilygo-tracking tem como alicerce para ser implementado a linguagem Go + as principais arquiteturas de software, framework e cleancode assim como S.O.L.I.D.
## Atenção! Todo este trabalho foi somente possível graças a ajuda e consulta da I.A Deepseek. Os estudos foi baseados em códigos de exemplo gerados e explicados pela I.A, Os demais detalhes como organização, comentários preparação de estrutura é dedicação do autor deste projeto.

## Toda a pesquisa pode ser conferida neste link https://chat.deepseek.com/share/gcw3kbzsvbnwdrrmfe .


## 🚀 FASE 1: APIs e Banco de Dados - Guia Prático 

Nesta fase é implementada os seguintes requisitos (muito pedido em vagas de emprego para área de Go). 

## 📋 Visão Geral da Fase 1 
```bash
Tarefa	                     	        Complexidade
1. Migrar para PostgreSQL	        	    Média
2. Adicionar autenticação (API Key)		    Baixa
3. Adicionar Swagger/OpenAPI		        Média
4. Adicionar testes unitários	        	Média
```
## 🟢 TAREFA 1: Migrar do DynamoDB para PostgreSQL  

Vamos confugrar o PostgresSQL localmente em um Docker (è preciso ter ele instalado). veja no prórprio arquivo docker-compose.yml os principais comando usado para docker-compose.yml

```bash
# docker-compose.yml (adicione este serviço)
services:
  postgres:
    image: postgres:15-alpine
    container_name: rastreador-postgres
    environment:
      POSTGRES_USER: rastreador
      POSTGRES_PASSWORD: rastreador123
      POSTGRES_DB: rastreador
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

volumes:
  postgres_data:
```

A conexão com o banco e o salvamento dos dados no banco, estão separados por questão de clean arquitetura e SOLID.
E a mina.go apenas orquesta os arquivos.

- È preciso go get github.com/lib/pq para abrir conexao com o banco



