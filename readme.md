# Alert Service

O Alert Service é um microsserviço consumidor responsável por processar eventos de alteração de preço de voos e notificar os usuários finais.

## Visão Geral

⚠️ **Projeto em Desenvolvimento** ⚠️

Este serviço segue os princípios da Arquitetura Limpa (Clean Architecture) e Hexagonal (Ports & Adapters) para isolar a lógica de negócio de dependências externas como mensageria e APIs. O serviço:

-   📥 Consome eventos de alteração de preço de uma fila no RabbitMQ.
-   ⚙️ Garante o processamento único de cada evento utilizando Redis.
-   🛡️ Aplica regras de negócio, como limites de envio (Rate Limiting), para evitar spam.
-   📧 Orquestra o envio de notificações de alerta para os usuários via um serviço de e-mail externo.
---

## Fluxo planejado (pode sofrer alterações)

```graph TD
    subgraph "Sistema Externo"
        A[Search Service (Java)]
    end

    subgraph "Infraestrutura de Mensageria e Cache"
        B(RabbitMQ)
        C(Redis Cache)
    end

    subgraph "Alert Service (Este Repositório)"
        D[Worker Go]
    end

    subgraph "Serviço de Terceiros"
        E[API de Notificação (E-mail)]
    end

    %% Definição do Fluxo de Dados
    A -- "1. Publica evento 'price.updated'" --> B
    B -- "2. Consome evento da fila" --> D
    D -- "3. Verifica idempotência e rate limit" --> C
    D -- "4. Envia notificação" --> E

    %% Estilização dos Componentes
    style A fill:#f9f,stroke:#333,stroke-width:2px
    style B fill:#FF6600,stroke:#333,stroke-width:2px
    style C fill:#DC382D,stroke:#333,stroke-width:2px
    style D fill:#00ADD8,stroke:#333,stroke-width:2px
    style E fill:#4CAF50,stroke:#333,stroke-width:2px
  ```

  ## Estrutura planejada (pode sofrer alterações)

  ```
alert-service/
├── cmd/
│   └── worker/
│
├── internal/
│   ├── modules/
│   │   └── alerting/
│   │       ├── entity.go     #   - Structs centrais do domínio.
│   │       ├── port.go       #   - Interfaces (Portas).
│   │       ├── service.go    #   - Lógica de aplicação e casos de uso.
│   │       └── error.go      #   - Erros customizados do domínio.
│   │
│   └── infra/
│       ├── consumer/         # Adaptadores de entrada
│       │   └── rabbitmq.go
│       ├── storage/          # Adaptadores de persistência/cache
│       │   └── redis.go
│       └── notification/     # Adaptadores de saída
│           └── email.go
│
├── dockerfile                # Instruções para containerização.
├── go.mod                    # Gerenciamento de dependências Go.
└── README.md                 # Esta documentação.
```