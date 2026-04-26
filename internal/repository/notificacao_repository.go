package repository

import (
	"database/sql"
	"notificacoes-rio/internal/model"
)

type NotificacaoRepository struct {
	DB *sql.DB
}

func (r *NotificacaoRepository) Salvar(n model.Notificacao) error {

	_, err := r.DB.Exec(`
		INSERT INTO notifications
		(chamado_id, tipo, cpf, cpf_hash, status_anterior, status_novo, titulo, descricao, timestamp)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
	`,
		n.ChamadoID,
		n.Tipo,
		n.CPF,
		n.CPFHash,
		n.StatusAnterior,
		n.StatusNovo,
		n.Titulo,
		n.Descricao,
		n.Timestamp,
	)

	return err
}

func (r *NotificacaoRepository) ListarPorCPF(cpfHash string) ([]model.Notificacao, error) {

	rows, err := r.DB.Query(`
		SELECT chamado_id, tipo, cpf, status_anterior, status_novo, titulo, descricao, timestamp
		FROM notifications
		WHERE cpf_hash = $1
		ORDER BY id DESC
	`, cpfHash)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lista []model.Notificacao

	for rows.Next() {
		var n model.Notificacao

		err := rows.Scan(
			&n.ChamadoID,
			&n.Tipo,
			&n.CPF,
			&n.StatusAnterior,
			&n.StatusNovo,
			&n.Titulo,
			&n.Descricao,
			&n.Timestamp,
		)

		if err != nil {
			return nil, err
		}

		lista = append(lista, n)
	}

	return lista, nil
}

func (r *NotificacaoRepository) MarcarComoLida(id int, cpfHash string) (bool, error) {

	result, err := r.DB.Exec(`
		UPDATE notifications
		SET lida = true
		WHERE id = $1 AND cpf_hash = $2
	`, id, cpfHash)

	if err != nil {
		return false, err
	}

	rows, _ := result.RowsAffected()
	return rows > 0, nil
}
