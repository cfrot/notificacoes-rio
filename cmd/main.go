package main

import (
	"log"
	"os"

	"notificacoes-rio/internal/database"
	"notificacoes-rio/internal/handler"
	"notificacoes-rio/internal/repository"
	"notificacoes-rio/internal/service"

	"github.com/gin-gonic/gin"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("erro ao carregar .env")
	}

	db := database.Conectar()
	defer db.Close()

	repo := repository.NewNotificacaoRepository(db)
	service := service.NewNotificacaoService(repo)

	webhookSecret := os.Getenv("WEBHOOK_SECRET")

	handler := handler.NewNotificacaoHandler(service, webhookSecret)

	router := gin.Default()

	router.POST("/webhook", handler.Webhook)

	router.GET("/notificacoes", handler.Listar)
	router.PATCH("/notificacoes/:id/read", handler.MarcarComoLida)
	router.GET("/notificacoes/unread-count", handler.UnreadCount)

	router.Run(":8080")
}
