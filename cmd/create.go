package cmd

import (
  _ "log"

  "github.com/spf13/cobra"
  "github.com/muxmuse/schema/schema"
)

func init() {
  rootCmd.AddCommand(createCmd)
}

var createCmd = &cobra.Command{
  Use:   "create <name> [<dir>]",
  Short: "create a new schema",
  Long:  `create a directory containing files that constitute a schema
          name: name of the schema
          dir: target directory. default ./<name>.schema
  `,
  Args: cobra.RangeArgs(1, 2),
  Run: func(cmd *cobra.Command, args []string) {
    if len(args) == 1 {
      schema.CreateNew(args[0], "./" + args[0] + ".schema")
    } else {
      schema.CreateNew(args[0], args[1])
    }
  },
}
