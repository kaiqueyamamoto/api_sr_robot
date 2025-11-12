# 🗑️ Limpeza do Banco de Dados - Conversas Sem userId

## ⚠️ IMPORTANTE

Após a correção de segurança, conversas antigas no banco de dados **não têm** o campo `userId`.

Isso significa que:
- ❌ Nenhum usuário consegue acessá-las (query filtra por userId)
- ❌ Ocupam espaço no banco
- ❌ Podem causar confusão

## 🎯 Solução: Limpar Conversas Antigas

### Opção 1: Script Automático (Recomendado)

```bash
# Conectar ao MongoDB e executar script
mongo "sua-connection-string/sr_robot" cleanup_conversations.js
```

### Opção 2: MongoDB Compass (Visual)

1. Abra MongoDB Compass
2. Conecte ao seu banco `sr_robot`
3. Abra a collection `conversations`
4. Use o filtro:
   ```json
   { "userId": { "$exists": false } }
   ```
5. Selecione todos e delete

### Opção 3: MongoDB Shell (Manual)

```javascript
// Conectar ao banco
use sr_robot

// Ver quantas conversas sem userId existem
db.conversations.countDocuments({userId: {$exists: false}})

// Listar alguns exemplos
db.conversations.find({userId: {$exists: false}}).limit(5)

// CUIDADO: Esta operação DELETA dados permanentemente!

// 1. Coletar IDs das conversas sem userId
const conversationIds = [];
db.conversations.find({userId: {$exists: false}}).forEach(function(c) {
  conversationIds.push(c._id);
});

print(`Encontradas ${conversationIds.length} conversas para deletar`);

// 2. Deletar mensagens associadas
const msgResult = db.messages.deleteMany({
  conversationId: {$in: conversationIds}
});

print(`Mensagens deletadas: ${msgResult.deletedCount}`);

// 3. Deletar conversas
const convResult = db.conversations.deleteMany({
  userId: {$exists: false}
});

print(`Conversas deletadas: ${convResult.deletedCount}`);
```

## 📊 Verificar Resultado

```javascript
// Deve retornar 0
db.conversations.countDocuments({userId: {$exists: false}})

// Ver todas as conversas restantes (devem ter userId)
db.conversations.find().pretty()
```

## 🔐 Opção Alternativa: Atribuir userId

Se você souber o dono das conversas antigas, pode atribuí-las:

```javascript
// Atribuir todas as conversas sem userId a um usuário específico
db.conversations.updateMany(
  {userId: {$exists: false}},
  {$set: {userId: "USER_ID_DO_PROPRIETARIO"}}
)
```

**Como obter user_id:**
```javascript
// Listar usuários
db.users.find({}, {email: 1, _id: 1})

// Copie o _id do usuário e use no comando acima
```

## ⚠️ ATENÇÃO

- ⛔ **Backup**: Faça backup antes de deletar!
- ⛔ **Irreversível**: Não há como recuperar após deletar
- ⛔ **Produção**: Teste em desenvolvimento primeiro

## 🎯 Quando Executar

Execute este script:
1. ✅ **Antes** de fazer o primeiro deploy da API corrigida
2. ✅ **Depois** de fazer backup do banco
3. ✅ **Uma única vez** (não é necessário repetir)

## 📝 Backup Antes de Deletar

```bash
# MongoDB Atlas
# Faça backup via console (Cloud Backups)

# MongoDB Local
mongodump --uri="your-connection-string" --out=backup-$(date +%Y%m%d)

# Restaurar se necessário
mongorestore --uri="your-connection-string" backup-20251112/sr_robot
```

## ✅ Checklist

- [ ] Backup do banco de dados feito
- [ ] Script testado em desenvolvimento
- [ ] Conversas sem userId identificadas
- [ ] Decisão tomada (deletar ou atribuir)
- [ ] Script executado
- [ ] Verificação feita (count deve ser 0)
- [ ] API com correção deployada
- [ ] Teste com múltiplos usuários realizado

---

**Após executar este script, cada usuário verá APENAS suas próprias conversas!**

