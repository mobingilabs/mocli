package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var stackCmd = &cobra.Command{
	Use:   "stack",
	Short: "manage your stack",
	Long:  `Manage your infrastructure/application stack.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Check `stack --help` for more information on supported subcommands.")
	},
}

func init() {
	rootCmd.AddCommand(stackCmd)
}
