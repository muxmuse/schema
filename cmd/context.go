package cmd

import (
  _ "log"
  "fmt"

  "github.com/spf13/cobra"

  _ "github.com/muxmuse/schema/mfa"
  "github.com/muxmuse/schema/schema"
)

func init() {
  rootCmd.AddCommand(contextCmd)
}

var contextCmd = &cobra.Command{
  Use:   "context use <",
  Short: "Select a connection",
  Long:  `Select a configured in ~/.schemapm/config to be used for all commands.
          This sets the property "selected" to true.`,
  Run: func(cmd *cobra.Command, args []string) {
    config := schema.GetConfig()

    for i := range config.Connections {
      if(len(args) > 0) {
        use := config.Connections[i].Name == args[0]
        config.Connections[i].Selected = use
      }

      if(config.Connections[i].Selected) {
        fmt.Print("> ")
      } else {
        fmt.Print("  ")
      }
      
      fmt.Println(config.Connections[i].Name)
    }

    schema.SaveConfig(config)
  },
}
