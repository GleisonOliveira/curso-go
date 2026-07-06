// Package testdb provides a reusable test database connection for integration tests.
//
// # Usage
//
//	func TestMain(m *testing.M) {
//		db, cleanup := testdb.Connect()
//		testdb.Migrate()
//		testDB = db
//		code := m.Run()
//		cleanup()
//		os.Exit(code)
//	}
//
// Connect loads .env.test, creates the test database if needed, and opens a
// GORM connection.  It is intentionally cheap so multiple test packages can
// call it without noticeable overhead.
//
// Migrate counts .up.sql files and compares against the clean version stored
// in schema_migrations.  Only the first test package where the database
// version is behind actually runs golang-migrate — subsequent packages see
// currentVersion >= expectedVersion and return after a single integer SELECT.
// Packages do not coordinate with each other; the database is the source of
// truth.
package testdb

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	setupOnce sync.Once
	closeOnce sync.Once
	globalDB  *gorm.DB
)

// Connect loads .env.test, creates the test database if it doesn't exist, and
// opens a GORM connection with a silent logger.
// It does NOT run migrations — call Migrate separately if you need the schema.
// Returns the *gorm.DB and a cleanup function that closes the underlying connection.
func Connect() (*gorm.DB, func()) {
	setupOnce.Do(func() {
		// 1. Descobre a raiz do projeto para localizar .env.test e migrations/
		root := projectRoot()

		// 2. Carrega as variáveis de ambiente do arquivo .env.test
		err := godotenv.Load(filepath.Join(root, ".env.test"))
		if err != nil {
			log.Fatalf("Error loading .env.test: %v", err)
		}

		// 3. Conecta ao PostgreSQL e cria o banco de testes se não existir
		ensureDBExists()

		// 4. Abre a conexão com o banco de testes via GORM (logger silencioso)
		globalDB = connectDB()
	})

	// A cleanup só fecha a conexão uma única vez, mesmo se chamada múltiplas vezes
	cleanup := func() {
		closeOnce.Do(func() {
			sqlDB, err := globalDB.DB()
			if err == nil {
				sqlDB.Close()
			}
		})
	}

	return globalDB, cleanup
}

// Migrate applies pending migrations only when the database version is
// behind the number of .up.sql files on disk.  Because the version is
// persisted in the schema_migrations table, only the first test package
// across any number of processes actually runs golang-migrate; subsequent
// packages skip the entire migration setup (no migrate.New, no second
// connection, no expensive version query — just a single integer SELECT).
func Migrate() {
	root := projectRoot()
	if isMigrated(root) {
		return
	}

	runMigrations(root)
}

// isMigrated counts .up.sql files on disk and compares the count against
// the highest clean version stored in schema_migrations.  If the table does
// not exist the query errors out and we return false (not migrated).
func isMigrated(root string) bool {
	expectedVersion := countUpFiles(filepath.Join(root, "migrations"))

	if expectedVersion == 0 {
		return true
	}

	var currentVersion int
	err := globalDB.Raw("SELECT COALESCE(MAX(version), 0) FROM schema_migrations WHERE dirty = false").Scan(&currentVersion).Error
	if err != nil {
		return false
	}

	return currentVersion >= expectedVersion
}

// countUpFiles returns how many *.up.sql files exist in the given directory.
func countUpFiles(dir string) int {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0
	}

	count := 0
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".up.sql") {
			count++
		}
	}
	return count
}

// projectRoot retorna o caminho absoluto da raiz do projeto,
// subindo 4 níveis a partir da localização deste arquivo fonte.
func projectRoot() string {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename) // infrastructure/database/testdb
	dir = filepath.Dir(dir)       // infrastructure/database
	dir = filepath.Dir(dir)       // infrastructure
	return filepath.Dir(dir)      // emailn (raiz do projeto)
}

// ensureDBExists conecta ao banco padrão "postgres" e cria o
// banco de testes definido em GORM_POSTGRES_DB caso não exista.
func ensureDBExists() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=postgres port=%s sslmode=disable TimeZone=%s",
		os.Getenv("GORM_POSTGRES_HOST"),
		os.Getenv("GORM_POSTGRES_USER"),
		os.Getenv("GORM_POSTGRES_PASSWORD"),
		os.Getenv("GORM_POSTGRES_PORT"),
		os.Getenv("GORM_POSTGRES_TIMEZONE"),
	)

	db, err := sql.Open("postgres", dsn)

	if err != nil {
		log.Fatalf("Failed to connect to postgres: %v", err)
	}

	defer db.Close()

	dbName := os.Getenv("GORM_POSTGRES_DB")

	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)", dbName).Scan(&exists)

	if err != nil {
		log.Fatalf("Failed to check database existence: %v", err)
	}

	if !exists {
		_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
		if err != nil {
			log.Fatalf("Failed to create database %s: %v", dbName, err)
		}
		fmt.Printf("Database %s created\n", dbName)
	}
}

// connectDB abre a conexão GORM com o banco de testes usando
// as variáveis de ambiente já carregadas.
func connectDB() *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s",
		os.Getenv("GORM_POSTGRES_HOST"),
		os.Getenv("GORM_POSTGRES_USER"),
		os.Getenv("GORM_POSTGRES_PASSWORD"),
		os.Getenv("GORM_POSTGRES_DB"),
		os.Getenv("GORM_POSTGRES_PORT"),
		os.Getenv("GORM_POSTGRES_TIMEZONE"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		log.Fatalf("Failed to connect to test database: %v", err)
	}

	return db
}

// runMigrations executa as migrations (arquivos .up.sql) via
// golang-migrate, criando as tabelas necessárias para os testes.
func runMigrations(root string) {
	migrationsPath := filepath.Join(root, "migrations")
	// Usa "file:" + caminho absoluto para evitar problemas de parsing
	// de URL com file:/// no Windows
	sourceURL := "file:" + filepath.ToSlash(migrationsPath)

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&TimeZone=%s",
		os.Getenv("GORM_POSTGRES_USER"),
		os.Getenv("GORM_POSTGRES_PASSWORD"),
		os.Getenv("GORM_POSTGRES_HOST"),
		os.Getenv("GORM_POSTGRES_PORT"),
		os.Getenv("GORM_POSTGRES_DB"),
		os.Getenv("GORM_POSTGRES_TIMEZONE"),
	)

	m, err := migrate.New(sourceURL, dsn)
	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Failed to run migrations: %v", err)
	}
}
