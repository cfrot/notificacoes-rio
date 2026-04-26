package service

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"notificacoes-rio/internal/model"
	"notificacoes-rio/internal/repository"
)

type NotificacaoService struct {
	Repo *repository.NotificacaoRepository
}

func hashCPF(cpf string) string {
	hash := sha256.Sum256([]byte(cpf))
	return hex.EncodeToString(hash[:])
}

func (s *NotificacaoService) ProcessarWebhook(payload model.WebhookPayload) error {

	n := model.Notificacao{
		ChamadoID:      payload.ChamadoID,
		Tipo:           payload.Tipo,
		CPF:            payload.CPF,
		CPFHash:        hashCPF(payload.CPF),
		StatusAnterior: payload.StatusAnterior,
		StatusNovo:     payload.StatusNovo,
		Titulo:         payload.Titulo,
		Descricao:      payload.Descricao,
		Timestamp:      payload.Timestamp,
	}

	err := s.Repo.Salvar(n)

	if err != nil {
		return errors.New("duplicado")
	}

	return nil
}

func (s *NotificacaoService) Listar(cpf string) ([]model.NotificacaoResponse, error) {

	cpfHash := hashCPF(cpf)

	lista, err := s.Repo.ListarPorCPF(cpfHash)
	if err != nil {
		return nil, err
	}

	resposta := make([]model.NotificacaoResponse, 0)

	for _, n := range lista {
		resposta = append(resposta, model.NotificacaoResponse{
			ChamadoID:      n.ChamadoID,
			Tipo:           n.Tipo,
			CPF:            n.CPF,
			StatusAnterior: n.StatusAnterior,
			StatusNovo:     n.StatusNovo,
			Titulo:         n.Titulo,
			Descricao:      n.Descricao,
			Timestamp:      n.Timestamp,
		})
	}

	return resposta, nil
}

func (s *NotificacaoService) MarcarComoLida(id int, cpf string) (bool, error) {

	cpfHash := hashCPF(cpf)

	return s.Repo.MarcarComoLida(id, cpfHash)
}
