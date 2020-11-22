package cmd

import "github.com/spf13/cobra"

var RootCmd = &cobra.Command{
	Use:   "task [command]",
	Short: "task is a CLI for managing your TODOs.",
}
