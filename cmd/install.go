package cmd

import (
  _ "log"

  "github.com/spf13/cobra"

  _ "github.com/muxmuse/schema/mfa"
  "github.com/muxmuse/schema/schema"
)

func init() {
  rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
  Use:   "install",
  Short: "Install a schema into a database",
  Long:  `Install a schema into a database`,
  Run: func(cmd *cobra.Command, args []string) {
		defer schema.DB.Close();
  },
}
