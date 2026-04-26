package handler

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"strconv"

	"notificacoes-rio/internal/model"
	"notificacoes-rio/internal/service"

	"github.com/gin-gonic/gin"
)

type NotificacaoHandler struct {
	Service *service.NotificacaoService
}

func (h *NotificacaoHandler) Webhook(c *gin.Context) {

	secret := "meu-secret"
	signature := c.GetHeader("X-Signature-256")

	body, _ := io.ReadAll(c.Request.Body)

	if !validarAssinatura(body, signature, secret) {
		c.JSON(http.StatusUnauthorized, gin.H{"erro": "assinatura invalida"})
		return
	}

	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	var payload model.WebhookPayload

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "json invalido"})
		return
	}

	err := h.Service.ProcessarWebhook(payload)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"mensagem": "evento duplicado ignorado"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"mensagem": "webhook processado"})
}

func (h *NotificacaoHandler) Listar(c *gin.Context) {

	cpf := c.Query("cpf")

	lista, err := h.Service.Listar(cpf)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": err.Error()})
		return
	}

	c.JSON(http.StatusOK, lista)
}

func (h *NotificacaoHandler) MarcarComoLida(c *gin.Context) {

	id, _ := strconv.Atoi(c.Param("id"))
	cpf := c.Query("cpf")

	ok, err := h.Service.MarcarComoLida(id, cpf)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": err.Error()})
		return
	}

	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"erro": "não encontrada"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"mensagem": "marcada como lida"})
}

func validarAssinatura(body []byte, signature string, secret string) bool {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)

	hash := "sha256=" + hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(hash), []byte(signature))
}
