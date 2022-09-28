package cmd

import (
	"crypto/tls"
	"fmt"

	"github.com/lnsp/ftp2p/security"
	"github.com/lnsp/ftp2p/tavern"

	"github.com/spf13/cobra"
)

var announcementTimeout int64

var serveAddr string
var serveCertificate, servePrivateKey string
var serveInsecure bool

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve as a tavern host and track peers",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		var (
			err       error
			tlsConfig *tls.Config
		)
		if serveInsecure {
			tlsConfig, err = security.NewBasicTLS()
		} else {
			tlsConfig, err = security.NewFileTLS(serveCertificate, servePrivateKey)
		}
		if err != nil {
			return fmt.Errorf("setup TLS: %v", err)
		}
		server := tavern.New(announcementTimeout)
		return server.ListenAndServe(serveAddr, tlsConfig)
	},
}

func init() {
	serveCmd.Flags().StringVarP(&serveAddr, "listen", "l", "localhost:6443", "Address to listen on")
	serveCmd.Flags().BoolVarP(&serveInsecure, "insecure", "x", true, "Use generated certificate")
	serveCmd.Flags().StringVarP(&serveCertificate, "cert", "c", "", "TLS certificate")
	serveCmd.Flags().StringVarP(&servePrivateKey, "key", "k", "", "TLS private key")
	serveCmd.Flags().Int64Var(&announcementTimeout, "timeout", 300, "Announcement timeout in seconds")
	rootCmd.AddCommand(serveCmd)
}
