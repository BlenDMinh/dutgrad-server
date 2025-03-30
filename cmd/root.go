package cmd

import (
	"os"

	"github.com/BlenDMinh/dutgrad-server/configs"
	"github.com/BlenDMinh/dutgrad-server/databases"
	"github.com/BlenDMinh/dutgrad-server/server"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dutgrad-server",
	Short: "Start the server",
	Long:  `Start the server to serve the application`,
	Run: func(cmd *cobra.Command, args []string) {
		configs.Init()
		databases.Init()
		server.Init()

		defer databases.Close()
		defer server.Close()
	},
}

func Execute() {
	rootCmd.AddCommand(seedCmd)
	rootCmd.AddCommand(migrateCmd)
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
