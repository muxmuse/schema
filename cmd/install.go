package cmd

import (
  _ "log"

  "github.com/spf13/cobra"

  _ "github.com/muxmuse/schema/mfa"
  "github.com/muxmuse/schema/schema"
)

func init() {
  rootCmd.AddCommand(installCmd)
}

var installCmd = &cobra.Command{
  Use:   "install getter source",
  Short: "Install a schema into a database",
  Long:  `Install a schema into a database
          getter: one of {git, file}
          source: url or filepath
  `,
  Args: cobra.ExactArgs(2),
  Run: func(cmd *cobra.Command, args []string) {
    var s schema.TSchema
    s.Getter = args[0]
    s.Url = args[1]
    
    schema.Connect()
    schema.Install(&s)

		defer schema.DB.Close();
  },
}
