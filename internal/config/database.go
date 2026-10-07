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

// ConnectDatabase inicializa a conexão com o PostgreSQL e roda as migrações
func ConnectDatabase() *gorm.DB {
	// Tenta carregar o arquivo .env (se estiver rodando localmente no Mac)
	// No Railway, ele vai ignorar isso e usar as variáveis da nuvem
	err := godotenv.Load()
	if err != nil {
		log.Println("Aviso: Arquivo .env não encontrado. Usando variáveis de ambiente do sistema.")
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("ERRO: DATABASE_URL não está configurada!")
	}

	// Conecta no PostgreSQL do Supabase
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Mostra os comandos SQL no terminal
	})
	if err != nil {
		log.Fatal("Falha ao conectar no banco de dados: ", err)
	}

	log.Println("🔥 Conectado ao PostgreSQL com sucesso!")

	// O AutoMigrate lê as nossas structs e cria as tabelas no banco de dados
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
