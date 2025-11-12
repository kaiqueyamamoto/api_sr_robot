# Mensagem de Commit Sugerida

## Para o Próximo Commit

```
feat: add Swagger documentation and complete chat API integration

### Features
- ✨ Add Swagger/OpenAPI documentation with interactive UI
- 🔐 Add JWT authentication endpoints (register/login) to Swagger
- 💬 Add complete chat API with N8N webhook integration
- 📝 Add conversation management (list, history, update, delete)
- 📊 Add Prometheus metrics for monitoring
- 🔄 Add user info endpoint with timestamps

### API Endpoints
- POST /auth/register - User registration with JWT
- POST /auth/login - User authentication (24h token expiration)
- GET /health - Health check endpoint
- POST /api/v1/chat - Send message (create/continue conversation)
- GET /api/v1/conversations - List all conversations
- GET /api/v1/conversations/:id - Get conversation history
- PUT /api/v1/conversations/:id - Update conversation title
- DELETE /api/v1/conversations/:id - Delete conversation
- GET /swagger/index.html - Interactive API documentation

### Technical Details
- Fix N8N webhook response parsing (array format)
- Add support for both 'output' and 'response' fields
- Add conversation ID generation and tracking
- Add user timestamps (created_at, updated_at)
- Add CORS middleware
- Add environment variable configuration (PORT, MONGODB_URL)

### Documentation
- Add comprehensive Swagger annotations
- Add API_EXAMPLES.md with cURL examples
- Add SWAGGER.md with complete guide
- Add PROMETHEUS_METRICS.md for monitoring
- Update request.http with all endpoints
- Add start.sh script for easy server startup

### Dependencies
- github.com/swaggo/swag - Swagger documentation
- github.com/swaggo/gin-swagger - Gin Swagger middleware
- github.com/golang-jwt/jwt/v5 - JWT authentication
- github.com/prometheus/client_golang - Metrics
- go.mongodb.org/mongo-driver - MongoDB driver

Breaking Changes: None
Migration Required: None
```

## Commit Convencional (Conventional Commits)

### Formato Resumido

```bash
git commit -m "feat: add Swagger docs and complete chat API" \
           -m "- Add interactive Swagger UI at /swagger/index.html" \
           -m "- Add JWT auth endpoints (register/login)" \
           -m "- Add chat API with N8N integration" \
           -m "- Add conversation management endpoints" \
           -m "- Add Prometheus metrics" \
           -m "- Fix N8N response parsing" \
           -m "- Add comprehensive documentation"
```

### Formato Detalhado (Recomendado)

```bash
git commit -m "feat: add Swagger documentation and complete chat API integration

Features:
- Add Swagger/OpenAPI 3.0 documentation
- Add JWT authentication (register/login)
- Add chat API with N8N webhook integration
- Add conversation CRUD operations
- Add Prometheus metrics integration
- Add user timestamps (created_at, updated_at)

Endpoints:
- POST /auth/register, /auth/login
- POST /api/v1/chat
- GET /api/v1/conversations, /api/v1/conversations/:id
- PUT /api/v1/conversations/:id
- DELETE /api/v1/conversations/:id
- GET /swagger/index.html

Fixes:
- N8N webhook response parsing (array format)
- Conversation ID generation and tracking

Documentation:
- Swagger annotations on all endpoints
- API examples and guides
- Monitoring metrics documentation"
```

## Comandos Git

### Verificar Mudanças

```bash
git status
git diff
git log --oneline -5
```

### Adicionar e Commitar

```bash
# Adicionar todos os arquivos
git add .

# Ou adicionar seletivamente
git add controllers/ models/ main.go docs/

# Commitar com mensagem
git commit -m "feat: add Swagger and chat API"

# Commitar com descrição detalhada
git commit
# (Abrirá editor para mensagem completa)
```

### Push

```bash
# Push para branch main
git push origin main

# Push para nova branch
git checkout -b feature/swagger-chat-api
git push origin feature/swagger-chat-api
```

## Tipos de Commit (Conventional Commits)

| Tipo       | Descrição           | Exemplo                           |
| ---------- | ------------------- | --------------------------------- |
| `feat`     | Nova funcionalidade | feat: add user authentication     |
| `fix`      | Correção de bug     | fix: correct N8N response parsing |
| `docs`     | Documentação        | docs: add API examples            |
| `style`    | Formatação          | style: format code with gofmt     |
| `refactor` | Refatoração         | refactor: improve error handling  |
| `perf`     | Performance         | perf: optimize database queries   |
| `test`     | Testes              | test: add chat controller tests   |
| `chore`    | Manutenção          | chore: update dependencies        |
| `ci`       | CI/CD               | ci: add GitHub Actions            |

## Exemplo Completo

```bash
# 1. Ver o que mudou
git status

# 2. Adicionar arquivos
git add .

# 3. Commitar
git commit -m "feat: add Swagger documentation and complete chat API

Added:
- Swagger UI at /swagger/index.html
- JWT authentication endpoints
- Chat API with N8N integration
- Conversation management
- Prometheus metrics
- Comprehensive documentation

Fixed:
- N8N response parsing
- Conversation ID tracking"

# 4. Push
git push origin main
```

## Notas

- ✅ Use mensagens descritivas e claras
- ✅ Siga o padrão Conventional Commits
- ✅ Inclua breaking changes se houver
- ✅ Referencie issues quando aplicável (#123)
- ✅ Mantenha commits atômicos (uma funcionalidade por commit)

---

**Arquivo gerado em:** 2025-11-12
**Status:** Pronto para commit
