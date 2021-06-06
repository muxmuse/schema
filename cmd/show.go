package cmd

import (
  _ "log"

  "github.com/spf13/cobra"

  _ "github.com/muxmuse/schema/mfa"
  "github.com/muxmuse/schema/schema"
)

func init() {
  rootCmd.AddCommand(showCmd)
}

var showCmd = &cobra.Command{
  Use:   "show <schema-name>",
  Short: "Show information about one installed schema",
  Long:  ``,
  Args: cobra.ExactArgs(1),
  Run: func(cmd *cobra.Command, args []string) {
    if(schema.SelectedConnectionConfig.Log < 2) {
      schema.SelectedConnectionConfig.Log = 2
    }
    schema.Connect()
    schema.Show(args[0])
    defer schema.DB.Close()
  },
}
