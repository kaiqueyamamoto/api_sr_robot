package controllers

import (
	"context"
	"net/http"
	"time"

	"chatserver/metrics"
	"chatserver/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var jwtSecret = []byte("your-secret-key-change-this-in-production") // Change this to env variable in production

type AuthController struct {
	userCollection *mongo.Collection
}

func NewAuthController(db *mongo.Database) *AuthController {
	return &AuthController{
		userCollection: db.Collection("users"),
	}
}

// Register godoc
// @Summary      Registrar novo usuário
// @Description  Cria um novo usuário com email e senha
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      models.RegisterRequest  true  "Dados de registro"
// @Success      201      {object}  models.AuthResponse
// @Failure      400      {object}  map[string]string
// @Failure      409      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /auth/register [post]
func (ac *AuthController) Register(c *gin.Context) {
	var req models.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user already exists
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	start := time.Now()
	var existingUser models.User
	err := ac.userCollection.FindOne(ctx, bson.M{"email": req.Email}).Decode(&existingUser)
	metrics.RecordDatabaseOperation("find", "users", "success", time.Since(start).Seconds())

	if err == nil {
		metrics.RecordAuthAttempt("register", "failure")
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	// Create new user
	now := time.Now()
	user := models.User{
		Email:     req.Email,
		Password:  req.Password,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Hash password
	if err := user.HashPassword(); err != nil {
		metrics.RecordAuthAttempt("register", "failure")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Insert user into database
	start = time.Now()
	result, err := ac.userCollection.InsertOne(ctx, user)
	if err != nil {
		metrics.RecordDatabaseOperation("insert", "users", "failure", time.Since(start).Seconds())
		metrics.RecordAuthAttempt("register", "failure")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}
	metrics.RecordDatabaseOperation("insert", "users", "success", time.Since(start).Seconds())

	// Get the inserted ID
	insertedID := result.InsertedID.(primitive.ObjectID).Hex()

	// Generate JWT token
	token, err := generateToken(req.Email, insertedID)
	if err != nil {
		metrics.RecordAuthAttempt("register", "failure")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Record successful registration
	metrics.RecordAuthAttempt("register", "success")
	metrics.RecordTokenIssued()

	c.JSON(http.StatusCreated, models.AuthResponse{
		Token:     token,
		Email:     req.Email,
		UserID:    insertedID,
		CreatedAt: user.CreatedAt,
	})
}

// Login godoc
// @Summary      Login de usuário
// @Description  Autentica um usuário e retorna token JWT (válido por 24 horas)
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      models.LoginRequest  true  "Credenciais de login"
// @Success      200      {object}  models.AuthResponse
// @Failure      400      {object}  map[string]string
// @Failure      401      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /auth/login [post]
func (ac *AuthController) Login(c *gin.Context) {
	var req models.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find user by email
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	start := time.Now()
	var user models.User
	err := ac.userCollection.FindOne(ctx, bson.M{"email": req.Email}).Decode(&user)

	if err != nil {
		metrics.RecordDatabaseOperation("find", "users", "failure", time.Since(start).Seconds())
		metrics.RecordAuthAttempt("login", "failure")
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	metrics.RecordDatabaseOperation("find", "users", "success", time.Since(start).Seconds())

	// Check password
	if err := user.CheckPassword(req.Password); err != nil {
		metrics.RecordAuthAttempt("login", "failure")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Generate JWT token
	token, err := generateToken(user.Email, user.ID)
	if err != nil {
		metrics.RecordAuthAttempt("login", "failure")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Record successful login
	metrics.RecordAuthAttempt("login", "success")
	metrics.RecordTokenIssued()

	c.JSON(http.StatusOK, models.AuthResponse{
		Token:     token,
		Email:     user.Email,
		UserID:    user.ID,
		CreatedAt: user.CreatedAt,
	})
}

// generateToken generates a JWT token with 1 day expiration
func generateToken(email, userID string) (string, error) {
	claims := models.Claims{
		Email:  email,
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // 1 day expiration
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ValidateToken validates a JWT token
func ValidateToken(tokenString string) (*models.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*models.Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}
