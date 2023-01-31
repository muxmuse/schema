package cmd

import (
  _ "log"

  "github.com/spf13/cobra"

  _ "github.com/muxmuse/schema/mfa"
  "github.com/muxmuse/schema/schema"
)

func init() {
  rootCmd.AddCommand(&cobra.Command{
    Use:   "dump-data-json",
    Short: "Print json based inserts of all data of all tables to stdout",
    Long:  `Not all datatypes are supported. Watch the printed warnings.`,
    Run: func(cmd *cobra.Command, args []string) {
      if(schema.SelectedConnectionConfig.Log < 2) {
        schema.SelectedConnectionConfig.Log = 2
      }
      schema.Connect()
      deleteBeforeInsert := true
      schema.DumpDataJson(deleteBeforeInsert)
      defer schema.DB.Close()
    },
  })
}
