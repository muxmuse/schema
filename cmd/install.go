package cmd

import (
  _ "log"

  "github.com/spf13/cobra"

  "github.com/muxmuse/schema/mfa"
  "github.com/muxmuse/schema/schema"
  "gopkg.in/src-d/go-git.v4/plumbing"
)

func init() {
  rootCmd.AddCommand(installCmd)
}

var installCmd = &cobra.Command{
  Use:   "install <repo> [<git tag>]",
  Short: "Install a schema into a database",
  Long:  `Install a schema into a database
          repo: url (or filepath) of git repsitory
          git tag: tag to checkout or current directory state if omited
  `,
  Args: cobra.RangeArgs(1, 2),
  Run: func(cmd *cobra.Command, args []string) {

    if(schema.SelectedConnectionConfig.Log < 2) {
      schema.SelectedConnectionConfig.Log = 2
    }
    schema.Connect()

    if len(args) == 1 {
      err, s := schema.CheckoutDev(args[0])
      mfa.CatchFatal(err)
      schema.Install(s)
    } else {
      err, s := schema.Checkout(args[0], plumbing.NewTagReferenceName(args[1]))
      mfa.CatchFatal(err)
      schema.Install(s)
    }
    
		defer schema.DB.Close();
  },
}
