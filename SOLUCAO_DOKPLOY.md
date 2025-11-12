# ✅ SOLUÇÃO RÁPIDA - Erro Dokploy

## 🎯 Commit Gerado! Agora Siga Este Passo a Passo:

### 📝 Commit Criado:

```
704e63f - fix: add multiple solutions for Dokploy deployment error
```

---

## 🚀 SOLUÇÃO RECOMENDADA (Use Esta!)

### **Usar Dockerfile em vez de Nixpacks**

#### Passo 1: Push do Código

```bash
git push origin main
```

#### Passo 2: No Painel do Dokploy

1. **Abra** seu aplicativo no Dokploy
2. **Vá em:** Settings → General
3. **Builder:** Mude de "Nixpacks" para **"Dockerfile"**
4. **Dockerfile Path:** Digite `Dockerfile.dokploy`
5. **Click em "Save"**

#### Passo 3: Configurar Variáveis

No Dokploy, adicione estas variáveis:

```
MONGODB_URL=mongodb+srv://sr_robot:brBBTUbOqnxVpN0S@conecta-tech.pajxycn.mongodb.net/
MONGODB_DATABASE=sr_robot
PORT=8080
ENV=production
GIN_MODE=release
```

#### Passo 4: Deploy

1. Click em **"Redeploy"**
2. Aguarde 3-5 minutos
3. ✅ **Sucesso!**

---

## 📊 O que vai acontecer:

```
🔧 Building with Dockerfile...
✅ Stage 1: Installing Go and dependencies
✅ Stage 2: Generating Swagger docs
✅ Stage 3: Building Go binary
✅ Stage 4: Creating final image
✅ Starting application...
✅ Connected to MongoDB!
🚀 Server running on port 8080
```

---

## 🔍 Verificar Deploy

Após deploy bem-sucedido:

```bash
# 1. Health Check
curl https://seu-app.dokploy.com/health

# Resultado esperado:
{"status":"ok","service":"sr_robot_api"}

# 2. Swagger
curl https://seu-app.dokploy.com/

# Redireciona para /swagger/index.html

# 3. Test Chat
curl -X POST https://seu-app.dokploy.com/api/v1/chat \
  -H "Content-Type: application/json" \
  -d '{"message":"Olá!"}'
```

---

## 🆘 Se o Dockerfile Não Funcionar

### Alternativa: Configuração Manual

1. **No Dokploy:**

   - Builder: **Nixpacks**
   - **Delete** o arquivo `nixpacks.toml` (ou renomeie)

2. **Build Command (no Dokploy):**

```bash
go mod download && go install github.com/swaggo/swag/cmd/swag@latest && /root/go/bin/swag init -g main.go --output ./docs || mkdir -p docs && go build -o chatserver main.go
```

3. **Start Command:**

```bash
./chatserver
```

4. **Redeploy**

---

## 📸 Captura de Tela Esperada

No Dokploy, você deve ver:

```
Settings
├── General
│   ├── Builder: [Dockerfile ▼]  ← MUDE AQUI
│   └── Dockerfile Path: Dockerfile.dokploy
├── Environment Variables
│   ├── MONGODB_URL=mongodb+srv://...
│   ├── MONGODB_DATABASE=sr_robot
│   ├── PORT=8080
│   ├── ENV=production
│   └── GIN_MODE=release
└── Deploy
    └── [Redeploy] ← CLICK AQUI
```

---

## ✅ RESUMO EXECUTIVO

**O QUE FAZER AGORA:**

1. `git push origin main`
2. No Dokploy: Mudar builder para "Dockerfile"
3. Path: `Dockerfile.dokploy`
4. Configurar variáveis de ambiente
5. Redeploy
6. ✅ **PRONTO!**

**TEMPO ESTIMADO:** 5 minutos

---

**Dúvidas?** Consulte `DOKPLOY_FIX.md` para guia detalhado!
