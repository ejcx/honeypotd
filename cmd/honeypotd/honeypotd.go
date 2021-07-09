package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ejcx/honeypotd/honeypots"
	"github.com/ejcx/honeypotd/honeypots/http"
	sshpot "github.com/ejcx/honeypotd/honeypots/ssh"
	"github.com/spf13/cobra"
)

const (
	VERSION = "0.0.0"
)

var (
	h = http.HTTPPot{}
	s = sshpot.SSHPot{}

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

	httpCmd = &cobra.Command{
		Use:   "http",
		Short: "Run a http honeypot",
		Run: func(cmd *cobra.Command, args []string) {
			honeypot := &honeypots.HoneyPot{
				Address: "",
				Port:    "8080",
			}
			h.Run(honeypot)

		},
	}

	envCmd = &cobra.Command{
		Use:   "env",
		Short: "Run a honeypot by environment variable",
		Run: func(cmd *cobra.Command, args []string) {
			if e, ok := os.LookupEnv("POT"); ok {
				switch e {
				case "HTTP":
					httpCmd.Run(cmd, args)
				case "SSH":
					sshCmd.Run(cmd, args)
				}
			}
			log.Fatal("Run honeypot with POT env set to HTTP or SSH")
		},
	}
	sshCmd = &cobra.Command{
		Use:   "ssh",
		Short: "Run a SSH honeypot",
		Run: func(cmd *cobra.Command, args []string) {
			honeypot := &honeypots.HoneyPot{
				Address: "",
				Port:    "2022",
			}
			s.Run(honeypot)
		},
	}
)

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(httpCmd)
	rootCmd.AddCommand(sshCmd)
	rootCmd.AddCommand(envCmd)
}

func main() {
	rootCmd.Execute()
}
