## Executando o projeto:

**Docker:**

1. renomeie seu arquivo `.env.docker` para `.env`
2. `docker compose up -d --build`

**Terminal:**

1.  renomeie seu arquivo `.env.example` para `.env`
2.  `docker compose up redis -d`
3.  `go mod download`
4.  `go run main.go`

## Testes:

**Docker:**

1.   renomeie seu arquivo `.env.docker` para `.env`
2.  `docker compose up -d --build`
3.  `docker exec -it api go test ./... -v`

**Terminal:**

1.  `docker compose up redis -d`
2.  `go test ./... -v`

## Configurações

**Por IP:**

1.  Executar uma requisição tipo  @GET para `localhost:8080`
2.  Valor padrão: 5 requisições por segundo, com bloqueio de 1 minuto caso ultrapassado
4.  Os valores padrões podem ser alterados através das variáveis de ambiente LIMIT_REQUEST_PER_SECOND_DEFAULT e CACHE_EXPIRATION

**Por API Key:**

1. Executar uma requisição tipo @GET para `localhost:8080`, utilizando o header `API_KEY`, as quais podem ser encontradas em `api-key.json`
4. O valor padrão de limite de requisições das API keys podem ser alterados através sa da variável  de ambiente LIMITER_REQUEST_PER_SECOND_API_KEY
