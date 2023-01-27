package cmd

import (
  _ "log"

  "github.com/spf13/cobra"

  _ "github.com/muxmuse/schema/mfa"
  "github.com/muxmuse/schema/schema"
)

func init() {
  rootCmd.AddCommand(execCmd)
}

var execCmd = &cobra.Command{
  Use:   "exec",
  Short: "exec batches from stdin (use this to restore dumped data)",
  Long:  ``,
  Run: func(cmd *cobra.Command, args []string) {
    if(schema.SelectedConnectionConfig.Log < 2) {
      schema.SelectedConnectionConfig.Log = 2
    }
    schema.Connect()
    schema.ExecFromStdin()
    defer schema.DB.Close()
  },
}
