# Serviço de Notificações — Prefeitura do Rio

API backend desenvolvida em Go para processar eventos de chamados urbanos e disponibilizar notificações aos cidadãos de forma segura, consistente e em tempo real.

---

## Objetivo

Este serviço é responsável por:

* Receber eventos de mudança de status via webhook
* Garantir integridade e idempotência dos dados
* Persistir notificações com segurança
* Disponibilizar consulta autenticada via JWT
* Entregar notificações em tempo real via WebSocket

---

## Arquitetura

O projeto segue uma arquitetura em camadas, com separação clara de responsabilidades:

```
cmd/
  main.go

internal/
  handler/        → camada HTTP (Gin) e WebSocket
  service/        → regras de negócio
  repository/     → acesso ao banco (SQL direto)
  model/          → entidades e DTOs
  database/       → conexão com PostgreSQL
```

### Princípios adotados

* Separação de responsabilidades
* Baixo acoplamento
* Código orientado a fluxo de dados
* Facilidade de manutenção e testes

---

## Tecnologias

* Go 1.24+
* Gin
* PostgreSQL
* Docker / Docker Compose
* HMAC SHA256
* JWT (autenticação)
* WebSocket (tempo real)

---

## Segurança

### Webhook (HMAC)

Cada requisição recebida deve conter o header:

```
X-Signature-256: sha256=<hash>
```

* A assinatura é validada com HMAC SHA256
* O segredo é definido via variável de ambiente (`WEBHOOK_SECRET`)
* Requisições inválidas são rejeitadas com `401`

---

### Autenticação (JWT)

Os endpoints REST exigem autenticação:

```
Authorization: Bearer <token>
```

* O CPF do usuário é extraído do claim `preferred_username`
* O acesso é restrito às notificações do próprio usuário

---

### Privacidade

* O CPF **não é armazenado em texto**
* Apenas o hash SHA256 do CPF é persistido no banco
* Garante proteção contra exposição de dados sensíveis

---

## Banco de Dados

Tabela principal:

```sql
notifications (
  id SERIAL PRIMARY KEY,
  chamado_id TEXT NOT NULL,
  cpf_hash TEXT NOT NULL,
  tipo TEXT,
  status_anterior TEXT,
  status_novo TEXT,
  titulo TEXT,
  descricao TEXT,
  timestamp TEXT,
  lida BOOLEAN DEFAULT FALSE,
  criada_em TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)
```

### Idempotência

```sql
UNIQUE (chamado_id, status_novo)
```

Evita duplicação de eventos em caso de reenvio do webhook.

---

## API REST

Todos os endpoints abaixo exigem JWT válido.

---

### Listar notificações

```
GET /notifications?page=1&limit=10
Authorization: Bearer <token>
```

Retorna as notificações do cidadão autenticado.

#### Exemplo de resposta:

```json
[
  {
    "chamado_id": "CH-2024-001234",
    "tipo": "status_change",
    "status_anterior": "em_analise",
    "status_novo": "em_execucao",
    "titulo": "Buraco na Rua — Atualização",
    "descricao": "Equipe designada",
    "timestamp": "2024-11-15T14:30:00Z",
    "lida": false
  }
]
```

---

###  Marcar como lida

PATCH /notifications/:id/read
Authorization: Bearer <token>

Marca uma notificação como lida.

* O usuário só pode alterar notificações que pertencem a ele

#### Resposta:

json
{
  "mensagem": "Notificação marcada como lida"
}

---

### Contador de não lidas

GET /notifications/unread-count
Authorization: Bearer <token>

Retorna o total de notificações não lidas.

#### Resposta:

json
{
  "unread": 3
}
---

## WebSocket

GET /ws

* Conexão autenticada via JWT
* Notificações são enviadas em tempo real
* Cada usuário recebe apenas seus próprios eventos

---

## Fluxo do sistema

1. Evento recebido via webhook
2. Validação da assinatura (HMAC)
3. Verificação de duplicidade (idempotência)
4. Persistência no banco
5. Envio em tempo real via WebSocket

---

## Como rodar o projeto

### 1. Subir ambiente completo

docker compose up --build
---

### 2. Variáveis de ambiente

Crie um arquivo `.env` baseado no `.env.example`:

DB_HOST=localhost
DB_PORT=5432
DB_USER=admin
DB_PASSWORD=admin
DB_NAME=notificacoes

WEBHOOK_SECRET=meu-secret
JWT_SECRET=meu-jwt-secret

---

### 3. Rodar aplicação

go run cmd/main.go

---

## Testes

just test

(Testes utilizam banco real para validar comportamento completo)

---

## Decisões técnicas

* **Arquitetura em camadas** → organização e manutenção
* **SQL direto** → maior controle e aderência ao desafio
* **HMAC no webhook** → segurança e integridade dos dados
* **JWT** → controle de acesso por usuário
* **Hash de CPF** → proteção de dados sensíveis
* **Idempotência no banco** → evita duplicidade
* **WebSocket** → elimina polling e melhora experiência

---

## Possíveis melhorias

* Redis para gerenciamento de conexões WebSocket
* Fila assíncrona (RabbitMQ / Kafka)
* Dead letter queue para falhas de webhook
* Observabilidade (OpenTelemetry)
* Rate limiting
* Deploy em Kubernetes

---

## Conclusão

O sistema foi projetado para atender aos requisitos do desafio com foco em:

* Confiabilidade no processamento de eventos
* Segurança e privacidade dos dados
* Clareza arquitetural
* Facilidade de evolução

---

## Autor

Daniel Pacheco
