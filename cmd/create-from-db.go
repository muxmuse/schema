package cmd

import (
  _ "log"

  "github.com/spf13/cobra"
  "github.com/muxmuse/schema/schema"
)

func init() {
  rootCmd.AddCommand(createFromDbCmd)
}

var createFromDbCmd = &cobra.Command{
  Use:   "create-from-db <name> [<dir>] [<version>]",
  Short: "create-from-db a new schema",
  Long:  `create a directory containing files that constitute a schema pulled from db
          name: name of the schema
          dir: target directory. default ./<name>.schema
          version: format v0.0.0. default v0.0.1
  `,
  Args: cobra.RangeArgs(1, 3),
  Run: func(cmd *cobra.Command, args []string) {
    if(schema.SelectedConnectionConfig.Log < 2) {
      schema.SelectedConnectionConfig.Log = 2
    }
    schema.Connect()

    dir := "./" + args[0] + ".schema"
    version := "v0.0.1"
    if len(args) > 0 {
      dir = args[1]
    }
    if len(args) > 1 {
      version = args[2]
    }
    schema.CreateFromDb(args[0], dir, version)
    defer schema.DB.Close();
  },
}
