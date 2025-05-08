package cmd

import (
	"fmt"
	"os"

	cobra "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
  Use:   "file-installer",
  Short: "File installer",
  Long: "File installer",
  Run: func(cmd *cobra.Command, args []string) {

  },
}

var installCmd = &cobra.Command{
  Use:   "install",
  Short: "Install files",
  Long:  "Install files",
  Run: func(cmd *cobra.Command, args []string) {
  },
}

func init() {
  rootCmd.AddCommand(installCmd)
}

func Execute() {
  if err := rootCmd.Execute(); err != nil {
    fmt.Fprintln(os.Stderr, err)
    os.Exit(1)
  }
}
