# 🏗️ Arquitetura da API SR Robot

## 📊 Diagrama de Fluxo

```
┌─────────────┐
│   Cliente   │
│ (Web/App)   │
└──────┬──────┘
       │ POST /api/v1/chat
       │ {message, conversationId?}
       ▼
┌──────────────────────────────────────┐
│         API Go (Gin Framework)       │
│                                      │
│  ┌────────────────────────────────┐ │
│  │   Chat Controller              │ │
│  │                                │ │
│  │  1. Validar request            │ │
│  │  2. Criar/Recuperar conversa   │ │
│  │  3. Salvar mensagem usuário    │ │
│  │  4. Buscar histórico           │ │
│  │  5. Chamar n8n webhook         │ │
│  │  6. Salvar resposta            │ │
│  │  7. Retornar resultado         │ │
│  └────────────────────────────────┘ │
│                                      │
└───────┬──────────────────┬───────────┘
        │                  │
        │                  │ HTTP POST
        │                  ▼
        │         ┌─────────────────┐
        │         │   n8n Webhook   │
        │         │   (Produção)    │
        │         │                 │
        │         │  - Processamento│
        │         │  - RAG/IA       │
        │         │  - Resposta     │
        │         └─────────────────┘
        │
        ▼
┌──────────────────────┐
│   MongoDB Atlas      │
│                      │
│  Collections:        │
│  • conversations     │
│  • messages          │
└──────────────────────┘
```

## 🗂️ Estrutura de Dados

### Relacionamento Entre Collections

```
conversations (1)  ──────  (N) messages
      │
      │ _id
      │
      └──> messages.conversationId
```

### Exemplo de Conversa Completa

```javascript
// Collection: conversations
{
  "_id": ObjectId("674a1b2c3d4e5f6789abcdef"),
  "userId": "",
  "title": "Nova Conversa",
  "createdAt": ISODate("2025-11-12T10:00:00Z"),
  "updatedAt": ISODate("2025-11-12T10:05:00Z")
}

// Collection: messages
[
  {
    "_id": ObjectId("674a1b2c3d4e5f6789abcd00"),
    "conversationId": ObjectId("674a1b2c3d4e5f6789abcdef"),
    "role": "user",
    "content": "Olá! Quanto tempo de experiência você tem com Node.js?",
    "createdAt": ISODate("2025-11-12T10:00:00Z")
  },
  {
    "_id": ObjectId("674a1b2c3d4e5f6789abcd01"),
    "conversationId": ObjectId("674a1b2c3d4e5f6789abcdef"),
    "role": "assistant",
    "content": "Tenho 5 anos de experiência com Node.js...",
    "latencyMs": 1250,
    "metadata": {
      "sources": ["linkedin", "github"]
    },
    "createdAt": ISODate("2025-11-12T10:00:01.250Z")
  },
  {
    "_id": ObjectId("674a1b2c3d4e5f6789abcd02"),
    "conversationId": ObjectId("674a1b2c3d4e5f6789abcdef"),
    "role": "user",
    "content": "E com Golang?",
    "createdAt": ISODate("2025-11-12T10:01:00Z")
  },
  {
    "_id": ObjectId("674a1b2c3d4e5f6789abcd03"),
    "conversationId": ObjectId("674a1b2c3d4e5f6789abcdef"),
    "role": "assistant",
    "content": "Tenho 2 anos de experiência com Go...",
    "latencyMs": 980,
    "createdAt": ISODate("2025-11-12T10:01:00.980Z")
  }
]
```

## 🔄 Fluxo Detalhado de uma Requisição

### 1. Cliente envia mensagem

```http
POST /api/v1/chat
Content-Type: application/json

{
  "message": "quanto tempo de experiência tenho com nodejs?",
  "conversationId": "674a1b2c3d4e5f6789abcdef"  // opcional
}
```

### 2. API processa (chat_controller.go)

```go
// 2.1. Validar request
if err := c.ShouldBindJSON(&req); err != nil {
    return BadRequest
}

// 2.2. Obter ou criar conversa
if req.ConversationID != "" {
    conversation = FindConversation(conversationID)
} else {
    conversation = CreateNewConversation()
}

// 2.3. Salvar mensagem do usuário
userMessage = Message{
    conversationId: conversation.ID,
    role: "user",
    content: req.Message
}
SaveMessage(userMessage)

// 2.4. Buscar histórico (últimas 10 mensagens)
history = GetConversationHistory(conversation.ID, limit: 10)

// 2.5. Chamar n8n
n8nResponse = CallN8NWebhook({
    message: req.Message,
    conversationId: conversation.ID,
    history: history
})

// 2.6. Salvar resposta do assistente
assistantMessage = Message{
    conversationId: conversation.ID,
    role: "assistant",
    content: n8nResponse.Response,
    latencyMs: calculatedLatency
}
SaveMessage(assistantMessage)

// 2.7. Retornar resposta
return {
    conversationId: conversation.ID,
    message: n8nResponse.Response,
    messageId: assistantMessage.ID,
    latencyMs: latency
}
```

### 3. n8n processa

```
n8n Webhook recebe:
├── message: "quanto tempo de experiência tenho com nodejs?"
├── conversationId: "674a1b2c3d4e5f6789abcdef"
└── history: [últimas 10 mensagens]

n8n executa:
├── Análise da pergunta
├── Busca em bases de dados (RAG)
├── Processamento com IA
├── Geração de resposta
└── Retorna: {response, metadata}
```

### 4. Cliente recebe resposta

```json
{
  "conversationId": "674a1b2c3d4e5f6789abcdef",
  "message": "Tenho 5 anos de experiência com Node.js...",
  "role": "assistant",
  "messageId": "674a1b2c3d4e5f6789abcd01",
  "latencyMs": 1250
}
```

## 📦 Camadas da Aplicação

```
┌─────────────────────────────────────┐
│          main.go                    │
│  • Inicialização                    │
│  • Configuração de rotas            │
│  • Middlewares                      │
│  • Conexão DB                       │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│       Controllers Layer             │
│  • chat_controller.go               │
│    - SendMessage()                  │
│    - GetConversationHistory()       │
│    - ListConversations()            │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│         Models Layer                │
│  • conversation.go                  │
│  • message.go                       │
│    - Structs                        │
│    - Validation                     │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│       Database Layer                │
│  • mongodb.go                       │
│    - Connect()                      │
│    - GetCollection()                │
│    - Disconnect()                   │
└─────────────────────────────────────┘
```

## 🌐 Endpoints da API

| Método | Endpoint | Descrição | Auth |
|--------|----------|-----------|------|
| `GET` | `/health` | Health check | ❌ |
| `POST` | `/api/v1/chat` | Enviar mensagem | ❌ |
| `GET` | `/api/v1/conversations` | Listar conversas | ❌ |
| `GET` | `/api/v1/conversations/:id` | Buscar conversa específica | ❌ |

## 🔐 Segurança (Futuro)

### Recomendações para Produção

```
┌────────────────────────────────────┐
│  Cliente                           │
└────────────┬───────────────────────┘
             │ JWT Token
             ▼
┌────────────────────────────────────┐
│  API Gateway / Load Balancer       │
│  • Rate Limiting                   │
│  • SSL/TLS                         │
│  • DDoS Protection                 │
└────────────┬───────────────────────┘
             │
             ▼
┌────────────────────────────────────┐
│  API Go                            │
│  • JWT Validation Middleware       │
│  • CORS                            │
│  • Input Validation                │
│  • Sanitization                    │
└────────────┬───────────────────────┘
             │
             ▼
┌────────────────────────────────────┐
│  MongoDB Atlas                     │
│  • Network Access Control          │
│  • Encryption at Rest              │
│  • Audit Logs                      │
└────────────────────────────────────┘
```

## 📈 Escalabilidade

### Horizontal Scaling

```
          ┌───────────────┐
          │ Load Balancer │
          └───────┬───────┘
                  │
        ┌─────────┼─────────┐
        │         │         │
    ┌───▼───┐ ┌──▼───┐ ┌──▼───┐
    │ API 1 │ │ API 2│ │ API 3│
    └───┬───┘ └──┬───┘ └──┬───┘
        │        │        │
        └────────┼────────┘
                 │
          ┌──────▼───────┐
          │ MongoDB      │
          │ (Sharded)    │
          └──────────────┘
```

### Performance Otimizations

1. **Caching** (Redis)
   - Cache de respostas frequentes
   - Session storage
   - Rate limiting

2. **Database Indexing**
   ```javascript
   // Índices sugeridos
   db.conversations.createIndex({ "createdAt": -1 })
   db.conversations.createIndex({ "userId": 1 })
   db.messages.createIndex({ "conversationId": 1, "createdAt": -1 })
   ```

3. **Connection Pooling**
   - Reuso de conexões MongoDB
   - Keep-alive HTTP connections

## 🔍 Monitoring & Observability

### Métricas Recomendadas

```
┌─────────────────────────────────────┐
│  Application Metrics                │
│  • Request rate                     │
│  • Response time (p50, p95, p99)    │
│  • Error rate                       │
│  • Active connections               │
└─────────────────────────────────────┘

┌─────────────────────────────────────┐
│  Business Metrics                   │
│  • Conversations created per hour   │
│  • Messages per conversation        │
│  • Average latency per message      │
│  • n8n webhook success rate         │
└─────────────────────────────────────┘

┌─────────────────────────────────────┐
│  Infrastructure Metrics             │
│  • CPU usage                        │
│  • Memory usage                     │
│  • MongoDB connections              │
│  • Network I/O                      │
└─────────────────────────────────────┘
```

### Logs Estruturados

```json
{
  "timestamp": "2025-11-12T10:00:00Z",
  "level": "info",
  "service": "sr_robot_api",
  "conversationId": "674a1b2c3d4e5f6789abcdef",
  "event": "message_sent",
  "latencyMs": 1250,
  "userId": "user123"
}
```

## 🚀 Deploy

### Docker

```dockerfile
# Build
docker build -t sr-robot-api .

# Run
docker run -p 8080:8080 \
  -e MONGODB_URL="..." \
  -e MONGODB_DATABASE="sr_robot" \
  sr-robot-api
```

### Kubernetes (futuro)

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: sr-robot-api
spec:
  replicas: 3
  selector:
    matchLabels:
      app: sr-robot-api
  template:
    metadata:
      labels:
        app: sr-robot-api
    spec:
      containers:
      - name: api
        image: sr-robot-api:latest
        ports:
        - containerPort: 8080
        env:
        - name: MONGODB_URL
          valueFrom:
            secretKeyRef:
              name: mongodb-secret
              key: url
```

## 📚 Tecnologias Utilizadas

| Tecnologia | Versão | Uso |
|------------|--------|-----|
| Go | 1.23+ | Linguagem principal |
| Gin | Latest | Framework web |
| MongoDB | 7.0+ | Banco de dados |
| n8n | - | Automação/IA (externo) |
| Docker | - | Containerização |
| Air | - | Hot reload (dev) |

## 🎯 Roadmap

### Fase 1: MVP ✅
- [x] Estrutura básica da API
- [x] Integração com MongoDB
- [x] Rota de chat funcionando
- [x] Persistência de conversas

### Fase 2: Melhorias
- [ ] Autenticação JWT
- [ ] Rate limiting
- [ ] Testes unitários
- [ ] CI/CD pipeline

### Fase 3: Escalabilidade
- [ ] Redis cache
- [ ] Horizontal scaling
- [ ] Load balancing
- [ ] Monitoring

### Fase 4: Features Avançadas
- [ ] Streaming de respostas (SSE)
- [ ] Busca em conversas
- [ ] Analytics dashboard
- [ ] Export de conversas

