package cmd

import (
  _ "log"

  "github.com/spf13/cobra"
  "github.com/muxmuse/schema/schema"
)

func init() {
  rootCmd.AddCommand(createCmd)
}

func defaultTo(a *string, b *string) *string {
  if a == nil {
    return b
  }
  return a
}

var createCmd = &cobra.Command{
  Use:   "create <name> [<dir>] [<version>]",
  Short: "create a new schema",
  Long:  `create a directory containing files that constitute a schema
          name: name of the schema
          dir: target directory. default ./<name>.schema
          version: format v0.0.0. default v0.0.1
  `,
  Args: cobra.RangeArgs(1, 3),
  Run: func(cmd *cobra.Command, args []string) {
    dir := "./" + args[0] + ".schema"
    version := "v0.0.1"
    if len(args) > 1 {
      dir = args[0]
    }
    if len(args) > 2 {
      version = args[2]
    }
    schema.CreateNew(args[0], dir, version)
  },
}
