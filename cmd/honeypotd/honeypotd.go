package main

import (
	"fmt"

	"github.com/ejcx/honeypotd/honeypots"
	"github.com/spf13/cobra"
)

const (
	VERSION = "0.0.0"
)

var (
	h       = honeypots.HTTPPot{}
	rootCmd = &cobra.Command{
		Use:   "honeypotd",
		Short: "A simple honeypot daemon",
	}

	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(VERSION)
		},
	}

	runCmd = &cobra.Command{
		Use:   "run",
		Short: "Run a honeypot",
		Run: func(cmd *cobra.Command, args []string) {
			honeypot := &honeypots.HoneyPot{
				Address: "",
				Port:    "8080",
			}
			h.Run(honeypot)

		},
	}
)

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(runCmd)
}

func main() {
	rootCmd.Execute()
}
