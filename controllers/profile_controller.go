package controllers

import (
	"context"
	"net/http"
	"time"

	"chatserver/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProfileController struct {
	userCollection *mongo.Collection
}

func NewProfileController(db *mongo.Database) *ProfileController {
	return &ProfileController{
		userCollection: db.Collection("users"),
	}
}

// GetProfile godoc
// @Summary      Obter perfil do usuário
// @Description  Retorna as informações de perfil do usuário autenticado (nome, bio, email)
// @Tags         profile
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  models.ProfileResponse
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /profile [get]
func (pc *ProfileController) GetProfile(c *gin.Context) {
	// Obter user_id do contexto (setado pelo middleware de auth)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	// Converter string para ObjectID
	objectID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de usuário inválido"})
		return
	}

	// Buscar usuário no banco
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User
	err = pc.userCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Usuário não encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar usuário"})
		return
	}

	// Retornar profile (com valores nulos se não existirem)
	c.JSON(http.StatusOK, models.ProfileResponse{
		Email: user.Email,
		Name:  user.Name,
		Bio:   user.Bio,
	})
}

// UpdateProfile godoc
// @Summary      Atualizar perfil do usuário
// @Description  Atualiza nome e bio do usuário autenticado (email não pode ser alterado)
// @Tags         profile
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      models.UpdateProfileRequest  true  "Dados do perfil"
// @Success      200      {object}  models.ProfileResponse
// @Failure      400      {object}  map[string]string
// @Failure      401      {object}  map[string]string
// @Failure      404      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /profile [put]
func (pc *ProfileController) UpdateProfile(c *gin.Context) {
	// Obter user_id do contexto
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	// Converter string para ObjectID
	objectID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de usuário inválido"})
		return
	}

	// Parse request body
	var req models.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Preparar update
	update := bson.M{
		"$set": bson.M{
			"updated_at": time.Now(),
		},
	}

	// Adicionar campos ao update apenas se foram fornecidos
	if req.Name != nil {
		update["$set"].(bson.M)["name"] = req.Name
	}
	if req.Bio != nil {
		update["$set"].(bson.M)["bio"] = req.Bio
	}

	// Atualizar no banco
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result := pc.userCollection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": objectID},
		update,
		// Opções para retornar o documento atualizado
	)

	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Usuário não encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar perfil"})
		return
	}

	// Buscar usuário atualizado
	var updatedUser models.User
	err = pc.userCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&updatedUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar usuário atualizado"})
		return
	}

	// Retornar perfil atualizado
	c.JSON(http.StatusOK, models.ProfileResponse{
		Email: updatedUser.Email,
		Name:  updatedUser.Name,
		Bio:   updatedUser.Bio,
	})
}

