package cmd

import (
  _ "log"

  "github.com/spf13/cobra"

  // "github.com/muxmuse/schema/mfa"
  "github.com/muxmuse/schema/schema"
)

func init() {
  rootCmd.AddCommand(uninstallCmd)
}

var uninstallCmd = &cobra.Command{
  Use:   "uninstall <name>",
  Short: "Remove schema from database",
  Long:  `Remove schema from database
          name: Name of the schema that shall be removed
  `,
  Args: cobra.ExactArgs(1),
  Run: func(cmd *cobra.Command, args []string) {

    if(schema.SelectedConnectionConfig.Log < 2) {
      schema.SelectedConnectionConfig.Log = 2
    }
    schema.Connect()
    schema.Uninstall(args[0])
    
		defer schema.DB.Close();
    schema.CleanUp()
  },
}
