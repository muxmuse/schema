package cmd

import (
  _ "log"
  "fmt"

  "github.com/spf13/cobra"

  _ "github.com/muxmuse/schema/mfa"
  "github.com/muxmuse/schema/schema"
)

func init() {
  rootCmd.AddCommand(useCmd)
}

var useCmd = &cobra.Command{
  Use:   "use <connection>",
  Short: "Select a connection",
  Long:  `Select a configured in ~/.schemapm/config to be used for all commands.
          This sets the property "selected" to true.`,
  Run: func(cmd *cobra.Command, args []string) {
    config := schema.GetConfig()

    for i := range config.Connections {
      use := config.Connections[i].Name == args[0]
      fmt.Println(config.Connections[i].Name, use)
      config.Connections[i].Selected = use
    }

    schema.SaveConfig(config)
  },
}
