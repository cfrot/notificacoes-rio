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
	service *service.NotificacaoService
	secret  string
}

func NewNotificacaoHandler(s *service.NotificacaoService, secret string) *NotificacaoHandler {
	return &NotificacaoHandler{
		service: s,
		secret:  secret,
	}
}

type WebhookPayload struct {
	ChamadoID      string `json:"chamado_id"`
	Tipo           string `json:"tipo"`
	CPF            string `json:"cpf"`
	StatusAnterior string `json:"status_anterior"`
	StatusNovo     string `json:"status_novo"`
	Titulo         string `json:"titulo"`
	Descricao      string `json:"descricao"`
	Timestamp      string `json:"timestamp"`
}

func validarAssinatura(body []byte, signature string, secret string) bool {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)

	hashEsperado := "sha256=" + hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(hashEsperado), []byte(signature))
}

func hashCPF(cpf string) string {
	hash := sha256.Sum256([]byte(cpf))
	return hex.EncodeToString(hash[:])
}

func (h *NotificacaoHandler) Webhook(c *gin.Context) {
	signature := c.GetHeader("X-Signature-256")

	body, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "erro ao ler body"})
		return
	}

	if !validarAssinatura(body, signature, h.secret) {
		c.JSON(http.StatusUnauthorized, gin.H{"erro": "assinatura invalida"})
		return
	}

	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	var payload WebhookPayload

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "json invalido"})
		return
	}

	notificacao := model.Notificacao{
		ChamadoID:      payload.ChamadoID,
		CpfHash:        hashCPF(payload.CPF),
		Tipo:           payload.Tipo,
		StatusAnterior: payload.StatusAnterior,
		StatusNovo:     payload.StatusNovo,
		Titulo:         payload.Titulo,
		Descricao:      payload.Descricao,
		Timestamp:      payload.Timestamp,
	}

	err = h.service.ProcessarWebhook(notificacao)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "erro ao salvar"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"mensagem": "evento recebido",
	})
}

func (h *NotificacaoHandler) Listar(c *gin.Context) {
	cpf := c.Query("cpf") // simplificado por enquanto

	notificacoes, err := h.service.Listar(cpf)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "erro ao buscar"})
		return
	}

	c.JSON(http.StatusOK, notificacoes)
}

func (h *NotificacaoHandler) MarcarComoLida(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "id invalido"})
		return
	}

	err = h.service.MarcarComoLida(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "erro ao atualizar"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"mensagem": "notificacao marcada como lida",
	})
}

func (h *NotificacaoHandler) UnreadCount(c *gin.Context) {
	cpf := c.Query("cpf") // depois vamos trocar por JWT

	total, err := h.service.ContarNaoLidas(cpf)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "erro ao contar"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"unread": total,
	})
}
