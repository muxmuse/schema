package cmd

import (
  _ "log"

  "github.com/spf13/cobra"

  _ "github.com/muxmuse/schema/mfa"
  "github.com/muxmuse/schema/schema"
)

func init() {
  rootCmd.AddCommand(uninstallCmd)
}

var uninstallCmd = &cobra.Command{
  Use:   "uninstall getter source",
  Short: "Remove a schema from database",
  Long:  `Remove a schema from a database
          getter: one of {git, file}
          source: url or filepath
  `,
  Args: cobra.ExactArgs(2),
  Run: func(cmd *cobra.Command, args []string) {
    var s schema.TSchema
    s.Getter = args[0]
    s.Url = args[1]
    
    schema.Uninstall(&s)

		defer schema.DB.Close();
  },
}
