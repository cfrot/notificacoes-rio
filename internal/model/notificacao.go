package model

type Notificacao struct {
	ID             int    `json:"id"`
	ChamadoID      string `json:"chamado_id"`
	CpfHash        string `json:"-"`
	Tipo           string `json:"tipo"`
	StatusAnterior string `json:"status_anterior"`
	StatusNovo     string `json:"status_novo"`
	Titulo         string `json:"titulo"`
	Descricao      string `json:"descricao"`
	Timestamp      string `json:"timestamp"`
	Lida           bool   `json:"lida"`
	CriadaEm       string `json:"criada_em"`
}
