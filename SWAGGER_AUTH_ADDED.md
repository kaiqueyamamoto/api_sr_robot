# ✅ Rotas de Autenticação Adicionadas ao Swagger!

## 🎉 O que Foi Feito

As rotas de **Login** e **Registro** agora estão documentadas no Swagger!

### Alterações Realizadas:

1. ✅ **Anotações Swagger** adicionadas em `controllers/auth.go`
2. ✅ **Rotas de Auth** registradas no `main.go`
3. ✅ **Documentação Swagger** regenerada
4. ✅ **request.http** atualizado e reorganizado

## 📖 Rotas de Auth no Swagger

Acesse: **http://localhost:8080/swagger/index.html**

### Novos Endpoints Documentados:

#### 1. POST /auth/register

**Registrar novo usuário**

```json
{
  "email": "user@example.com",
  "password": "senha123"
}
```

**Resposta:**

```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "email": "user@example.com",
  "user_id": "507f1f77bcf86cd799439011",
  "created_at": "2025-11-12T20:30:00Z"
}
```

#### 2. POST /auth/login

**Login de usuário**

```json
{
  "email": "user@example.com",
  "password": "senha123"
}
```

**Resposta:**

```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "email": "user@example.com",
  "user_id": "507f1f77bcf86cd799439011",
  "created_at": "2025-11-12T20:30:00Z"
}
```

**Nota:** Token JWT válido por **24 horas**

## 📋 Todos os Endpoints no Swagger

### Health

- `GET /health` - Verificação de saúde

### Autenticação

- `POST /auth/register` - Registrar usuário
- `POST /auth/login` - Login

### Chat

- `POST /api/v1/chat` - Enviar mensagem
- `GET /api/v1/conversations/{id}` - Ver histórico
- `GET /api/v1/conversations` - Listar conversas

## 🎯 Como Testar no Swagger

### 1. Acessar Swagger UI

```
http://localhost:8080/swagger/index.html
```

### 2. Registrar Usuário

1. Expanda `POST /auth/register`
2. Click "Try it out"
3. Modifique o JSON:

```json
{
  "email": "meuusuario@example.com",
  "password": "minhasenha123"
}
```

4. Click "Execute"
5. **Copie o token** da resposta!

### 3. Usar o Token (se necessário no futuro)

- Para endpoints protegidos
- Click no botão "Authorize" (cadeado no topo)
- Cole: `Bearer SEU_TOKEN_AQUI`
- Click "Authorize"

### 4. Testar Chat

1. Expanda `POST /api/v1/chat`
2. Click "Try it out"
3. Cole:

```json
{
  "message": "Olá! Qual é o meu nome?"
}
```

4. Click "Execute"
5. Veja a resposta do chatbot!

## 📝 Estrutura Atualizada

```
controllers/
├── auth.go          ✅ Com anotações Swagger
└── chat_controller.go  ✅ Com anotações Swagger

main.go              ✅ Rotas de auth registradas

docs/                ✅ Documentação regenerada
├── docs.go
├── swagger.json
└── swagger.yaml

request.http         ✅ Reorganizado e atualizado
```

## 🔧 Modelos Documentados

O Swagger agora documenta:

### Auth Models

- `RegisterRequest` - Dados de registro
- `LoginRequest` - Dados de login
- `AuthResponse` - Resposta com token

### Chat Models

- `ChatRequest` - Requisição de chat
- `ChatResponse` - Resposta do chat
- `Message` - Modelo de mensagem
- `Conversation` - Modelo de conversa

## 💡 Benefícios

✅ **Documentação Completa** - Todas as rotas documentadas  
✅ **Testes Rápidos** - Teste direto no navegador  
✅ **Modelos Visíveis** - Veja estrutura dos dados  
✅ **Validação Automática** - Swagger valida requisições  
✅ **Exportável** - Baixe para Postman/Insomnia

## 🚀 Comandos Úteis

### Iniciar Servidor

```bash
./start.sh
```

### Regenerar Swagger (após mudanças)

```bash
/go/bin/swag init -g main.go --output ./docs
go build -o chatserver main.go
./chatserver
```

### Testar Endpoints

```bash
# Registrar
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"test123"}'

# Login
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"test123"}'

# Chat
curl -X POST http://localhost:8080/api/v1/chat \
  -H "Content-Type: application/json" \
  -d '{"message":"Olá!"}'
```

## 📊 Tags no Swagger

As rotas estão organizadas por tags:

- **health** - Health check
- **auth** - Autenticação (Register, Login)
- **chat** - Chat e conversas

## 🎨 Visualização no Swagger

No Swagger UI você verá:

```
health
  ↳ GET /health - Health Check

auth
  ↳ POST /auth/register - Registrar novo usuário
  ↳ POST /auth/login - Login de usuário

chat
  ↳ POST /api/v1/chat - Enviar mensagem para o chatbot
  ↳ GET /api/v1/conversations/{id} - Obter histórico de conversa
  ↳ GET /api/v1/conversations - Listar conversas

Schemas
  ↳ RegisterRequest
  ↳ LoginRequest
  ↳ AuthResponse
  ↳ ChatRequest
  ↳ ChatResponse
  ↳ ...
```

## ✨ Pronto!

Agora seu Swagger está completo com:

- ✅ Rotas de Autenticação
- ✅ Rotas de Chat
- ✅ Modelos documentados
- ✅ Exemplos funcionais

**Acesse agora**: http://localhost:8080/swagger/index.html 🚀

---

**Dúvidas?** Consulte [SWAGGER.md](SWAGGER.md) para documentação completa.
