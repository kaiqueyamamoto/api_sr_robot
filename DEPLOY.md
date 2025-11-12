# 🚀 Guia de Deploy - SR Robot API

## Deploy com Nixpacks

Este projeto está configurado para deploy com Nixpacks, compatível com:
- **Railway** ✅
- **Render** ✅
- **Fly.io** ✅
- **Heroku** ✅

## 📋 Pré-requisitos

### Variáveis de Ambiente Obrigatórias

```bash
MONGODB_URL=mongodb+srv://user:pass@cluster.mongodb.net/
MONGODB_DATABASE=sr_robot
PORT=8080
ENV=production
GIN_MODE=release
```

### Variáveis Opcionais

```bash
JWT_SECRET=your-super-secret-key-here
```

## 🚂 Deploy na Railway

### 1. Via CLI

```bash
# Instalar Railway CLI
npm i -g @railway/cli

# Login
railway login

# Iniciar projeto
railway init

# Adicionar variáveis de ambiente
railway variables set MONGODB_URL="mongodb+srv://..."
railway variables set MONGODB_DATABASE="sr_robot"
railway variables set PORT="8080"
railway variables set ENV="production"
railway variables set GIN_MODE="release"

# Deploy
railway up
```

### 2. Via Dashboard

1. Acesse https://railway.app
2. Click em "New Project"
3. Selecione "Deploy from GitHub repo"
4. Escolha seu repositório
5. Configure as variáveis de ambiente:
   - `MONGODB_URL`
   - `MONGODB_DATABASE`
   - `PORT`
   - `ENV`
   - `GIN_MODE`
6. Click em "Deploy"

### 3. Configuração Automática

O projeto já está configurado com:
- ✅ `railway.toml` - Configuração Railway
- ✅ `nixpacks.toml` - Build configuration
- ✅ `Procfile` - Start command
- ✅ Health check em `/health`

## 🎨 Deploy no Render

### 1. Via Dashboard

1. Acesse https://render.com
2. Click em "New +"
3. Selecione "Web Service"
4. Conecte seu repositório
5. Configure:
   - **Name**: sr-robot-api
   - **Environment**: Go
   - **Build Command**: `swag init -g main.go --output ./docs && go build -o chatserver main.go`
   - **Start Command**: `./chatserver`
   - **Port**: 8080

6. Adicione variáveis de ambiente:
   ```
   MONGODB_URL=mongodb+srv://...
   MONGODB_DATABASE=sr_robot
   PORT=8080
   ENV=production
   GIN_MODE=release
   ```

7. Click em "Create Web Service"

### 2. Via render.yaml

Crie `render.yaml`:

```yaml
services:
  - type: web
    name: sr-robot-api
    env: go
    buildCommand: swag init -g main.go --output ./docs && go build -o chatserver main.go
    startCommand: ./chatserver
    healthCheckPath: /health
    envVars:
      - key: PORT
        value: 8080
      - key: ENV
        value: production
      - key: GIN_MODE
        value: release
      - key: MONGODB_URL
        sync: false
      - key: MONGODB_DATABASE
        value: sr_robot
```

## 🪂 Deploy no Fly.io

### 1. Instalação

```bash
# Instalar Fly CLI
curl -L https://fly.io/install.sh | sh

# Login
fly auth login

# Iniciar projeto
fly launch
```

### 2. Configuração

O `fly.toml` será gerado automaticamente. Edite:

```toml
app = "sr-robot-api"
primary_region = "gru" # São Paulo

[build]
  builder = "paketobuildpacks/builder:base"

[env]
  PORT = "8080"
  ENV = "production"
  GIN_MODE = "release"

[http_service]
  internal_port = 8080
  force_https = true
  auto_stop_machines = true
  auto_start_machines = true
  min_machines_running = 1
  
  [[http_service.checks]]
    grace_period = "10s"
    interval = "30s"
    method = "GET"
    timeout = "5s"
    path = "/health"
```

### 3. Secrets

```bash
fly secrets set MONGODB_URL="mongodb+srv://..."
fly secrets set MONGODB_DATABASE="sr_robot"
```

### 4. Deploy

```bash
fly deploy
```

## 🐳 Deploy com Docker

### 1. Build

```bash
docker build -t sr-robot-api .
```

### 2. Run Local

```bash
docker run -p 8080:8080 \
  -e MONGODB_URL="mongodb+srv://..." \
  -e MONGODB_DATABASE="sr_robot" \
  -e PORT="8080" \
  -e ENV="production" \
  -e GIN_MODE="release" \
  sr-robot-api
```

### 3. Docker Compose

```yaml
version: '3.8'

services:
  api:
    build: .
    ports:
      - "8080:8080"
    environment:
      - MONGODB_URL=${MONGODB_URL}
      - MONGODB_DATABASE=${MONGODB_DATABASE}
      - PORT=8080
      - ENV=production
      - GIN_MODE=release
    restart: unless-stopped
```

## ☁️ Deploy no Heroku

### 1. Via CLI

```bash
# Login
heroku login

# Criar app
heroku create sr-robot-api

# Adicionar buildpack
heroku buildpacks:set heroku/go

# Configurar variáveis
heroku config:set MONGODB_URL="mongodb+srv://..."
heroku config:set MONGODB_DATABASE="sr_robot"
heroku config:set PORT="8080"
heroku config:set ENV="production"
heroku config:set GIN_MODE="release"

# Deploy
git push heroku main
```

## 🔧 Verificação de Deploy

### 1. Health Check

```bash
curl https://seu-app.railway.app/health
```

**Resposta esperada:**
```json
{
  "status": "ok",
  "service": "sr_robot_api"
}
```

### 2. Swagger

```
https://seu-app.railway.app/swagger/index.html
```

### 3. Teste de API

```bash
# Registrar
curl -X POST https://seu-app.railway.app/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"test123"}'

# Chat
curl -X POST https://seu-app.railway.app/api/v1/chat \
  -H "Content-Type: application/json" \
  -d '{"message":"Olá!"}'
```

## 📊 Monitoramento

### Railway

- Logs: `railway logs`
- Metrics: Dashboard da Railway
- Health: Configurado em `railway.toml`

### Render

- Logs: Dashboard → Logs
- Metrics: Dashboard → Metrics
- Health: `/health` endpoint

### Fly.io

```bash
# Logs
fly logs

# Status
fly status

# Metrics
fly dashboard
```

## 🐛 Troubleshooting

### Erro: "Port already in use"

Certifique-se que a variável `PORT` está configurada corretamente.

### Erro: "MongoDB connection failed"

Verifique:
1. `MONGODB_URL` está correta
2. IP do servidor está na whitelist do MongoDB Atlas
3. Credenciais estão corretas

### Build falhou

```bash
# Limpar cache (Railway)
railway run --clean

# Ver logs de build
railway logs --build
```

### Swagger não carrega

Certifique-se que o build command inclui:
```bash
swag init -g main.go --output ./docs
```

## 🔐 Segurança em Produção

### 1. Variáveis de Ambiente

✅ Use secrets para dados sensíveis  
✅ Não commite arquivos `.env`  
✅ Use JWT secret forte  

### 2. MongoDB

✅ Use MongoDB Atlas com autenticação  
✅ Configure IP whitelist  
✅ Use conexão TLS/SSL  

### 3. API

✅ Configure rate limiting  
✅ Use HTTPS (automático na maioria das plataformas)  
✅ Monitore logs de autenticação  

## 📈 Escalabilidade

### Railway

```bash
# Escalar horizontalmente
railway scale --replicas 3

# Escalar verticalmente
railway scale --size medium
```

### Render

- Auto-scaling disponível no plano Pro
- Configure em Instance Type

### Fly.io

```toml
[http_service]
  min_machines_running = 2
  max_machines_running = 10
```

## 🎯 Checklist de Deploy

- [ ] Variáveis de ambiente configuradas
- [ ] MongoDB Atlas configurado
- [ ] Health check funcionando
- [ ] Swagger acessível
- [ ] Testes de API passando
- [ ] Logs configurados
- [ ] Monitoramento ativo
- [ ] HTTPS habilitado
- [ ] Domínio customizado (opcional)

## 📚 Recursos

- [Nixpacks Docs](https://nixpacks.com/docs)
- [Railway Docs](https://docs.railway.app)
- [Render Docs](https://render.com/docs)
- [Fly.io Docs](https://fly.io/docs)

## 🆘 Suporte

**Logs em tempo real:**

```bash
# Railway
railway logs --follow

# Render
# Dashboard → Logs

# Fly.io
fly logs --follow
```

**Restart:**

```bash
# Railway
railway run restart

# Fly.io
fly apps restart
```

---

**Status**: ✅ Pronto para deploy em produção!

