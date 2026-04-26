package repository

import (
	"database/sql"
	"notificacoes-rio/internal/model"
)

type NotificacaoRepository struct {
	DB *sql.DB
}

func NewNotificacaoRepository(db *sql.DB) *NotificacaoRepository {
	return &NotificacaoRepository{DB: db}
}

// =====================
// SALVAR
// =====================
func (r *NotificacaoRepository) Salvar(n model.Notificacao) error {
	_, err := r.DB.Exec(`
		INSERT INTO notifications (
			chamado_id,
			cpf_hash,
			tipo,
			status_anterior,
			status_novo,
			titulo,
			descricao,
			timestamp
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
	`,
		n.ChamadoID,
		n.CpfHash,
		n.Tipo,
		n.StatusAnterior,
		n.StatusNovo,
		n.Titulo,
		n.Descricao,
		n.Timestamp,
	)

	return err
}

// =====================
// LISTAR
// =====================
func (r *NotificacaoRepository) ListarPorCPF(cpfHash string) ([]model.Notificacao, error) {
	rows, err := r.DB.Query(`
		SELECT id, chamado_id, tipo, status_anterior, status_novo, titulo, descricao, timestamp, lida, criada_em
		FROM notifications
		WHERE cpf_hash = $1
		ORDER BY id DESC
	`, cpfHash)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notificacoes []model.Notificacao

	for rows.Next() {
		var n model.Notificacao

		err := rows.Scan(
			&n.ID,
			&n.ChamadoID,
			&n.Tipo,
			&n.StatusAnterior,
			&n.StatusNovo,
			&n.Titulo,
			&n.Descricao,
			&n.Timestamp,
			&n.Lida,
			&n.CriadaEm,
		)

		if err != nil {
			return nil, err
		}

		notificacoes = append(notificacoes, n)
	}

	return notificacoes, nil
}

// =====================
// MARCAR COMO LIDA
// =====================
func (r *NotificacaoRepository) MarcarComoLida(id int) error {
	_, err := r.DB.Exec(`
		UPDATE notifications
		SET lida = true
		WHERE id = $1
	`, id)

	return err
}

// =====================
// CONTAR NÃO LIDAS
// =====================
func (r *NotificacaoRepository) ContarNaoLidas(cpfHash string) (int, error) {
	var total int

	err := r.DB.QueryRow(`
		SELECT COUNT(*)
		FROM notifications
		WHERE cpf_hash = $1 AND lida = false
	`, cpfHash).Scan(&total)

	return total, err
}
