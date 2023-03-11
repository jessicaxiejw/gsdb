package cmd

import (
	"gsdb/internal/sheet"
	"gsdb/internal/sql"
	"io/ioutil"
	"log"
	"net"

	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Running the server that receives and processes SQL commands",
	Run: func(cmd *cobra.Command, args []string) {
		// fetching flags
		var credentialPath, rootDir, host, port, protocol string
		var err error
		if credentialPath, err = cmd.Flags().GetString("credential-json-path"); err != nil {
			log.Fatal(err) // TODO: wrap error
		}
		if protocol, err = cmd.Flags().GetString("protocol"); err != nil {
			log.Fatal(err) // TODO: wrap error
		}
		if host, err = cmd.Flags().GetString("host"); err != nil {
			log.Fatal(err) // TODO: wrap error
		}
		if port, err = cmd.Flags().GetString("port"); err != nil {
			log.Fatal(err) // TODO: wrap error
		}
		if rootDir, err = cmd.Flags().GetString("root-directory"); err != nil {
			log.Fatal(err) // TODO: wrap error
		}

		// generating clients
		cred, err := ioutil.ReadFile(credentialPath)
		if err != nil {
			log.Fatal(err) // TODO: wrap error
		}
		sheetClient, err := sheet.New(cred, rootDir)
		if err != nil {
			log.Fatal(err) // TODO: wrap error
		}
		client := sql.NewPostgreSQL(sheetClient)

		// listening to incoming requests
		listen, err := net.Listen(protocol, host+":"+port)
		if err != nil {
			log.Fatal(err)
		}
		defer listen.Close()

		// TODO: support multiple connections
		for {
			conn, err := listen.Accept()
			if err != nil {
				log.Fatal(err)
			}
			go func() {
				defer conn.Close()
				buffer := make([]byte, 1024) // TODO: better buffer allocation
				_, err := conn.Read(buffer)
				if err != nil {
					log.Fatal(err) // TODO: wrap error
				}
				err = client.Execute(string(buffer))
				if err != nil {
					log.Fatal(err) // TODO: wrap error
				}
			}()
		}

	},
}

func init() {
	serverCmd.Flags().StringP("credential-json-path", "c", "", "path to the .json file downloaded from the Google Cloud Console under the service account")
	serverCmd.MarkFlagRequired("credential-json-path")
	serverCmd.Flags().StringP("root-directory", "r", "gsdb", "the path of the shared folder in Google Drive with the service account. This will be the path which all your data will be stored")
	serverCmd.Flags().StringP("host", "", "localhost", "the host of which the database is listening to for commands")
	serverCmd.Flags().StringP("port", "p", "9001", "the port of which the database is listening to for commands")
	serverCmd.Flags().StringP("protocol", "", "tcp", "the type of request the database is listening to for commands")

	rootCmd.AddCommand(serverCmd)
}
