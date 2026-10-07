package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/jvbenetti/lorekeeper-backend.git/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ConnectDatabase init connection with PostgreSQL and run migrations
func ConnectDatabase() *gorm.DB {
	err := godotenv.Load()
	if err != nil {
		log.Println("Aviso: Arquivo .env não encontrado. Usando variáveis de ambiente do sistema.")
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("ERRO: DATABASE_URL não está configurada!")
	}

	// Connect PostgreSQL in Supabase
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Falha ao conectar no banco de dados: ", err)
	}

	log.Println("🔥 Conectado ao PostgreSQL com sucesso!")

	// Create db tables
	err = db.AutoMigrate(
		&models.Location{},
		&models.Entity{},
		&models.Item{},
	)
	if err != nil {
		log.Fatal("Falha ao rodar as migrações: ", err)
	}

	log.Println("✨ Migrações concluídas! Tabelas criadas/atualizadas.")

	return db
}
