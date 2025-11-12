# Dev Container for SR Robot API

Este Dev Container fornece um ambiente de desenvolvimento completo para a API Go do SR Robot.

## 🚀 O que está incluído

### Serviços
- **Go API** - Ambiente de desenvolvimento Go 1.23
- **MongoDB Atlas** - Banco de dados NoSQL na nuvem (conexão direta) [[memory:10890365]]

### Ferramentas Go
- `gopls` - Language server
- `dlv` - Debugger
- `golangci-lint` - Linter
- `goimports` - Organização de imports
- `air` - Hot reload
- `swag` - Geração de documentação Swagger

### Extensões VS Code
- Go (oficial)
- GitHub Copilot
- Docker
- GitLens
- Thunder Client (teste de APIs)
- MongoDB para VS Code

## 📦 Como usar

### 1. Pré-requisitos
- Docker Desktop instalado e rodando
- VS Code com extensão "Dev Containers" instalada

### 2. Abrir no Dev Container

**Opção A: Via Command Palette**
1. Abra a pasta `api` no VS Code
2. Pressione `F1` ou `Ctrl+Shift+P`
3. Digite: `Dev Containers: Reopen in Container`
4. Aguarde a build e inicialização (primeira vez pode demorar)

**Opção B: Via notificação**
1. Abra a pasta `api` no VS Code
2. Clique em "Reopen in Container" quando aparecer a notificação

### 3. Verificar se está funcionando

Após o container iniciar, abra um terminal integrado e execute:

```bash
# Verificar Go
go version

# Verificar conexão com MongoDB Atlas
mongosh "mongodb+srv://sr_robot:brBBTUbOqnxVpN0S@conecta-tech.pajxycn.mongodb.net/?appName=Conecta-Tech" --eval "db.version()"
```

## 🛠️ Comandos úteis

### Desenvolvimento

```bash
# Instalar dependências
make deps

# Rodar com hot reload
make dev

# Build
make build

# Rodar testes
make test

# Rodar testes com cobertura
make test-cover

# Formatar código
make format

# Rodar linters
make lint
```

### Banco de dados

```bash
# Conectar ao MongoDB Atlas
mongosh "mongodb+srv://sr_robot:brBBTUbOqnxVpN0S@conecta-tech.pajxycn.mongodb.net/sr_robot?appName=Conecta-Tech"

# Ver collections
mongosh "mongodb+srv://sr_robot:brBBTUbOqnxVpN0S@conecta-tech.pajxycn.mongodb.net/sr_robot?appName=Conecta-Tech" --eval "db.getCollectionNames()"
```

### Verificar serviços

```bash
# Status do container da API
docker ps

# Logs da API
docker logs <container_id>
```

## 🗄️ Estrutura do banco de dados

O MongoDB Atlas está configurado com:
- ✅ Database `sr_robot` no MongoDB Atlas
- ✅ Conexão direta via connection string
- ✅ Collections: users, profiles, sources, chunks, conversations, messages

Para inicializar as collections localmente, você pode executar o script `init-scripts/01-init-mongodb.js` manualmente no MongoDB Atlas via Mongo Shell.

## 🔌 Portas expostas

| Serviço | Porta | URL                   |
|---------|-------|-----------------------|
| API     | 8080  | http://localhost:8080 |

## 🌍 Variáveis de ambiente

As seguintes variáveis estão pré-configuradas no Dev Container:

```bash
MONGODB_URL=mongodb://localhost:27017/sr_robot
MONGODB_DATABASE=sr_robot
GO111MODULE=on
GOPATH=/go
```

Para adicionar mais variáveis, edite `devcontainer.json` na seção `remoteEnv`.

## 🐛 Debug

O Dev Container está configurado para debug com Delve:

1. Adicione breakpoints no código
2. Pressione `F5` ou vá em "Run and Debug"
3. Selecione "Launch Package" ou "Attach to Process"

## 📝 Hot Reload

O projeto usa [Air](https://github.com/cosmtrek/air) para hot reload:

```bash
# Inicia servidor com hot reload
make dev

# Ou diretamente
air
```

Configuração em `.air.toml`.

## 🔧 Personalização

### Adicionar extensões VS Code

Edite `.devcontainer/devcontainer.json`:

```json
"extensions": [
  "golang.go",
  "sua-extensao-aqui"
]
```

### Adicionar ferramentas Go

Edite `.devcontainer/Dockerfile`:

```dockerfile
RUN go install -v seu-pacote@latest
```

### Adicionar serviços Docker

Edite `.devcontainer/docker-compose.yml`:

```yaml
services:
  seu-servico:
    image: imagem:tag
    # ...
```

## 🚨 Troubleshooting

### Container não inicia

```bash
# Limpar containers e volumes
docker-compose -f .devcontainer/docker-compose.yml down -v

# Rebuild do container
F1 > Dev Containers: Rebuild Container
```

### MongoDB Atlas não conecta

```bash
# Testar conexão com MongoDB Atlas
mongosh "mongodb+srv://sr_robot:brBBTUbOqnxVpN0S@conecta-tech.pajxycn.mongodb.net/?appName=Conecta-Tech" --eval "db.adminCommand('ping')"

# Verificar se há problemas de firewall/rede
# Certifique-se que seu IP está na whitelist do MongoDB Atlas
```

### Go modules com erro

```bash
# Limpar cache
go clean -modcache

# Re-download
go mod download
go mod tidy
```

### Permissões no Windows

Se tiver problemas com permissões de arquivo no Windows:

1. Certifique-se que o Docker Desktop está usando WSL 2
2. Clone o repositório dentro do WSL (não em `/mnt/c/`)

## 📚 Recursos

- [Dev Containers Documentation](https://code.visualstudio.com/docs/devcontainers/containers)
- [Go Dev Container](https://github.com/devcontainers/templates/tree/main/src/go)
- [MongoDB Documentation](https://www.mongodb.com/docs/)
- [MongoDB Go Driver](https://www.mongodb.com/docs/drivers/go/current/)
- [Air Documentation](https://github.com/cosmtrek/air)

## 🎯 Próximos passos

Após o ambiente estar rodando:

1. Inicialize o Go module: `make init`
2. Instale as dependências: `make deps`
3. Rode os testes: `make test`
4. Inicie o servidor: `make dev`
5. Acesse: http://localhost:8080

---

**Dica**: Use `make help` para ver todos os comandos disponíveis!

