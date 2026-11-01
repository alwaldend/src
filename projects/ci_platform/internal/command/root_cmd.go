package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "ci_platform",
	Short: "ci_platform",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("config", "c", "", "config path")
	err := viper.BindPFlag(
		"config", rootCmd.PersistentFlags().Lookup("config"),
	)
	if err != nil {
		panic(err)
	}
}
