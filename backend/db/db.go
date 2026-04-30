//conexion con la base de datos

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	var err error
	DB, err = aql.Open("postgres", connStr)
	if err != il {
		log.Fatal("Error de conexion:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("Error conectando a la base de datos:", err)
	}
}