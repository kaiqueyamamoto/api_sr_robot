# SR Robot - API Go

API em Go para o chatbot SR Robot com integração ao MongoDB Atlas e n8n.

## 📁 Estrutura do Projeto

```
api/
├── controllers/
│   └── chat_controller.go    # Controlador de chat
├── models/
│   ├── conversation.go        # Model de conversa
│   └── message.go             # Model de mensagem
├── database/
│   └── mongodb.go             # Conexão com MongoDB
├── main.go                    # Entry point
├── go.mod                     # Dependências
├── Makefile                   # Comandos úteis
├── .air.toml                  # Configuração hot reload
└── .devcontainer/             # Dev Container config
```

## 🚀 Instalação

### 1. Instalar dependências

```bash
go mod download
go mod tidy
```

### 2. Configurar variáveis de ambiente

As variáveis já estão configuradas no Dev Container, mas você pode criar um `.env`:

```bash
MONGODB_URL=mongodb+srv://sr_robot:brBBTUbOqnxVpN0S@conecta-tech.pajxycn.mongodb.net/?appName=Conecta-Tech
MONGODB_DATABASE=sr_robot
PORT=8080
ENV=development
```

### 3. Rodar a aplicação

```bash
# Com hot reload
make dev
# ou
air

# Sem hot reload
go run main.go

# Build
make build
./bin/api
```

## 📡 API Endpoints

### Base URL
```
http://localhost:8080/api/v1
```

### 1. Enviar Mensagem (Chat)

**POST** `/api/v1/chat`

Envia uma mensagem para o chatbot. Se não fornecer `conversationId`, uma nova conversa é criada.

**Request:**
```json
{
  "message": "quanto tempo de experiência tenho com nodejs?",
  "conversationId": "optional-conversation-id"
}
```

**Response:**
```json
{
  "conversationId": "674a1b2c3d4e5f6789abcdef",
  "message": "Você tem X anos de experiência com Node.js...",
  "role": "assistant",
  "messageId": "674a1b2c3d4e5f6789abcd00",
  "latencyMs": 1250
}
```

**Status Codes:**
- `200` - Sucesso
- `400` - Requisição inválida
- `404` - Conversa não encontrada
- `500` - Erro interno

### 2. Buscar Histórico de Conversa

**GET** `/api/v1/conversations/:id`

Retorna todas as mensagens de uma conversa específica.

**Response:**
```json
{
  "conversation": {
    "id": "674a1b2c3d4e5f6789abcdef",
    "userId": "",
    "title": "Nova Conversa",
    "createdAt": "2025-11-12T10:00:00Z",
    "updatedAt": "2025-11-12T10:05:00Z"
  },
  "messages": [
    {
      "id": "674a1b2c3d4e5f6789abcd00",
      "conversationId": "674a1b2c3d4e5f6789abcdef",
      "role": "user",
      "content": "quanto tempo de experiência tenho com nodejs?",
      "createdAt": "2025-11-12T10:00:00Z"
    },
    {
      "id": "674a1b2c3d4e5f6789abcd01",
      "conversationId": "674a1b2c3d4e5f6789abcdef",
      "role": "assistant",
      "content": "Você tem X anos de experiência com Node.js...",
      "latencyMs": 1250,
      "createdAt": "2025-11-12T10:00:01Z"
    }
  ]
}
```

### 3. Listar Todas as Conversas

**GET** `/api/v1/conversations`

Lista todas as conversas armazenadas.

**Response:**
```json
{
  "conversations": [
    {
      "id": "674a1b2c3d4e5f6789abcdef",
      "userId": "",
      "title": "Nova Conversa",
      "createdAt": "2025-11-12T10:00:00Z",
      "updatedAt": "2025-11-12T10:05:00Z"
    }
  ],
  "total": 1
}
```

### 4. Health Check

**GET** `/health`

Verifica se a API está funcionando.

**Response:**
```json
{
  "status": "ok",
  "service": "sr_robot_api"
}
```

## 🗄️ Estrutura do MongoDB

### Collections

#### `conversations`
```javascript
{
  "_id": ObjectId,
  "userId": String,          // Opcional
  "title": String,
  "createdAt": Date,
  "updatedAt": Date
}
```

#### `messages`
```javascript
{
  "_id": ObjectId,
  "conversationId": ObjectId,
  "role": String,            // "user", "assistant", "system"
  "content": String,
  "tokens": Number,          // Opcional
  "latencyMs": Number,       // Opcional
  "metadata": Object,        // Opcional
  "createdAt": Date
}
```

## 🔗 Integração com n8n

A API chama o webhook do n8n em produção:

**URL:** `https://galaxy.conecta-tech.com.br/webhook/conversation`

**Payload enviado:**
```json
{
  "message": "mensagem do usuário",
  "conversationId": "id-da-conversa",
  "history": [
    // Array com últimas 10 mensagens
  ]
}
```

**Resposta esperada do n8n:**
```json
{
  "response": "resposta do chatbot",
  "metadata": {
    // metadados opcionais
  }
}
```

## 🧪 Testando a API

### Usando cURL

```bash
# Criar nova conversa
curl -X POST http://localhost:8080/api/v1/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "quanto tempo de experiência tenho com nodejs?"
  }'

# Continuar conversa existente
curl -X POST http://localhost:8080/api/v1/chat \
  -H "Content-Type: application/json" \
  -d '{
    "conversationId": "674a1b2c3d4e5f6789abcdef",
    "message": "e com golang?"
  }'

# Buscar histórico
curl http://localhost:8080/api/v1/conversations/674a1b2c3d4e5f6789abcdef

# Listar conversas
curl http://localhost:8080/api/v1/conversations
```

### Usando arquivo HTTP (VS Code REST Client)

Crie um arquivo `api.http`:

```http
### Health Check
GET http://localhost:8080/health

### Criar nova conversa
POST http://localhost:8080/api/v1/chat
Content-Type: application/json

{
  "message": "quanto tempo de experiência tenho com nodejs?"
}

### Continuar conversa
POST http://localhost:8080/api/v1/chat
Content-Type: application/json

{
  "conversationId": "{{conversationId}}",
  "message": "e com golang?"
}

### Buscar histórico
GET http://localhost:8080/api/v1/conversations/{{conversationId}}

### Listar conversas
GET http://localhost:8080/api/v1/conversations
```

## 🛠️ Comandos Make

```bash
make help          # Ver todos os comandos
make dev           # Rodar com hot reload
make build         # Build da aplicação
make run           # Rodar aplicação
make test          # Rodar testes
make lint          # Rodar linters
make format        # Formatar código
make clean         # Limpar arquivos temporários
```

## 📦 Dependências Principais

- **Gin** - Framework web
- **MongoDB Go Driver** - Cliente MongoDB
- **godotenv** - Carregar variáveis de ambiente

## 🚀 Deploy

Para produção, você pode usar Docker:

```bash
# Build
docker build -t sr-robot-api .

# Run
docker run -p 8080:8080 \
  -e MONGODB_URL="mongodb+srv://..." \
  -e MONGODB_DATABASE="sr_robot" \
  sr-robot-api
```

## 📝 Notas

- As conversas são criadas automaticamente na primeira mensagem
- O histórico das últimas 10 mensagens é enviado para o n8n como contexto
- Todas as mensagens e respostas são persistidas no MongoDB
- A latência de cada resposta é medida e armazenada

## 🔒 Segurança

Para produção, considere adicionar:
- Autenticação JWT
- Rate limiting
- Validação de input mais robusta
- HTTPS
- Logs estruturados
- Monitoramento

## 📚 Recursos

- [Gin Documentation](https://gin-gonic.com/docs/)
- [MongoDB Go Driver](https://www.mongodb.com/docs/drivers/go/current/)
- [Go by Example](https://gobyexample.com/)
