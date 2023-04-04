// cmd/pages/main.go
// CDF pages server.

package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/cubeflix/cdf/pages"
	"github.com/spf13/cobra"
)

var addr, path string
var certFile, keyFile string

var rootCmd = &cobra.Command{
	Use:   "cdf-pages",
	Short: "the Cubeflix Document Format pages server",
	Long:  `cdf-pages is a HTTP server for hosting static sites written in CDF.`,
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "start the server",
	Long:  `Start serving the cdf-pages server. If cert and key are set, the server runs with TLS.`,
	Run: func(cmd *cobra.Command, args []string) {
		s, err := pages.LoadServer(path)
		if err != nil {
			fmt.Println("cdf-pages:", err)
			os.Exit(1)
			return
		}
		if len(certFile) != 0 && len(keyFile) != 0 {
			// Use TLS.
			err = http.ListenAndServeTLS(addr, certFile, keyFile, s)
		} else {
			// No TLS.
			err = http.ListenAndServe(addr, s)
		}
		fmt.Println("cdf-pages:", err)
		return
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "print the cdf-pages version",
	Long:  `Print the cdf-pages version.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("cdf-pages version 1.0.0")
	},
}

func Execute() {
	// Add arguments.
	serveCmd.PersistentFlags().StringVar(&addr, "addr", ":80", "the address to serve to (defaults to :80)")
	serveCmd.PersistentFlags().StringVar(&path, "path", "", "the path to the cdf-pages project (defaults to current directory)")
	serveCmd.PersistentFlags().StringVar(&certFile, "cert", "", "the certificate file")
	serveCmd.PersistentFlags().StringVar(&keyFile, "key", "", "the key file")

	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(versionCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println("cdf-pages:", err)
		os.Exit(1)
	}
}

func main() {
	Execute()
}
