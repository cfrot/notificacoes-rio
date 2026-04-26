package model

type Notificacao struct {
	ID             int
	ChamadoID      string
	Tipo           string
	CPF            string
	CPFHash        string
	StatusAnterior string
	StatusNovo     string
	Titulo         string
	Descricao      string
	Timestamp      string
	Lida           bool
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
