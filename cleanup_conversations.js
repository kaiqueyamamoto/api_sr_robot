// Script para limpar conversas sem userId do banco de dados
// Execute com: mongo <connection_string>/sr_robot cleanup_conversations.js

// Conectar ao banco (já conectado se executado via mongo cli)

print("🔍 Verificando conversas sem userId...");

// Contar conversas sem userId
const countWithoutUserId = db.conversations.countDocuments({
  userId: { $exists: false }
});

print(`📊 Encontradas ${countWithoutUserId} conversas sem userId`);

if (countWithoutUserId === 0) {
  print("✅ Nenhuma conversa sem userId encontrada. Banco de dados já está limpo!");
  quit();
}

// Listar alguns exemplos
print("\n📋 Exemplos de conversas sem userId:");
db.conversations.find({ userId: { $exists: false } }).limit(5).forEach(function(conv) {
  print(`  - ID: ${conv._id}, Título: ${conv.title}, Criado em: ${conv.createdAt}`);
});

print("\n⚠️  ATENÇÃO: Estas conversas serão DELETADAS!");
print("⚠️  Esta operação NÃO pode ser desfeita!");
print("\n");

// Perguntar confirmação (se executado interativamente)
// Para execução automática, comente o bloco abaixo

/*
const confirmation = readline().trim().toLowerCase();
if (confirmation !== 'sim') {
  print("❌ Operação cancelada pelo usuário.");
  quit();
}
*/

print("🗑️  Iniciando limpeza...");

// Coletar IDs das conversas que serão deletadas
const conversationIds = [];
db.conversations.find({ userId: { $exists: false } }).forEach(function(conv) {
  conversationIds.push(conv._id);
});

print(`📦 ${conversationIds.length} IDs de conversas coletados`);

// Deletar mensagens associadas
print("🗑️  Deletando mensagens associadas...");
const messagesResult = db.messages.deleteMany({
  conversationId: { $in: conversationIds }
});
print(`✅ ${messagesResult.deletedCount} mensagens deletadas`);

// Deletar conversas
print("🗑️  Deletando conversas...");
const conversationsResult = db.conversations.deleteMany({
  userId: { $exists: false }
});
print(`✅ ${conversationsResult.deletedCount} conversas deletadas`);

// Verificar resultado
const remainingWithoutUserId = db.conversations.countDocuments({
  userId: { $exists: false }
});

print("\n" + "=".repeat(50));
print("📊 RESUMO DA LIMPEZA:");
print("=".repeat(50));
print(`Conversas deletadas: ${conversationsResult.deletedCount}`);
print(`Mensagens deletadas: ${messagesResult.deletedCount}`);
print(`Conversas sem userId restantes: ${remainingWithoutUserId}`);
print("=".repeat(50));

if (remainingWithoutUserId === 0) {
  print("✅ Limpeza concluída com sucesso!");
  print("✅ Banco de dados está seguro agora.");
} else {
  print("⚠️  Ainda existem conversas sem userId. Execute novamente se necessário.");
}

print("\n🔐 Agora todas as conversas novas terão userId associado.");
print("🔐 Usuários só poderão ver suas próprias conversas.");

