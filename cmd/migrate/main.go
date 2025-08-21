package migrate

import (
	"log"

	"github.com/sahmaragaev/lunaria-backend/internal/config"
	"github.com/sahmaragaev/lunaria-backend/internal/database/mongodb"
	"github.com/sahmaragaev/lunaria-backend/internal/database/postgres"
	"github.com/spf13/cobra"
)

var MigrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			log.Fatal("Failed to load config:", err)
		}
		postgresDB, err := postgres.NewPostgresConnection(cfg.Postgres)
		if err != nil {
			log.Fatal("Failed to connect to PostgreSQL:", err)
		}
		defer postgresDB.Close()
		if err := postgres.RunMigrations(postgresDB.DB); err != nil {
			log.Fatal("Postgres migrations failed:", err)
		}
		mongoDB, err := mongodb.NewMongoConnection(cfg.MongoDB)
		if err != nil {
			log.Fatal("Failed to connect to MongoDB:", err)
		}
		defer mongoDB.Close()
		if err := mongodb.RunMigrations(mongoDB.Database); err != nil {
			log.Fatal("MongoDB migrations failed:", err)
		}
		log.Println("Migrations completed successfully.")
	},
}
