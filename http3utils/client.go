package http3utils

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"

	"github.com/lucas-clemente/quic-go/http3"
)

func TrustedClient(expectedPublicKey []byte) (*http.Client, error) {
	verifyPeer := func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
		// Expect exactly the given certificate
		if len(rawCerts) > 1 {
			return fmt.Errorf("expected single certificate")
		}
		// Parse certificate
		certificate, err := x509.ParseCertificate(rawCerts[0])
		if err != nil {
			return fmt.Errorf("failed to parse certificate: %w", err)
		}
		// Marshal back public key into comparable form
		publicKey, err := x509.MarshalPKIXPublicKey(certificate.PublicKey)
		if err != nil {
			return fmt.Errorf("marshal public key back to bytes: %w", err)
		}
		// Make sure that public key matches
		if !bytes.Equal(publicKey, expectedPublicKey) {
			return fmt.Errorf("certificates do not match")
		}
		return nil
	}
	return &http.Client{
		Transport: &http3.RoundTripper{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify:    true,
				VerifyPeerCertificate: verifyPeer,
			},
		},
	}, nil
}

var DefaultClient = &http.Client{
	Transport: &http3.RoundTripper{},
}

var DefaultClientInsecure = &http.Client{
	Transport: &http3.RoundTripper{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	},
}
