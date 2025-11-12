package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MessageRole define os papéis possíveis em uma mensagem
type MessageRole string

const (
	RoleUser      MessageRole = "user"
	RoleAssistant MessageRole = "assistant"
	RoleSystem    MessageRole = "system"
)

// Message representa uma mensagem dentro de uma conversa
type Message struct {
	ID             primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ConversationID primitive.ObjectID `json:"conversationId" bson:"conversationId"`
	Role           MessageRole        `json:"role" bson:"role"`                     // user, assistant, system
	Content        string             `json:"content" bson:"content"`               // Conteúdo da mensagem
	Tokens         int                `json:"tokens,omitempty" bson:"tokens,omitempty"` // Quantidade de tokens (opcional)
	LatencyMs      int64              `json:"latencyMs,omitempty" bson:"latencyMs,omitempty"` // Latência da resposta em ms
	Metadata       map[string]interface{} `json:"metadata,omitempty" bson:"metadata,omitempty"` // Metadados adicionais
	CreatedAt      time.Time          `json:"createdAt" bson:"createdAt"`
}

// NewMessage cria uma nova mensagem
func NewMessage(conversationID primitive.ObjectID, role MessageRole, content string) *Message {
	return &Message{
		ID:             primitive.NewObjectID(),
		ConversationID: conversationID,
		Role:           role,
		Content:        content,
		CreatedAt:      time.Now(),
	}
}

