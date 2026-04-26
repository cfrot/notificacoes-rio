package main

import (
	"notificacoes-rio/internal/database"
	"notificacoes-rio/internal/handler"
	"notificacoes-rio/internal/repository"
	"notificacoes-rio/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	// 🔥 CARREGA .env
	err := godotenv.Load()
	if err != nil {
		panic("erro ao carregar .env")
	}

	db := database.Conectar()
	defer db.Close()

	repo := &repository.NotificacaoRepository{DB: db}
	service := &service.NotificacaoService{Repo: repo}
	handler := &handler.NotificacaoHandler{Service: service}

	router := gin.Default()

	router.POST("/webhook", handler.Webhook)
	router.GET("/notificacoes", handler.Listar)
	router.PATCH("/notificacoes/:id/read", handler.MarcarComoLida)

	router.Run(":8080")
}
