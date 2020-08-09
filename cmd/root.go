package cmd

import (
  "log"
  "os"

  "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
  Use:   "schema",
  Short: "SchemaPM is a package manager for Microsoft SQL Server",
  Long:  ``,
  Run: func(cmd *cobra.Command, args []string) { },
}

func Execute() {
  if err := rootCmd.Execute(); err != nil {
    log.Println(err)
    os.Exit(1)
  }
}

