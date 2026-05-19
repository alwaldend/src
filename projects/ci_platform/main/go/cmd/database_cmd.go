package cmd

import (
	infra "git.alwaldend.com/alwaldend/src/projects/ci_platform/main/go/infrastructure"
	"github.com/spf13/cobra"
)

var databaseCmd = &cobra.Command{
	Use:   "database",
	Short: "database commands",
	Long:  `description`,
}

var databaseMigrate = &cobra.Command{
	Use:   "migrate",
	Short: "migrate the database",
	Long:  "description",
	Run: func(command *cobra.Command, args []string) {
		infra.NewApp().Migrate()
	},
}

func init() {
	rootCmd.AddCommand(databaseCmd)
	databaseCmd.AddCommand(databaseMigrate)
}
