package service

import (
	"crypto/sha256"
	"encoding/hex"
	"notificacoes-rio/internal/model"
	"notificacoes-rio/internal/repository"
)

type NotificacaoService struct {
	repo *repository.NotificacaoRepository
}

func NewNotificacaoService(r *repository.NotificacaoRepository) *NotificacaoService {
	return &NotificacaoService{repo: r}
}

func hashCPF(cpf string) string {
	hash := sha256.Sum256([]byte(cpf))
	return hex.EncodeToString(hash[:])
}

func (s *NotificacaoService) ProcessarWebhook(n model.Notificacao) error {
	return s.repo.Salvar(n)
}

func (s *NotificacaoService) Listar(cpf string) ([]model.Notificacao, error) {
	hash := hashCPF(cpf)
	return s.repo.ListarPorCPF(hash)
}

func (s *NotificacaoService) MarcarComoLida(id int) error {
	return s.repo.MarcarComoLida(id)
}

func (s *NotificacaoService) ContarNaoLidas(cpf string) (int, error) {
	hash := hashCPF(cpf)
	return s.repo.ContarNaoLidas(hash)
}
