# 🚀 DEPLOY DOKPLOY - PASSO A PASSO COMPLETO

## ❌ Erro Atual:
```
Error: Writing app
Caused by: Is a directory (os error 21)
```

**CAUSA:** Dokploy está usando **Nixpacks** (que tem o bug)

**SOLUÇÃO:** Mudar para **Dockerfile** ✅

---

## 📋 PASSO A PASSO COMPLETO

### ✅ Etapa 1: Verificar que o código foi enviado

Você já fez isso! ✅ O Dokploy clonou com sucesso:
```
Cloned github.com/kaiqueyamamoto/api_sr_robot.git: ✅
```

---

### ⚙️ Etapa 2: MUDAR BUILDER NO DOKPLOY

**ESTE É O PASSO CRUCIAL!**

#### 2.1. Acesse o Dokploy Dashboard
```
https://seu-dokploy.com/dashboard
```

#### 2.2. Localize seu Aplicativo
- Nome: **srrobot-api-ejhf6d** (ou similar)
- Status: ❌ Build Failed

#### 2.3. Entre nas Configurações

1. **Click no seu app** → `srrobot-api-ejhf6d`
2. **Click em "Settings"** (no menu lateral)
3. **Click em "General"**

#### 2.4. MUDAR O BUILDER ⭐

Na seção **"Build Configuration"**:

```
┌─────────────────────────────────────┐
│ Build Configuration                 │
├─────────────────────────────────────┤
│ Builder: [Nixpacks ▼]      ← AQUI! │
│                                     │
│ Mude para:                          │
│ Builder: [Dockerfile ▼]    ← ISTO! │
│                                     │
│ Dockerfile Path:                    │
│ [Dockerfile.dokploy]       ← ISTO! │
└─────────────────────────────────────┘
```

**Passos exatos:**
1. Click no dropdown **"Builder"**
2. Selecione **"Dockerfile"**
3. No campo **"Dockerfile Path"**, digite: `Dockerfile.dokploy`
4. **Scroll para baixo** e click **"Save Changes"**

---

### 🔐 Etapa 3: Configurar Variáveis de Ambiente

Na mesma página de Settings:

1. **Click em "Environment"** (no menu lateral)
2. **Add Variable** (para cada uma abaixo):

```env
Nome: MONGODB_URL
Valor: mongodb+srv://sr_robot:brBBTUbOqnxVpN0S@conecta-tech.pajxycn.mongodb.net/

Nome: MONGODB_DATABASE
Valor: sr_robot

Nome: PORT
Valor: 8080

Nome: ENV
Valor: production

Nome: GIN_MODE
Valor: release

Nome: JWT_SECRET
Valor: seu_secret_super_seguro_aqui_123456
```

**Importante:** Use um JWT_SECRET forte em produção!

3. **Click "Save"** depois de adicionar todas

---

### 🚀 Etapa 4: REDEPLOY

1. **Volte para a página principal do app**
2. **Click no botão "Redeploy"** ou **"Rebuild"**
   - Pode estar no canto superior direito
   - Ou no menu de ações

3. **Aguarde o build** (2-5 minutos)

---

## 📊 O que vai acontecer (BUILD BEM-SUCEDIDO):

```
✅ Cloning Repo github.com/kaiqueyamamoto/api_sr_robot.git
✅ Build with Dockerfile: Dockerfile.dokploy
✅ [1/4] Building stage: builder
✅ [2/4] Installing Go and dependencies
✅ [3/4] Generating Swagger documentation
✅ [4/4] Building application binary
✅ Creating final image
✅ Starting container
✅ Container running
✅ Health check passed
```

**Nos logs você verá:**
```
✅ Conectado ao MongoDB Atlas!
🚀 Servidor rodando na porta 8080
📖 Documentação Swagger disponível em: http://localhost:8080/swagger/index.html
```

---

## ✅ Etapa 5: Verificar Deploy

### 5.1. Obter a URL do App

No Dokploy, você verá algo como:
```
https://srrobot-api-ejhf6d.your-domain.com
```

### 5.2. Testar Endpoints

#### Health Check:
```bash
curl https://sua-url.com/health
```

**Resposta esperada:**
```json
{
  "status": "ok",
  "service": "sr_robot_api"
}
```

#### Swagger UI:
```
https://sua-url.com/
```

Deve redirecionar automaticamente para:
```
https://sua-url.com/swagger/index.html
```

#### Test Chat:
```bash
curl -X POST https://sua-url.com/api/v1/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "Olá, SR Robot!"
  }'
```

---

## 🎯 CHECKLIST RÁPIDO

Marque cada item conforme completar:

- [ ] **Passo 1:** Código enviado para GitHub ✅ (JÁ FEITO!)
- [ ] **Passo 2:** Acessar Dokploy Dashboard
- [ ] **Passo 3:** Ir em Settings → General
- [ ] **Passo 4:** Mudar Builder de "Nixpacks" para "Dockerfile"
- [ ] **Passo 5:** Definir Dockerfile Path: `Dockerfile.dokploy`
- [ ] **Passo 6:** Salvar mudanças
- [ ] **Passo 7:** Ir em Settings → Environment
- [ ] **Passo 8:** Adicionar todas as variáveis de ambiente
- [ ] **Passo 9:** Salvar variáveis
- [ ] **Passo 10:** Click em "Redeploy"
- [ ] **Passo 11:** Aguardar build completar (2-5 min)
- [ ] **Passo 12:** Verificar logs (deve mostrar "Conectado ao MongoDB")
- [ ] **Passo 13:** Testar /health endpoint
- [ ] **Passo 14:** Acessar Swagger UI
- [ ] **Passo 15:** Testar chat endpoint

---

## 🖼️ VISUAL GUIDE

### Onde está o Builder?

```
Dokploy Dashboard
└── Projects
    └── Seu Projeto
        └── srrobot-api-ejhf6d
            └── Settings (menu lateral)
                └── General (tab)
                    └── Build Configuration
                        ├── Builder: [Dockerfile ▼]  ← MUDE AQUI!
                        └── Dockerfile Path: Dockerfile.dokploy
```

---

## 🆘 TROUBLESHOOTING

### Se ainda der erro:

#### Opção A: Verificar Dockerfile
```bash
# Certifique-se que o arquivo existe:
ls -la Dockerfile.dokploy

# Deve aparecer no resultado
```

#### Opção B: Usar Build Manual

Se por algum motivo o Dockerfile não funcionar:

1. **Volte para Nixpacks** (ou use "Custom Build")
2. **Delete o arquivo** `nixpacks.toml` do repositório
3. **Configure manualmente no Dokploy:**

**Build Command:**
```bash
go mod download && \
go install github.com/swaggo/swag/cmd/swag@latest && \
/root/go/bin/swag init -g main.go --output ./docs || mkdir -p docs && \
go build -o chatserver main.go
```

**Start Command:**
```bash
./chatserver
```

4. **Redeploy**

---

## 📞 SUPORTE

### Logs do Dokploy

Para ver logs detalhados:
1. No app, click em **"Logs"** ou **"Build Logs"**
2. Procure por linhas com ❌ ou "Error"
3. Se ver "Connected to MongoDB" = ✅ sucesso!

### Variáveis de Ambiente

Verifique se todas estão configuradas:
```
Settings → Environment → [Lista de variáveis]
```

Deve ter **7 variáveis** no mínimo:
- MONGODB_URL
- MONGODB_DATABASE
- PORT
- ENV
- GIN_MODE
- JWT_SECRET
- (outras opcionais)

---

## ✅ RESUMO EXECUTIVO

**O PROBLEMA:**
- ❌ Dokploy está usando Nixpacks (tem bug)

**A SOLUÇÃO:**
- ✅ Mudar para Dockerfile

**ONDE MUDAR:**
- Dokploy → Seu App → Settings → General → Builder → **"Dockerfile"**

**O QUE DIGITAR:**
- Dockerfile Path: `Dockerfile.dokploy`

**DEPOIS:**
- Adicionar variáveis de ambiente
- Click "Redeploy"
- ✅ **Vai funcionar!**

**TEMPO TOTAL:** 5-10 minutos

---

## 🎉 SUCESSO!

Quando funcionar, você verá:

```
✅ Build successful
✅ Container running
✅ Health check passed

Access your app at:
https://sua-url.com
```

**Swagger disponível em:**
```
https://sua-url.com/swagger/index.html
```

---

**Qualquer dúvida, consulte:**
- `DOKPLOY_FIX.md` - Troubleshooting detalhado
- `SOLUCAO_DOKPLOY.md` - Guia visual
- `Dockerfile.dokploy` - O arquivo que será usado

**BOA SORTE! 🚀**

