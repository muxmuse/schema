package cmd

import (
  _ "log"

  "github.com/spf13/cobra"

  _ "github.com/muxmuse/schema/mfa"
  "github.com/muxmuse/schema/schema"
)

func init() {
  rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
  Use:   "diff",
  Short: "List modified objects",
  Long:  ``,
  Run: func(cmd *cobra.Command, args []string) {
    if(schema.SelectedConnectionConfig.Log < 2) {
      schema.SelectedConnectionConfig.Log = 2
    }
    schema.Connect()
    schema.Diff()
    defer schema.DB.Close()
  },
}
