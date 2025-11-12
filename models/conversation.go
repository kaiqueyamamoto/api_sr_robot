package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Conversation representa uma conversa entre o usuário e o chatbot
type Conversation struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID    string             `json:"userId,omitempty" bson:"userId,omitempty"` // Opcional: para usuários autenticados
	Title     string             `json:"title,omitempty" bson:"title,omitempty"`   // Título da conversa (pode ser gerado automaticamente)
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
}

// NewConversation cria uma nova conversa
func NewConversation(userID string) *Conversation {
	now := time.Now()
	return &Conversation{
		ID:        primitive.NewObjectID(),
		UserID:    userID,
		Title:     "Nova Conversa",
		CreatedAt: now,
		UpdatedAt: now,
	}
}

