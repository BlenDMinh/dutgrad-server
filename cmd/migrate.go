package cmd

import (
	"github.com/BlenDMinh/dutgrad-server/configs"
	"github.com/BlenDMinh/dutgrad-server/databases"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate the database",
	Long:  `Migrate the database to the latest version`,
	Run: func(cmd *cobra.Command, args []string) {
		configs.Init()
		databases.Init()

		db, err := databases.GetDB().DB()

		if err != nil {
			panic(err)
		}

		goose.SetDialect("postgres")

		migrationDir := "databases/migrations"

		if err := goose.Up(db, migrationDir); err != nil {
			panic(err)
		}

		defer databases.Close()
	},
}
