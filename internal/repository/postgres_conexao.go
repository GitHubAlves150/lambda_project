package repository

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	_ "github.com/lib/pq"
	"github.com/subosito/gotenv"
)

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

func Config_DB() string {

	//dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName)

	config := LerDotEnv()
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", config.Host, config.Port, config.User, config.Password, config.DBName)
	return dsn
}

func LerDotEnv() Config {
	//Carrego o arquivo .env
	err := gotenv.Load()
	if err != nil {
		log.Printf("Erro ao ler as variaveis de ambiente: ", err)
	}

	//os.Getenv lẽ o valor das variaveis de ambiente do .env
	//Converto a PORT para int
	port_str := os.Getenv("DB_PORT")
	port, err := strconv.Atoi(port_str)
	if err != nil {
		log.Fatal("Err ao converter DB_PORT para inteiro")
	}

	//Popula a struct com as variaveis de ambiente
	cfg := Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     port,
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
	}

	return cfg
}

func ConectaDB() (*sql.DB, error) {
	// ─────────────────────────────────────
	// CONEXÃO COM O BANCO "rastreador"
	// ─────────────────────────────────────

	//1. Conecta com o banco "rastreador"
	//dsn := "host=localhost port=5432 user=rastreador password=rastreador123 dbname=rastreador sslmode=disable"

	dsn := Config_DB()

	//sql.Open() NÂO conect com o banco ainda, ele só prepara a conexão
	// "postgres" = nome do driver registrado pelo _ "github.com/lib/pq"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Erro ao conectar:", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("Erro no ping", err)
	}
	//se chegou até aqui é por que deu conexão
	log.Printf("\nConexão bem suceedida\n")
	return db, nil

}
