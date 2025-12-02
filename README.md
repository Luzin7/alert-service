# Alert Service

Microsserviço consumidor responsável por processar eventos de alteração de preço de voos e notificar usuários por e-mail. Este projeto foi inspirado pelo [Search Service](https://github.com/maxsonferovante/search-service) do [@maxsonferovante](https://github.com/maxsonferovante), que realiza a busca e monitoramento de preços de passagens aéreas.

## Visão Geral

O Alert Service atua como um **worker assíncrono** que consome eventos de mudança de preço publicados no RabbitMQ. Quando detecta que o preço de um voo caiu para o valor esperado pelo usuário, o serviço:

1. **Recebe** eventos `price.updated` da fila RabbitMQ
2. **Valida** o payload e converte para entidades de domínio
3. **Gera** links para o Google Flights com os dados do voo
4. **Busca** o e-mail do usuário no banco de dados
5. **Envia** notificação por e-mail via SMTP

O projeto segue os princípios da **Clean Architecture** (Arquitetura Limpa) e **Hexagonal Architecture** (Ports & Adapters), garantindo alta testabilidade, baixo acoplamento e independência de frameworks externos.

---

## Arquitetura

### Diagrama de Fluxo

```mermaid
graph TD
    A[Search Service] -->|Publica price.updated| B[RabbitMQ Queue]
    B -->|Consume mensagem| C[Worker Consumer]
    C -->|Deserializa JSON| D[Handler]
    D -->|Converte para domínio| E[ProcessAlert UseCase]
    E -->|Busca e-mail| F[PostgreSQL]
    E -->|Gera link| G[Google Flights Provider]
    E -->|Envia e-mail| H[SMTP Server]
    
    style A fill:#4A90E2,stroke:#333,stroke-width:2px,color:#fff
    style B fill:#FF6B6B,stroke:#333,stroke-width:2px,color:#fff
    style C fill:#00D8B8,stroke:#333,stroke-width:2px,color:#fff
    style E fill:#FFD93D,stroke:#333,stroke-width:2px
    style F fill:#6C5CE7,stroke:#333,stroke-width:2px,color:#fff
    style H fill:#A29BFE,stroke:#333,stroke-width:2px,color:#fff
```

### Camadas da Arquitetura

#### 1. **Domain (Domínio)**
Camada central contendo as regras de negócio e entidades puras, sem dependências externas.

```
internal/domain/
├── alert.go          # Entidade principal do alerta de preço
├── alert_email.go    # Entidade de e-mail de notificação
└── contract.go       # Interfaces (ports) do domínio
```

**Exemplo de Entidade:**
```go
type Alert struct {
    ID           int64
    MessageID    string
    Origin       string
    Destination  string
    OutboundDate time.Time
    ReturnDate   time.Time
    NewPrice     float64
    OldPrice     float64
    TargetPrice  float64
    Currency     string
    CheckedAt    time.Time
    Link         string
}
```

**Contratos (Interfaces):**
```go
type LinkGenerator interface {
    Generate(origin, dest string, out, ret time.Time) string
}

type AlertRepository interface {
    GetUserEmail(ctx context.Context, alertID int64) (string, error)
}
```

#### 2. **Use Cases (Casos de Uso)**
Orquestração da lógica de negócio através de casos de uso específicos.

```
internal/usecases/
├── process_alert.go       # Caso de uso principal
└── process_alert_test.go  # Testes unitários
```

**Responsabilidades:**
- Receber alerta do handler
- Gerar link do Google Flights
- Buscar e-mail do usuário
- Enviar notificação por e-mail

#### 3. **Transport (Camada de Transporte)**
Responsável pela comunicação externa: consumo de mensagens e APIs HTTP (quando necessário).

```
internal/transport/
└── consumer/
    ├── worker.go         # Worker RabbitMQ
    ├── handler.go        # Handler de mensagens
    ├── payload.go        # Estrutura do payload JSON
    ├── handler_test.go
    └── payload_test.go
```

**Exemplo de Payload:**
```json
{
  "messageId": "abc-123",
  "alertId": 42,
  "origin": "GRU",
  "destination": "JFK",
  "outboundDate": "2025-12-15",
  "returnDate": "2025-12-20",
  "oldPrice": 2500.00,
  "newPrice": 1800.00,
  "currency": "BRL",
  "targetPrice": 2000.00,
  "toleranceUp": 100.00,
  "checkedAt": "2025-12-02T10:00:00Z"
}
```

#### 4. **Infrastructure (Infraestrutura)**
Implementações concretas dos adapters (database, cache, SMTP, providers).

```
internal/infra/
├── cache/
│   └── connection.go       # Conexão Redis (preparado para idempotência)
├── database/
│   ├── connection.go       # Pool de conexões PostgreSQL
│   ├── repository.go       # Implementação do AlertRepository
│   └── repository_test.go
├── messenger/
│   └── connection.go       # Conexão RabbitMQ
├── providers/
│   ├── google_flights.go   # Gerador de links do Google Flights
│   └── google_flights_test.go
└── smtp/
    ├── connection.go       # Conexão SMTP
    └── sender.go           # Implementação do envio de e-mails
```

---

## Tecnologias

| Tecnologia | Versão | Uso |
|-----------|---------|-----|
| **Go** | 1.24.3 | Linguagem principal |
| **RabbitMQ** | 3.x | Mensageria assíncrona |
| **PostgreSQL** | Alpine | Banco de dados relacional |
| **Redis** | Alpine | Cache e idempotência (preparado) |
| **SMTP** | - | Envio de e-mails |
| **Docker** | - | Containerização |
| **pgx** | v5 | Driver PostgreSQL nativo |
| **amqp091-go** | v1.10 | Cliente RabbitMQ oficial |

---

## Configuração de Ambiente

### Variáveis de Ambiente

Crie um arquivo `.env` na raiz do projeto baseado no `.env.example`:

```env
# Ambiente
ENV=development

# Database
DATABASE_URL=postgres://admin:admin@localhost:5432/search_service?sslmode=disable

# RabbitMQ
MESSENGER_USERNAME=admin
MESSENGER_PASSWORD=admin
MESSENGER_HOST=localhost
MESSENGER_PORT=5672
QUEUE_NAME=price-alerts

# SMTP (exemplo com Gmail)
SMTP_SERVER=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=seu-email@gmail.com
SMTP_PASSWORD=sua-senha-de-app

# Redis (em desenvolvimento)
CACHE_ADDR=localhost:6379
CACHE_USERNAME=
CACHE_PASSWORD=
CACHE_DB=0
```

### Docker Compose (Desenvolvimento)

O `docker-compose.dev.yml` sobe as dependências localmente:

```bash
docker-compose -f docker-compose.dev.yml up -d
```

**Serviços disponíveis:**
- PostgreSQL: `localhost:5432`
- RabbitMQ: `localhost:5672` (Management UI: `http://localhost:15672`)
- Redis: `localhost:6379`

---

## Como Executar

### Pré-requisitos

- Go 1.24+
- Docker e Docker Compose
- Conta SMTP configurada (Gmail, SendGrid, etc.)

### 1. Clonar o repositório

```bash
git clone https://github.com/Luzin7/alert-service.git
cd alert-service
```

### 2. Configurar variáveis de ambiente

```bash
cp .env.example .env
# Edite o arquivo .env com suas credenciais
```

### 3. Subir dependências

```bash
docker-compose -f docker-compose.dev.yml up -d
```

### 4. Instalar dependências Go

```bash
go mod download
```

### 5. Executar o worker

```bash
go run cmd/worker/main.go
```

Você verá logs como:
```
Starting worker on queue: price-alerts
Esperando mensagens na fila price-alerts...
```

---

## Estrutura Completa de Diretórios

```
alert-service/
├── cmd/
│   └── worker/
│       └── main.go                 # Entrypoint do worker
├── internal/
│   ├── domain/                     # Entidades e contratos
│   │   ├── alert.go
│   │   ├── alert_email.go
│   │   └── contract.go
│   ├── errors/                     # Erros customizados
│   │   └── api_error.go
│   ├── infra/                      # Implementações de infraestrutura
│   │   ├── cache/
│   │   │   └── connection.go
│   │   ├── database/
│   │   │   ├── connection.go
│   │   │   ├── repository.go
│   │   │   └── repository_test.go
│   │   ├── messenger/
│   │   │   └── connection.go
│   │   ├── providers/
│   │   │   ├── google_flights.go
│   │   │   └── google_flights_test.go
│   │   └── smtp/
│   │       ├── connection.go
│   │       └── sender.go
│   ├── transport/                  # Camada de transporte
│   │   ├── consumer/
│   │   │   ├── handler.go
│   │   │   ├── handler_test.go
│   │   │   ├── payload.go
│   │   │   ├── payload_test.go
│   │   │   └── worker.go
│   │   └── http/
│   │       └── server.go
│   └── usecases/                   # Casos de uso
│       ├── process_alert.go
│       └── process_alert_test.go
├── .env.example
├── .gitignore
├── docker-compose.dev.yml
├── go.mod
├── go.sum
└── README.md
```

---

## Testes

### Executar todos os testes

```bash
go test ./...
```

### Executar testes com cobertura

```bash
go test -cover ./...
```

### Testes unitários incluem:
- Payload parsing e conversão para domínio
- Gerador de links do Google Flights
- Repository de alertas (com pgxmock)
- Handler de mensagens

---

## Integração com Search Service

Este serviço foi desenvolvido para trabalhar em conjunto com o [Search Service](https://github.com/maxsonferovante/search-service), que realiza:

- Busca de preços de voos em APIs externas
- Monitoramento contínuo de alterações de preço
- Publicação de eventos `price.updated` no RabbitMQ

**Fluxo completo:**
1. Search Service detecta mudança de preço
2. Publica evento no RabbitMQ
3. Alert Service consome e notifica usuário

---

## Roadmap

- [x] Estrutura básica da Clean Architecture
- [x] Integração com RabbitMQ
- [x] Estrutura para envio de e-mails via SMTP
- [x] Geração de links do Google Flights
- [x] Testes unitários
- [x] Health check HTTP endpoint
- [ ] Implementar idempotência com Redis
- [ ] Métricas e observabilidade
- [ ] CI/CD pipeline