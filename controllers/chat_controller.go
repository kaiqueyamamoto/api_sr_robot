package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"chatserver/database"
	"chatserver/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	N8N_WEBHOOK_URL = "https://galaxy.conecta-tech.com.br/webhook/conversation"
)

// ChatRequest representa a requisição de chat
type ChatRequest struct {
	ConversationID string `json:"conversationId,omitempty"` // Opcional: se não fornecido, cria nova conversa
	Message        string `json:"message" binding:"required"`
}

// ChatResponse representa a resposta do chat
type ChatResponse struct {
	ConversationID string             `json:"conversationId"`
	Message        string             `json:"message"` // Resposta do assistente
	Role           models.MessageRole `json:"role"`
	MessageID      string             `json:"messageId"`
	LatencyMs      int64              `json:"latencyMs"`
}

// N8NRequest representa a requisição para o n8n
type N8NRequest struct {
	Message        string           `json:"message"`
	ConversationID string           `json:"conversationId"`
	History        []models.Message `json:"history,omitempty"` // Histórico das últimas mensagens
}

// N8NResponse representa a resposta do n8n
type N8NResponse struct {
	Output   string                 `json:"output"`   // Campo retornado pelo N8N
	Response string                 `json:"response"` // Alternativa (compatibilidade)
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// GetResponse retorna a resposta (output ou response)
func (n *N8NResponse) GetResponse() string {
	if n.Output != "" {
		return n.Output
	}
	return n.Response
}

// ChatController gerencia as conversas
type ChatController struct {
	conversationsCollection *mongo.Collection
	messagesCollection      *mongo.Collection
}

// NewChatController cria uma nova instância do controller
func NewChatController() *ChatController {
	return &ChatController{
		conversationsCollection: database.GetCollection("conversations"),
		messagesCollection:      database.GetCollection("messages"),
	}
}

// SendMessage godoc
// @Summary      Enviar mensagem para o chatbot
// @Description  Envia uma mensagem e recebe a resposta do chatbot. Cria nova conversa ou continua existente.
// @Tags         chat
// @Accept       json
// @Produce      json
// @Param        request  body      ChatRequest  true  "Mensagem do usuário"
// @Success      200      {object}  ChatResponse
// @Failure      400      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /api/v1/chat [post]
func (ctrl *ChatController) SendMessage(c *gin.Context) {
	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := context.Background()
	startTime := time.Now()

	// 1. Obter ou criar conversa
	var conversationID primitive.ObjectID
	var err error

	if req.ConversationID != "" {
		conversationID, err = primitive.ObjectIDFromHex(req.ConversationID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID de conversa inválido"})
			return
		}

		// Verificar se a conversa existe
		var conversation models.Conversation
		err = ctrl.conversationsCollection.FindOne(ctx, bson.M{"_id": conversationID}).Decode(&conversation)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Conversa não encontrada"})
			return
		}

		// Atualizar updatedAt
		ctrl.conversationsCollection.UpdateOne(
			ctx,
			bson.M{"_id": conversationID},
			bson.M{"$set": bson.M{"updatedAt": time.Now()}},
		)
	} else {
		// Criar nova conversa
		conversation := models.NewConversation("")
		result, err := ctrl.conversationsCollection.InsertOne(ctx, conversation)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar conversa"})
			return
		}
		conversationID = result.InsertedID.(primitive.ObjectID)
	}

	// 2. Salvar mensagem do usuário
	userMessage := models.NewMessage(conversationID, models.RoleUser, req.Message)
	_, err = ctrl.messagesCollection.InsertOne(ctx, userMessage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao salvar mensagem do usuário"})
		return
	}

	// 3. Buscar histórico recente (últimas 10 mensagens)
	history, err := ctrl.getConversationHistory(ctx, conversationID, 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar histórico"})
		return
	}

	// 4. Chamar webhook do n8n
	n8nRequest := N8NRequest{
		Message:        req.Message,
		ConversationID: conversationID.Hex(),
		History:        history,
	}

	n8nResponse, err := ctrl.callN8NWebhook(n8nRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Erro ao chamar n8n: %v", err)})
		return
	}

	// 5. Calcular latência
	latencyMs := time.Since(startTime).Milliseconds()

	// 6. Salvar resposta do assistente
	botResponse := n8nResponse.GetResponse()
	assistantMessage := models.NewMessage(conversationID, models.RoleAssistant, botResponse)
	assistantMessage.LatencyMs = latencyMs
	assistantMessage.Metadata = n8nResponse.Metadata

	result, err := ctrl.messagesCollection.InsertOne(ctx, assistantMessage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao salvar resposta do assistente"})
		return
	}

	// 7. Retornar resposta
	response := ChatResponse{
		ConversationID: conversationID.Hex(),
		Message:        botResponse,
		Role:           models.RoleAssistant,
		MessageID:      result.InsertedID.(primitive.ObjectID).Hex(),
		LatencyMs:      latencyMs,
	}

	c.JSON(http.StatusOK, response)
}

// GetConversationHistory godoc
// @Summary      Obter histórico de conversa
// @Description  Retorna todas as mensagens de uma conversa específica
// @Tags         chat
// @Accept       json
// @Produce      json
// @Param        id    path      string  true  "Conversation ID"
// @Success      200   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]string
// @Failure      404   {object}  map[string]string
// @Router       /api/v1/conversations/{id} [get]
func (ctrl *ChatController) GetConversationHistory(c *gin.Context) {
	conversationID := c.Param("id")

	objectID, err := primitive.ObjectIDFromHex(conversationID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de conversa inválido"})
		return
	}

	ctx := context.Background()

	// Verificar se a conversa existe
	var conversation models.Conversation
	err = ctrl.conversationsCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&conversation)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Conversa não encontrada"})
		return
	}

	// Buscar todas as mensagens da conversa
	messages, err := ctrl.getConversationHistory(ctx, objectID, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar mensagens"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"conversation": conversation,
		"messages":     messages,
	})
}

// ListConversations godoc
// @Summary      Listar conversas
// @Description  Lista todas as conversas ordenadas por data de atualização
// @Tags         chat
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]string
// @Router       /api/v1/conversations [get]
func (ctrl *ChatController) ListConversations(c *gin.Context) {
	ctx := context.Background()

	// Ordenar por updatedAt descendente (mais recentes primeiro)
	findOptions := options.Find().SetSort(bson.D{{Key: "updatedAt", Value: -1}})

	cursor, err := ctrl.conversationsCollection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar conversas"})
		return
	}
	defer cursor.Close(ctx)

	var conversations []models.Conversation
	if err = cursor.All(ctx, &conversations); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao decodificar conversas"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"conversations": conversations,
		"total":         len(conversations),
	})
}

// UpdateConversationTitle godoc
// @Summary      Atualizar título da conversa
// @Description  Atualiza o título de uma conversa existente
// @Tags         chat
// @Accept       json
// @Produce      json
// @Param        id       path      string  true  "Conversation ID"
// @Param        request  body      map[string]string  true  "Novo título"
// @Success      200      {object}  models.Conversation
// @Failure      400      {object}  map[string]string
// @Failure      404      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /api/v1/conversations/{id} [put]
func (ctrl *ChatController) UpdateConversationTitle(c *gin.Context) {
	conversationID := c.Param("id")

	objectID, err := primitive.ObjectIDFromHex(conversationID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de conversa inválido"})
		return
	}

	var request struct {
		Title string `json:"title" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Título é obrigatório"})
		return
	}

	ctx := context.Background()

	// Atualizar o título
	update := bson.M{
		"$set": bson.M{
			"title":     request.Title,
			"updatedAt": time.Now(),
		},
	}

	result := ctrl.conversationsCollection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": objectID},
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)

	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Conversa não encontrada"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar conversa"})
		return
	}

	var conversation models.Conversation
	if err := result.Decode(&conversation); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao decodificar conversa"})
		return
	}

	c.JSON(http.StatusOK, conversation)
}

// DeleteConversation godoc
// @Summary      Deletar conversa
// @Description  Deleta uma conversa e todas suas mensagens
// @Tags         chat
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Conversation ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/v1/conversations/{id} [delete]
func (ctrl *ChatController) DeleteConversation(c *gin.Context) {
	conversationID := c.Param("id")

	objectID, err := primitive.ObjectIDFromHex(conversationID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de conversa inválido"})
		return
	}

	ctx := context.Background()

	// Verificar se a conversa existe
	var conversation models.Conversation
	err = ctrl.conversationsCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&conversation)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Conversa não encontrada"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar conversa"})
		return
	}

	// Deletar todas as mensagens da conversa
	_, err = ctrl.messagesCollection.DeleteMany(ctx, bson.M{"conversationId": objectID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao deletar mensagens"})
		return
	}

	// Deletar a conversa
	_, err = ctrl.conversationsCollection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao deletar conversa"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Conversa deletada com sucesso"})
}

// getConversationHistory busca o histórico de mensagens de uma conversa
func (ctrl *ChatController) getConversationHistory(ctx context.Context, conversationID primitive.ObjectID, limit int64) ([]models.Message, error) {
	filter := bson.M{"conversationId": conversationID}

	findOptions := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}})
	if limit > 0 {
		findOptions.SetLimit(limit)
	}

	cursor, err := ctrl.messagesCollection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var messages []models.Message
	if err = cursor.All(ctx, &messages); err != nil {
		return nil, err
	}

	// Reverter a ordem para ter as mensagens mais antigas primeiro
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}

// callN8NWebhook chama o webhook do n8n
func (ctrl *ChatController) callN8NWebhook(request N8NRequest) (*N8NResponse, error) {
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(N8N_WEBHOOK_URL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("n8n retornou status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Tentar primeiro como array (formato do N8N)
	var n8nArray []N8NResponse
	if err := json.Unmarshal(body, &n8nArray); err == nil && len(n8nArray) > 0 {
		return &n8nArray[0], nil
	}

	// Se não for array, tentar como objeto direto
	var n8nResponse N8NResponse
	if err := json.Unmarshal(body, &n8nResponse); err != nil {
		return nil, fmt.Errorf("erro ao parsear resposta do n8n: %v (body: %s)", err, string(body))
	}

	return &n8nResponse, nil
}
