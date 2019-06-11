package main

import (
	"fmt"

	"github.com/4ever9/tq"
	"github.com/spf13/cobra"
)

var all bool

func init() {
	versionCMD.Flags().BoolVarP(&all, "all", "a", false, "show all version info")
}

var versionCMD = &cobra.Command{
	Use:   "version [flags]",
	Short: "Show version about app",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Print(getVersion(all))
	},
}

func getVersion(all bool) string {
	version := fmt.Sprintf("TQ version: %s-%s-%s\n", tq.CurrentVersion, tq.CurrentBranch, tq.CurrentCommit)
	if all {
		version += fmt.Sprintf("App build date: %s\n", tq.BuildDate)
		version += fmt.Sprintf("System version: %s\n", tq.Platform)
		version += fmt.Sprintf("Golang version: %s\n", tq.GoVersion)
	}

	return version
}
