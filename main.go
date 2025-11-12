package main

import (
	"log"
	"os"

	"chatserver/controllers"
	"chatserver/database"
	_ "chatserver/docs" // Importa a documentação gerada pelo Swagger
	"chatserver/middleware"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           SR Robot API
// @version         1.0
// @description     API para chatbot SR Robot com JWT authentication e Prometheus metrics
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  support@example.com

// @license.name  MIT
// @license.url   http://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Bearer token (add "Bearer " prefix)

func main() {
	// Carregar variáveis de ambiente
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  Arquivo .env não encontrado, usando variáveis de ambiente do sistema")
	}

	// Obter configurações
	mongoURI := os.Getenv("MONGODB_URL")
	if mongoURI == "" {
		log.Fatal("❌ MONGODB_URL não configurado")
	}

	dbName := os.Getenv("MONGODB_DATABASE")
	if dbName == "" {
		dbName = "sr_robot"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Conectar ao MongoDB
	if err := database.Connect(mongoURI, dbName); err != nil {
		log.Fatalf("❌ Erro ao conectar ao MongoDB: %v", err)
	}
	defer database.Disconnect()

	// Configurar Gin
	if os.Getenv("ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Middleware CORS
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health check
	router.GET("/health", healthCheck)

	// Auth routes
	authController := controllers.NewAuthController(database.Database)
	auth := router.Group("/auth")
	{
		auth.POST("/register", authController.Register)
		auth.POST("/login", authController.Login)
	}

	// Profile routes (protegidas com autenticação)
	profileController := controllers.NewProfileController(database.Database)
	profile := router.Group("/profile")
	profile.Use(middleware.AuthMiddleware())
	{
		profile.GET("", profileController.GetProfile)
		profile.PUT("", profileController.UpdateProfile)
	}

	// Rotas da API
	api := router.Group("/api/v1")
	{
		// Chat routes
		chatController := controllers.NewChatController()

		// Enviar mensagem (criar ou continuar conversa)
		api.POST("/chat", chatController.SendMessage)

		// Buscar histórico de uma conversa
		api.GET("/conversations/:id", chatController.GetConversationHistory)

		// Listar todas as conversas
		api.GET("/conversations", chatController.ListConversations)

		// Atualizar título da conversa
		api.PUT("/conversations/:id", chatController.UpdateConversationTitle)

		// Deletar conversa
		api.DELETE("/conversations/:id", chatController.DeleteConversation)
	}

	// Iniciar servidor
	log.Printf("🚀 Servidor rodando na porta %s", port)
	log.Printf("📖 Documentação Swagger disponível em: http://localhost:%s/swagger/index.html", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("❌ Erro ao iniciar servidor: %v", err)
	}
}

// healthCheck godoc
// @Summary      Health Check
// @Description  Verifica se o servidor está rodando
// @Tags         health
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]string
// @Router       /health [get]
func healthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":  "ok",
		"service": "sr_robot_api",
	})
}
