package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

func main() {
	down := flag.Bool("down", false, "rollback de 1 step")
	steps := flag.Int("steps", 0, "número de steps (positivo=sobe, negativo=desce)")
	flag.Parse()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Erro ao carregar .env")
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&TimeZone=%s",
		os.Getenv("GORM_POSTGRES_USER"),
		os.Getenv("GORM_POSTGRES_PASSWORD"),
		os.Getenv("GORM_POSTGRES_HOST"),
		os.Getenv("GORM_POSTGRES_PORT"),
		os.Getenv("GORM_POSTGRES_DB"),
		os.Getenv("GORM_POSTGRES_TIMEZONE"),
	)

	m, err := migrate.New("file://migrations", dsn)
	if err != nil {
		log.Fatalf("Erro ao criar migrate: %v", err)
	}

	switch {
	case *down:
		if err := m.Steps(-1); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Erro ao executar rollback: %v", err)
		}
	case *steps != 0:
		if err := m.Steps(*steps); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Erro ao executar steps: %v", err)
		}
	default:
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Erro ao rodar migrations: %v", err)
		}
	}

	fmt.Println("Migrations executadas com sucesso!")
}
