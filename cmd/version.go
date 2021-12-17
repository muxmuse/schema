package cmd

import (
  "log"
  "fmt"

  "github.com/spf13/cobra"

  _ "github.com/muxmuse/schema/mfa"
  "github.com/muxmuse/schema/schema"
)

func init() {
  rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
  Use:   "version ",
  Short: "Increment version and create migration scripts",
  Long:  `Increment version and create migration scripts`,
  Run: func(cmd *cobra.Command, args []string) {
    config := schema.GetConfig()
    validContextName := len(args) == 0

    for i := range config.Connections {
      validContextName = validContextName || config.Connections[i].Name == args[0]
    }
    
    if(!validContextName) {
      log.Println("No such context ", args[0])
    }

    for i := range config.Connections {
      if(len(args) > 0 && validContextName) {
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
