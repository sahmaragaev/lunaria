package server

import (
	"log"
	"net/http"
	_ "net/http/pprof"

	jsoniter "github.com/json-iterator/go"

	"github.com/sahmaragaev/lunaria-backend/internal/config"
	"github.com/sahmaragaev/lunaria-backend/internal/database/mongodb"
	"github.com/sahmaragaev/lunaria-backend/internal/database/postgres"
	"github.com/sahmaragaev/lunaria-backend/internal/router"
	"github.com/spf13/cobra"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

var ServerCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the Lunaria backend server",
	Run: func(cmd *cobra.Command, args []string) {
		go func() {
			log.Println(http.ListenAndServe(":6060", nil))
		}()

		cfg, err := config.Load()
		if err != nil {
			log.Fatal("Failed to load config:", err)
		}

		postgresDB, err := postgres.NewPostgresConnection(cfg.Postgres)
		if err != nil {
			log.Fatal("Failed to connect to PostgreSQL:", err)
		}
		defer postgresDB.Close()

		mongoDB, err := mongodb.NewMongoConnection(cfg.MongoDB)
		if err != nil {
			log.Fatal("Failed to connect to MongoDB:", err)
		}
		defer mongoDB.Close()

		router := router.SetupRouter(cfg, postgresDB, mongoDB)
		log.Printf("Starting Lunaria backend on port %s", cfg.Server.Port)
		if err := router.Run(":" + cfg.Server.Port); err != nil {
			log.Fatal("Failed to start server:", err)
		}
	},
}
