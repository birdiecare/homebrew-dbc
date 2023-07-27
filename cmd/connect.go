/*
Copyright © 2023 Harry M harry.morgan@birdie.care
*/
package cmd

import (
	"flag"
	"fmt"

	"github.com/birdiecare/dbc/handler"
	"github.com/common-nighthawk/go-figure"
	"github.com/spf13/cobra"
)

var host string
var port string
var region string
var user string
var localport string
var iam bool

// connectCmd represents the connect command
var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Connect to a DB",
	Long: `Connect to a DB.
Opens a connection with the given database @ the given port at localhost:5432

	dbc connect

Password Authentication:
Use an existing database user password to authenticate against the dtaabase once the connection is open.

	dbc connect -h $host -u $user 

Then connect to the DB using the password for $user

	psql -h localhost -p 5432 -U $user -d $database --password

Then paste your password for $user.

IAM Authentication: 
Use IAM Authentication to authenticate as a User to the DB

	dbc connect -h $host -u $user --iam

This command will output a Token to use as a password when connecting to your database.	

	psql -h localhost -p ${localport} -U ${user} -d ${database} --password

Then paste the token`,

	Run: func(cmd *cobra.Command, args []string) {
		flag.Parse()

		// Assert AWS Creds
		handler.AssertCredentials()

		if host == "" {
			host = handler.FuzzEndpoints(iam)
		}

		if iam {
			fmt.Println("")
			myFigure := figure.NewFigure("DBC Connect IAM", "", true)
			myFigure.Print()

			fmt.Println("")
			handler.GenerateToken(host, port, region, user)
		} else {
			myFigure := figure.NewFigure("DBC Connect", "", true)
			myFigure.Print()

			fmt.Println("")
		}
		// Start Port-Forwarding Session
		handler.Handler(region, host, port, localport)
	},
}

func init() {
	rootCmd.AddCommand(connectCmd)

	//Flags
	connectCmd.Flags().StringVarP(&host, "host", "H", "", "Hostname of the Database to open a connection to. If a hostname is not provided, a fuzzyfind list with be presented to select a database, and subsequently, a user to connect to`")
	connectCmd.Flags().StringVarP(&port, "port", "p", "5432", "Port of the Datbase to open a connection to (default 5432)")
	connectCmd.Flags().StringVarP(&region, "region", "r", "eu-west-2", "Region of the Datbase to open a connection to (default eu-west-2)")
	connectCmd.Flags().StringVarP(&user, "user", "u", "", "The DB User to open a connection with")
	connectCmd.Flags().StringVarP(&localport, "localport", "l", "5432", "Local Port to forward database connection to (default 5432)")
	connectCmd.Flags().BoolVarP(&iam, "iam", "I", false, "Bool: Use IAM Authentication for Database Connection - Generates a password token using IAM Authentication")
}
