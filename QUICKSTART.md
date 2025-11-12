# 🚀 Quick Start - SR Robot API

Guia rápido para começar a desenvolver a API do SR Robot.

## 📋 Pré-requisitos

- Go 1.23+ instalado
- MongoDB Atlas configurado
- VS Code com extensão "Dev Containers" (recomendado)

## 🏃 Início Rápido

### Opção 1: Dev Container (Recomendado)

1. **Abrir no Dev Container**
   ```bash
   # No VS Code:
   # F1 > "Dev Containers: Reopen in Container"
   ```

2. **As dependências serão instaladas automaticamente**
   
3. **Rodar a API**
   ```bash
   make dev
   # ou
   air
   ```

4. **Testar**
   - Abra o arquivo `api.http`
   - Clique em "Send Request" acima de cada requisição

### Opção 2: Local (Windows)

1. **Instalar dependências**
   ```bash
   cd api
   go mod download
   ```

2. **Rodar a API**
   ```bash
   # Com hot reload
   go install github.com/cosmtrek/air@latest
   air

   # Ou sem hot reload
   go run main.go
   ```

3. **Testar**
   ```bash
   curl http://localhost:8080/health
   ```

## 🧪 Testando a API

### 1. Health Check

```bash
curl http://localhost:8080/health
```

Resposta esperada:
```json
{
  "status": "ok",
  "service": "sr_robot_api"
}
```

### 2. Enviar primeira mensagem (criar conversa)

```bash
curl -X POST http://localhost:8080/api/v1/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "Olá! Quanto tempo de experiência você tem com Node.js?"
  }'
```

Resposta esperada:
```json
{
  "conversationId": "674a1b2c3d4e5f6789abcdef",
  "message": "Resposta do chatbot aqui...",
  "role": "assistant",
  "messageId": "674a1b2c3d4e5f6789abcd00",
  "latencyMs": 1250
}
```

### 3. Continuar a conversa

```bash
curl -X POST http://localhost:8080/api/v1/chat \
  -H "Content-Type: application/json" \
  -d '{
    "conversationId": "674a1b2c3d4e5f6789abcdef",
    "message": "E com Golang?"
  }'
```

### 4. Ver histórico da conversa

```bash
curl http://localhost:8080/api/v1/conversations/674a1b2c3d4e5f6789abcdef
```

### 5. Listar todas as conversas

```bash
curl http://localhost:8080/api/v1/conversations
```

## 📁 Estrutura Criada

```
api/
├── controllers/
│   └── chat_controller.go       # ✅ Lógica do chat
├── models/
│   ├── conversation.go          # ✅ Model de conversa
│   └── message.go               # ✅ Model de mensagem
├── database/
│   └── mongodb.go               # ✅ Conexão MongoDB
├── main.go                      # ✅ Entry point configurado
├── go.mod                       # ✅ Dependências
├── api.http                     # ✅ Arquivo de testes HTTP
├── README.md                    # ✅ Documentação completa
├── Dockerfile                   # ✅ Para deploy
└── Makefile                     # ✅ Comandos úteis
```

## 🗄️ Collections MongoDB

A API cria automaticamente as seguintes collections:

### `conversations`
```javascript
{
  "_id": ObjectId("..."),
  "userId": "",                    // Opcional
  "title": "Nova Conversa",
  "createdAt": ISODate("..."),
  "updatedAt": ISODate("...")
}
```

### `messages`
```javascript
{
  "_id": ObjectId("..."),
  "conversationId": ObjectId("..."),
  "role": "user",                  // ou "assistant"
  "content": "texto da mensagem",
  "latencyMs": 1250,
  "metadata": {},
  "createdAt": ISODate("...")
}
```

## 🔗 Fluxo de uma Mensagem

1. **Cliente envia POST** para `/api/v1/chat`
2. **API verifica** se `conversationId` existe:
   - Se sim → usa conversa existente
   - Se não → cria nova conversa
3. **API salva** a mensagem do usuário no MongoDB
4. **API busca** últimas 10 mensagens da conversa (contexto)
5. **API chama** o webhook do n8n:
   - URL: `https://galaxy.conecta-tech.com.br/webhook/conversation`
   - Payload: `{message, conversationId, history}`
6. **n8n processa** e retorna resposta
7. **API salva** resposta do assistente no MongoDB
8. **API retorna** resposta para o cliente

## 🛠️ Comandos Úteis

```bash
# Ver todos os comandos disponíveis
make help

# Desenvolvimento com hot reload
make dev

# Build da aplicação
make build

# Rodar testes
make test

# Formatar código
make format

# Rodar linters
make lint

# Limpar arquivos temporários
make clean
```

## 🔍 Debugging

### VS Code

1. Pressione `F5` ou vá em "Run and Debug"
2. Selecione "Launch Package"
3. Adicione breakpoints no código
4. Execute requisições

### Logs

A API imprime logs no console:
```
✅ Conectado ao MongoDB Atlas!
🚀 Servidor rodando na porta 8080
```

## 📝 Próximos Passos

1. ✅ API funcionando com MongoDB Atlas
2. ✅ Rota de chat criada
3. ✅ Conversas e mensagens sendo salvas
4. ✅ Integração com n8n configurada

### Melhorias Sugeridas:

- [ ] Adicionar autenticação (JWT)
- [ ] Implementar rate limiting
- [ ] Adicionar testes unitários
- [ ] Implementar cache com Redis
- [ ] Adicionar métricas e monitoring
- [ ] Implementar busca de conversas
- [ ] Adicionar paginação nas listagens
- [ ] Implementar soft delete
- [ ] Adicionar validação mais robusta
- [ ] Implementar streaming de respostas (SSE)

## 🐛 Troubleshooting

### Erro: "MONGODB_URL não configurado"

No Dev Container, a variável já está configurada. Se rodar localmente, defina:

```bash
export MONGODB_URL="mongodb+srv://sr_robot:brBBTUbOqnxVpN0S@conecta-tech.pajxycn.mongodb.net/?appName=Conecta-Tech"
```

### Erro: "go: command not found"

Certifique-se que o Go está instalado:
```bash
go version
```

### Erro ao conectar no MongoDB Atlas

Verifique se seu IP está na whitelist do MongoDB Atlas.

### API não responde

Verifique se está rodando:
```bash
curl http://localhost:8080/health
```

## 📚 Documentação Adicional

- [README.md](./README.md) - Documentação completa da API
- [api.http](./api.http) - Exemplos de requisições
- [.devcontainer/README.md](./.devcontainer/README.md) - Guia do Dev Container

## 🎉 Pronto!

Sua API está configurada e pronta para desenvolvimento!

Para testar rapidamente:

1. Abra `api.http` no VS Code
2. Clique em "Send Request" na primeira requisição (Health Check)
3. Se retornar OK, está funcionando!
4. Teste as outras requisições
