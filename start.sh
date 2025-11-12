#!/bin/bash

# Script de inicialização do SR Robot API
# Compila e inicia o servidor com todas as configurações

set -e

echo "🤖 SR Robot API - Inicialização"
echo "================================"
echo ""

# Cores para output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 1. Verificar se .env existe
if [ ! -f .env ]; then
    echo -e "${YELLOW}⚠️  Arquivo .env não encontrado${NC}"
    echo "   Usando variáveis de ambiente do sistema..."
fi

# 2. Verificar variáveis necessárias
if [ -z "$MONGODB_URL" ]; then
    echo -e "${YELLOW}⚠️  MONGODB_URL não configurado${NC}"
    echo "   Defina a variável ou crie um arquivo .env"
fi

# 3. Parar servidor anterior (se existir)
echo -e "${BLUE}🛑 Parando servidor anterior...${NC}"
lsof -ti:8080 | xargs kill -9 2>/dev/null || true
sleep 1

# 4. Gerar documentação Swagger
echo -e "${BLUE}📖 Gerando documentação Swagger...${NC}"
/go/bin/swag init -g main.go --output ./docs

# 5. Compilar projeto
echo -e "${BLUE}🔨 Compilando projeto...${NC}"
go build -o chatserver main.go

# 6. Iniciar servidor
echo ""
echo -e "${GREEN}✅ Iniciando servidor...${NC}"
echo ""
./chatserver &

# 7. Aguardar servidor iniciar
sleep 5

# 8. Verificar se está rodando
if curl -s http://localhost:8080/health > /dev/null; then
    echo ""
    echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${GREEN}✅ Servidor rodando com sucesso!${NC}"
    echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo ""
    echo -e "${BLUE}📍 Endpoints disponíveis:${NC}"
    echo "   • API: http://localhost:8080"
    echo "   • Health: http://localhost:8080/health"
    echo -e "   • ${GREEN}Swagger: http://localhost:8080/swagger/index.html${NC}"
    echo ""
    echo -e "${BLUE}📚 Documentação:${NC}"
    echo "   • SWAGGER.md - Guia do Swagger"
    echo "   • API_EXAMPLES.md - Exemplos da API"
    echo "   • README.md - Documentação geral"
    echo ""
else
    echo ""
    echo -e "${YELLOW}⚠️  Servidor não respondeu${NC}"
    echo "   Verifique os logs em /tmp/chatserver.log"
    echo ""
fi

