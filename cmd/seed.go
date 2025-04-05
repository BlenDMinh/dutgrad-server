package cmd

import (
	"log"

	"github.com/BlenDMinh/dutgrad-server/configs"
	"github.com/BlenDMinh/dutgrad-server/databases"
	"github.com/BlenDMinh/dutgrad-server/seeders"
	"github.com/spf13/cobra"
)

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Seed the database",
	Long:  `Seed the database with initial data`,
	Run: func(cmd *cobra.Command, args []string) {
		configs.Init()
		databases.Init()

		seeder_list := []seeders.Seeder{
			&seeders.SpaceRoleSeeder{},
			&seeders.MockAccountSeeder{},
			&seeders.SpaceSeeder{},
		}

		for _, seeder := range seeder_list {
			log.Printf("Seeding %s...\n", seeder.Name())
			log.Printf("%s: Truncating...\n", seeder.Name())
			if err := seeder.Truncate(); err != nil {
				panic(err)
			}

			log.Printf("%s: Seeding...\n", seeder.Name())
			if err := seeder.Seed(); err != nil {
				panic(err)
			}

			log.Printf("%s: Seeded successfully\n", seeder.Name())
		}

		defer databases.Close()
	},
}
