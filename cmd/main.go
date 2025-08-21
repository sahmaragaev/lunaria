package main

import (
	"log"
	"os"

	"github.com/spf13/cobra"

	migrate "github.com/sahmaragaev/lunaria-backend/cmd/migrate"
	server "github.com/sahmaragaev/lunaria-backend/cmd/server"
)

var rootCmd = &cobra.Command{
	Use:   "lunaria-backend",
	Short: "Lunaria Backend CLI",
}

func main() {
	rootCmd.AddCommand(server.ServerCmd)
	rootCmd.AddCommand(migrate.MigrateCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
