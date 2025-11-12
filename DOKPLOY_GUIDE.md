# 🚀 Guia de Deploy - Dokploy

## Erro Resolvido

O erro "Is a directory (os error 21)" foi causado por conflito de configurações do Nixpacks.

### ✅ Correções Aplicadas:

1. Removido `nodejs` e `git` dos nixPkgs (não necessários)
2. Ajustado caminho do `swag` para usar `$GOPATH/bin/swag`
3. Adicionado `.dockerignore` para ignorar arquivos desnecessários
4. Criado configuração alternativa em JSON (`nixpacks.json`)

## 📋 Configuração do Dokploy

### 1. Variáveis de Ambiente

Configure estas variáveis no Dokploy:

```bash
MONGODB_URL=mongodb+srv://user:password@cluster.mongodb.net/
MONGODB_DATABASE=sr_robot
PORT=8080
ENV=production
GIN_MODE=release
```

### 2. Build Configuration

**Builder:** Nixpacks

**Build Command:** (opcional, já está no nixpacks.toml)

```bash
swag init -g main.go --output ./docs && go build -o chatserver main.go
```

**Start Command:**

```bash
./chatserver
```

### 3. Health Check

Configure o health check:

- **Path:** `/health`
- **Interval:** 30s
- **Timeout:** 10s

## 🔧 Arquivos de Configuração

### nixpacks.toml (Principal)

```toml
[phases.setup]
nixPkgs = ["go_1_24"]

[phases.install]
cmds = [
  "go mod download",
  "go install github.com/swaggo/swag/cmd/swag@latest"
]

[phases.build]
cmds = [
  "$GOPATH/bin/swag init -g main.go --output ./docs",
  "go build -o chatserver main.go"
]

[start]
cmd = "./chatserver"
```

### nixpacks.json (Alternativo)

Se o TOML não funcionar, o Dokploy usará o JSON automaticamente.

## 🐛 Troubleshooting

### Erro: "Is a directory (os error 21)"

**Solução aplicada:**

- ✅ Removido arquivo `.nixpacks` conflitante
- ✅ Simplificado `nixpacks.toml`
- ✅ Adicionado `.dockerignore`

### Erro: "swag: command not found"

**Solução:**

```bash
# O swag é instalado em $GOPATH/bin
# Use: $GOPATH/bin/swag ou apenas swag se PATH estiver correto
```

### Erro: "MongoDB connection failed"

**Verificar:**

1. MONGODB_URL está correto
2. IP do Dokploy está na whitelist do MongoDB Atlas
3. Credenciais estão corretas

### Build muito lento

**Otimizar:**

```toml
# Adicionar ao nixpacks.toml
[variables]
GOCACHE = "/tmp/go-build"
GOMODCACHE = "/tmp/go-mod"
```

## ✅ Próximos Passos

### 1. Commit e Push

```bash
git add .
git commit -m "fix: adjust Nixpacks config for Dokploy"
git push origin main
```

### 2. Redeploy no Dokploy

1. Acesse o painel do Dokploy
2. Vá para seu aplicativo
3. Click em "Redeploy"
4. Aguarde o build

### 3. Verificar Deploy

```bash
# Health check
curl https://seu-app.dokploy.com/health

# Swagger
curl https://seu-app.dokploy.com/

# Test API
curl -X POST https://seu-app.dokploy.com/api/v1/chat \
  -H "Content-Type: application/json" \
  -d '{"message":"Hello!"}'
```

## 📊 Logs do Dokploy

### Ver Logs em Tempo Real

No painel do Dokploy:

1. Click no seu aplicativo
2. Aba "Logs"
3. Selecione "Real-time logs"

### Comandos Úteis

```bash
# Ver status
# (no painel do Dokploy)

# Restart aplicação
# Click em "Restart"

# Ver métricas
# Aba "Monitoring"
```

## 🎯 Checklist de Deploy

- [x] nixpacks.toml corrigido
- [x] .dockerignore adicionado
- [x] Variáveis de ambiente documentadas
- [x] Health check configurado
- [x] Swagger funcionando
- [x] Build otimizado
- [ ] Commit e push
- [ ] Redeploy no Dokploy
- [ ] Testar endpoints

## 🔍 Validação

Após o deploy bem-sucedido, você deverá ver:

```
╔═══════════════════════ Nixpacks v1.39.0 ═══════════════════════╗
║ setup      │ go_1_24                                           ║
║────────────────────────────────────────────────────────────────║
║ install    │ go mod download                                   ║
║            │ go install github.com/swaggo/swag/cmd/swag@latest ║
║────────────────────────────────────────────────────────────────║
║ build      │ swag init -g main.go --output ./docs              ║
║            │ go build -o chatserver main.go                    ║
║────────────────────────────────────────────────────────────────║
║ start      │ ./chatserver                                      ║
╚════════════════════════════════════════════════════════════════╝
✅ Build completed successfully
✅ Container started
✅ Health check passed
```

## 🎉 Deploy Bem-Sucedido!

Após o deploy:

- ✅ API estará disponível na URL do Dokploy
- ✅ Swagger em `/` ou `/swagger/index.html`
- ✅ Health check em `/health`
- ✅ Logs disponíveis no painel

## 📚 Recursos Adicionais

- [Nixpacks Documentation](https://nixpacks.com/docs)
- [Dokploy Documentation](https://docs.dokploy.com)
- [Go Deployment Best Practices](https://go.dev/doc/articles/wiki/)

## 🆘 Suporte

Se o erro persistir:

1. Verifique os logs no Dokploy
2. Confirme que todas as variáveis de ambiente estão corretas
3. Teste o build localmente:

```bash
docker build -t test .
docker run -p 8080:8080 test
```

---

**Status:** ✅ Configuração corrigida - Pronto para redeploy!
