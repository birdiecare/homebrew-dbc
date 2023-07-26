/*
Copyright © 2023 Harry M harry.morgan@birdie.care
*/
package cmd

import (
	"flag"

	"github.com/birdiecare/dbc/handler"
	"github.com/spf13/cobra"
)

var proxyPort string

// proxyCommand represents the proxy command
var proxyCommand = &cobra.Command{
	Use:   "proxy",
	Short: "Open a SOCKS5 proxy to our private infrastructure",
	Long: `Open a SOCKS5 proxy to our private infrastructure.
Opens a SOCKS5 proxy on local port 1080 that allows you to connect to our private infrastructure.

	dbc proxy

Usage:
Open a proxy:

	dbc proxy

Then connect to the DB using a Postgres client that supports SOCKS5 proxies by setting the ALL_PROXY environment variable to export ALL_PROXY=socks5://127.0.0.1:1080

Use a custom local port:

	db proxy -p 1081

`,

	Run: func(cmd *cobra.Command, args []string) {
		flag.Parse()
		// Assert AWS Creds
		handler.AssertCredentials()
		handler.Proxy(proxyPort)
	},
}

func init() {
	rootCmd.AddCommand(proxyCommand)

	//Flags
	proxyCommand.Flags().StringVarP(&proxyPort, "port", "p", "1080", "Local port to open the SOCKS5 proxy on (default 1080)")
}
