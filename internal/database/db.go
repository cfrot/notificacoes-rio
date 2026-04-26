package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func Conectar() *sql.DB {

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	// 🔍 DEBUG (pode apagar depois)
	fmt.Println("DB CONFIG:")
	fmt.Println("HOST:", host)
	fmt.Println("PORT:", port)
	fmt.Println("USER:", user)
	fmt.Println("DB:", dbname)

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	// 🔥 VALIDA CONEXÃO DE VERDADE
	if err := db.Ping(); err != nil {
		panic("erro ao conectar no banco: " + err.Error())
	}

	fmt.Println("✅ Conectado ao banco com sucesso")

	return db
}
