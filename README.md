# 🚦 Desafio Rate Limiter

Este projeto implementa um **Rate Limiter** em Go, capaz de limitar requisições por IP ou token (API\_KEY), com uso de Redis como backend e arquitetura desacoplada.

## 📦 Estrutura do Projeto

```
desafio_rate_limiter/
├── cmd/server            # Ponto de entrada da aplicação
├── internal/limiter      # Implementação da lógica de rate limiting
│   ├── strategy          # Estratégia de backend (ex: Redis)
│   └── config            # Configuração por env
├── Dockerfile            # Build da aplicação
├── docker-compose.yml    # Ambiente local com Redis
├── .env                  # Configurações de ambiente
├── go.mod / go.sum       # Dependências Go
└── README.md             # Este arquivo
```

---

## 🚀 Como executar

### 🔧 Requisitos

* Go 1.24
* Docker e Docker Compose
* Redis

### 🛠️ Configuração do `.env`

Crie um arquivo `.env` com as variáveis:

```env
# Limites gerais por IP
RATE_LIMIT_IP=5
BLOCK_DURATION_IP=300

# Limites por Token (pode ser sobrescrito no Redis)
RATE_LIMIT_TOKEN=10
BLOCK_DURATION_TOKEN=300

# Redis
REDIS_ADDR=redis:6379
REDIS_DB=0
REDIS_PASSWORD=
```

---

### ✅ Rodando localmente

#### 1. Suba a Aplicação

```bash
go run cmd/server/main.go 
```

A aplicação estará disponível em: http://localhost:8080

#### 2. Execute as chamadas na porta 8080

```bash
curl -H "API_KEY: abc123" http://localhost:8080
```

---

### 🐳 Rodando via Docker

```bash
docker build -t rate-limiter .
docker-compose up --build
```

---

## 🧪 Rodando os Testes

### Testes unitários:

```bash
go test -v ./...
```
---

## 🧠 Como funciona

O middleware de rate limiting usa:

* IP do cliente ou `API_KEY` no header
* Redis para controle de contagem
* Estratégia desacoplada para facilitar futuras trocas de backend

---

## ✍️ Exemplos de uso

Requisição sem `API_KEY`:

```bash
curl http://localhost:8080
```

Requisição com token:

```bash
curl -H "API_KEY: abc123" http://localhost:8080
```

---

## 🔐 Licença

Este projeto está licenciado sob a [MIT License](LICENSE).
