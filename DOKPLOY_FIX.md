# 🔧 Correção do Erro "Is a directory" - Dokploy

## ❌ Erro Atual

```
Error: Writing app
Caused by: Is a directory (os error 21)
```

## ✅ 3 Soluções Possíveis

### **Solução 1: Usar Dockerfile (MAIS CONFIÁVEL)**

No painel do Dokploy:

1. **Mude o Builder** de "Nixpacks" para **"Dockerfile"**
2. **Dockerfile Path:** `Dockerfile.dokploy`
3. **Configure variáveis de ambiente:**
   ```
   MONGODB_URL=mongodb+srv://...
   MONGODB_DATABASE=sr_robot
   PORT=8080
   ENV=production
   GIN_MODE=release
   ```
4. **Redeploy**

✅ Esta é a solução mais estável!

---

### **Solução 2: Build Manual no Dokploy**

No painel do Dokploy, configure manualmente:

**Build Command:**
```bash
go mod download && go install github.com/swaggo/swag/cmd/swag@latest && /root/go/bin/swag init -g main.go --output ./docs || mkdir -p docs && go build -o chatserver main.go
```

**Start Command:**
```bash
./chatserver
```

**Delete o arquivo** `nixpacks.toml` ou renomeie para `nixpacks.toml.bak`

---

### **Solução 3: Usar build.sh (Script)**

1. **Renomear configuração:**
```bash
mv nixpacks.toml nixpacks-backup.toml
mv nixpacks-simple.toml nixpacks.toml
```

2. **Commit:**
```bash
git add .
git commit -m "fix: use build script for Dokploy"
git push origin main
```

3. **Redeploy no Dokploy**

---

## 🚀 Solução Recomendada (Dockerfile)

### Passo a Passo:

#### 1. No Dokploy Dashboard

- Vá para seu aplicativo
- **Settings** → **General**
- **Builder:** Selecione **"Dockerfile"**
- **Dockerfile Path:** `Dockerfile.dokploy`
- **Save**

#### 2. Configurar Variáveis

**Environment Variables:**
```
MONGODB_URL=mongodb+srv://sr_robot:brBBTUbOqnxVpN0S@conecta-tech.pajxycn.mongodb.net/
MONGODB_DATABASE=sr_robot
PORT=8080
ENV=production
GIN_MODE=release
```

#### 3. Deploy

- Click em **"Redeploy"** ou **"Rebuild"**
- Aguarde o build completar (2-5 minutos)
- ✅ Deve funcionar!

### Por que Dockerfile é melhor?

✅ Mais controle sobre o processo de build  
✅ Multi-stage build (imagem menor)  
✅ Sem dependência de paths específicos  
✅ Funciona em qualquer plataforma  
✅ Mais confiável  

---

## 📝 Arquivos Disponíveis

| Arquivo | Uso |
|---------|-----|
| **Dockerfile.dokploy** | ✅ RECOMENDADO para Dokploy |
| **nixpacks.toml** | Versão simplificada |
| **nixpacks-simple.toml** | Com build.sh |
| **build.sh** | Script de build standalone |
| **Dockerfile** | Docker padrão |

---

## 🔍 Debug do Erro

O erro "Is a directory (os error 21)" no Nixpacks geralmente ocorre por:

1. ❌ Conflito de configuração múltipla
2. ❌ Path incorreto do swag
3. ❌ Problema com variável $GOPATH
4. ❌ Permissões de diretório

**Solução:** Usar Dockerfile elimina todos esses problemas!

---

## ✅ Teste Local do Dockerfile

Antes de fazer deploy, teste localmente:

```bash
# Build
docker build -f Dockerfile.dokploy -t sr-robot-api .

# Test run
docker run -p 8080:8080 \
  -e MONGODB_URL="mongodb+srv://..." \
  -e MONGODB_DATABASE="sr_robot" \
  -e PORT="8080" \
  sr-robot-api

# Test
curl http://localhost:8080/health
```

---

## 🎯 Checklist de Deploy

- [ ] Mudar builder para "Dockerfile" no Dokploy
- [ ] Definir Dockerfile path: `Dockerfile.dokploy`
- [ ] Configurar variáveis de ambiente
- [ ] Redeploy
- [ ] Verificar logs
- [ ] Testar health check
- [ ] Acessar Swagger

---

## 📊 Resultado Esperado

```
✅ Building with Dockerfile
✅ Stage 1: Builder (installing deps)
✅ Stage 2: Final image (copying binary)
✅ Container created
✅ Application started
✅ Health check passed
```

Então você verá nos logs:
```
✅ Conectado ao MongoDB Atlas!
🚀 Servidor rodando na porta 8080
📖 Documentação Swagger disponível em: http://localhost:8080/swagger/index.html
```

---

## 💡 Recomendação Final

**USE O DOCKERFILE!**

1. Mude para builder "Dockerfile" no Dokploy
2. Path: `Dockerfile.dokploy`
3. Redeploy
4. ✅ Funciona!

É a solução mais confiável e compatível! 🚀

