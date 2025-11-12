# Swagger Documentation

## 📖 Acesso à Documentação

A documentação interativa da API está disponível através do Swagger UI:

**URL:** `http://localhost:8080/swagger/index.html`

## 🚀 Recursos do Swagger

### Interface Interativa

- **Explorar Endpoints**: Visualize todos os endpoints disponíveis
- **Testar Requisições**: Execute requisições diretamente da interface
- **Ver Modelos**: Veja a estrutura dos modelos de dados
- **Validação**: Validação automática de requisições

### Formato OpenAPI

A documentação segue o padrão OpenAPI 3.0 e está disponível em:

- **JSON**: `http://localhost:8080/swagger/doc.json`
- **YAML**: `/docs/swagger.yaml`

## 📝 Endpoints Documentados

### Health

- `GET /health` - Verificação de saúde do servidor

### Chat

- `POST /api/v1/chat` - Enviar mensagem para o chatbot
- `GET /api/v1/conversations/{id}` - Obter histórico de uma conversa
- `GET /api/v1/conversations` - Listar todas as conversas

## 🔧 Como Usar

### 1. Acessar a Interface

Abra o navegador em: `http://localhost:8080/swagger/index.html`

### 2. Testar um Endpoint

#### Exemplo: Enviar Mensagem

1. Clique em `POST /api/v1/chat`
2. Clique em "Try it out"
3. Preencha o corpo da requisição:

```json
{
  "message": "Olá, qual é o meu nome?"
}
```

4. Clique em "Execute"
5. Veja a resposta abaixo

#### Exemplo: Continuar Conversa

```json
{
  "conversationId": "507f1f77bcf86cd799439011",
  "message": "E onde eu trabalho?"
}
```

### 3. Ver Modelos de Dados

Role até a seção "Schemas" para ver todos os modelos:

- `ChatRequest` - Requisição de chat
- `ChatResponse` - Resposta do chat
- `Message` - Modelo de mensagem
- `Conversation` - Modelo de conversa

## 🛠️ Anotações Swagger no Código

### Anotação Principal (main.go)

```go
// @title           SR Robot API
// @version         1.0
// @description     API para chatbot SR Robot com JWT authentication e Prometheus metrics
// @host            localhost:8080
// @BasePath        /
```

### Anotação de Endpoint

```go
// @Summary      Enviar mensagem para o chatbot
// @Description  Envia uma mensagem e recebe a resposta do chatbot
// @Tags         chat
// @Accept       json
// @Produce      json
// @Param        request  body      ChatRequest  true  "Mensagem do usuário"
// @Success      200      {object}  ChatResponse
// @Router       /api/v1/chat [post]
```

## 🔄 Atualizar Documentação

Quando modificar os endpoints ou adicionar novos:

```bash
# Gerar documentação atualizada
/go/bin/swag init -g main.go --output ./docs

# Recompilar o projeto
go build -o chatserver main.go

# Reiniciar o servidor
./chatserver
```

## 📦 Dependências

```bash
# Instalar dependências do Swagger
go get -u github.com/swaggo/swag/cmd/swag
go get -u github.com/swaggo/gin-swagger
go get -u github.com/swaggo/files

# Instalar CLI do swag
go install github.com/swaggo/swag/cmd/swag@latest
```

## 🎯 Exemplos de Uso

### 1. Criar Nova Conversa

```bash
curl -X POST http://localhost:8080/api/v1/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "Qual é o meu nome?"
  }'
```

**Resposta:**

```json
{
  "conversationId": "6914c018a1f818f796048b16",
  "userMessage": "Qual é o meu nome?",
  "assistantMessage": "Seu nome é João Silva.",
  "createdAt": "2025-11-12T20:30:00Z"
}
```

### 2. Continuar Conversa Existente

```bash
curl -X POST http://localhost:8080/api/v1/chat \
  -H "Content-Type: application/json" \
  -d '{
    "conversationId": "6914c018a1f818f796048b16",
    "message": "E onde eu trabalho?"
  }'
```

### 3. Listar Conversas

```bash
curl http://localhost:8080/api/v1/conversations
```

### 4. Ver Histórico de Conversa

```bash
curl http://localhost:8080/api/v1/conversations/6914c018a1f818f796048b16
```

## 📋 Modelos de Dados

### ChatRequest

```go
type ChatRequest struct {
    ConversationID string `json:"conversationId,omitempty"` // Opcional
    Message        string `json:"message" binding:"required"`
}
```

### ChatResponse

```go
type ChatResponse struct {
    ConversationID     string    `json:"conversationId"`
    UserMessage        string    `json:"userMessage"`
    AssistantMessage   string    `json:"assistantMessage"`
    CreatedAt          time.Time `json:"createdAt"`
}
```

### Message

```go
type Message struct {
    ID             primitive.ObjectID `json:"id" bson:"_id,omitempty"`
    ConversationID primitive.ObjectID `json:"conversationId" bson:"conversationId"`
    Role           MessageRole        `json:"role" bson:"role"`
    Content        string             `json:"content" bson:"content"`
    CreatedAt      time.Time          `json:"createdAt" bson:"createdAt"`
}
```

### Conversation

```go
type Conversation struct {
    ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
    CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
    UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
}
```

## 🔍 Recursos Avançados

### Exportar Documentação

```bash
# Gerar JSON
curl http://localhost:8080/swagger/doc.json > api-docs.json

# Usar com outras ferramentas
# - Postman: Importar OpenAPI
# - Insomnia: Importar OpenAPI
# - API Testing Tools
```

### Integração com CI/CD

```yaml
# .github/workflows/swagger.yml
- name: Generate Swagger
  run: |
    go install github.com/swaggo/swag/cmd/swag@latest
    swag init -g main.go --output ./docs

- name: Publish Docs
  # Publicar documentação em ambiente de staging/produção
```

## 🎨 Personalização

### Temas e Estilos

A interface Swagger pode ser personalizada através de configurações:

```go
// Configurar Swagger com opções customizadas
url := ginSwagger.URL("http://localhost:8080/swagger/doc.json")
router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
```

## ⚡ Dicas

1. **Documentação Sempre Atualizada**: Rode `swag init` após modificar endpoints
2. **Testes Rápidos**: Use o Swagger UI para testar rapidamente
3. **Compartilhar API**: Compartilhe a URL do Swagger com outros desenvolvedores
4. **Exportar para Postman**: Baixe o JSON e importe no Postman
5. **Validação**: O Swagger valida automaticamente os tipos de dados

## 🐛 Troubleshooting

### Documentação Não Atualiza

```bash
# Limpar e regenerar
rm -rf docs/
/go/bin/swag init -g main.go --output ./docs
go build -o chatserver main.go
./chatserver
```

### Erro 404 no Swagger

- Verifique se o import `_ "chatserver/docs"` está presente
- Confirme que a pasta `docs/` foi gerada
- Reinicie o servidor

### Tipos Não Aparecem

- Adicione comentários nos structs exportados
- Use tags JSON nos campos
- Rode `swag init` novamente

## 📚 Referências

- [Swagger Official](https://swagger.io/)
- [gin-swagger](https://github.com/swaggo/gin-swagger)
- [swag Documentation](https://github.com/swaggo/swag)
- [OpenAPI Specification](https://swagger.io/specification/)
