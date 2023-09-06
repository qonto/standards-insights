package tls

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
)

func CreateConfig(keyPath string, certPath string, cacertPath string, insecure bool) (*tls.Config, error) {
	tlsConfig := &tls.Config{} //nolint
	if keyPath != "" {
		cert, err := tls.LoadX509KeyPair(certPath, keyPath)
		if err != nil {
			return nil, fmt.Errorf("Fail to load certificates: %w", err)
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
	}
	if cacertPath != "" {
		caCert, err := os.ReadFile(cacertPath) //nolint
		if err != nil {
			return nil, fmt.Errorf("Fail to load the ca certificate: %w", err)
		}
		caCertPool := x509.NewCertPool()
		result := caCertPool.AppendCertsFromPEM(caCert)
		if !result {
			return nil, fmt.Errorf("fail to read ca certificate on %s", certPath)
		}
		tlsConfig.ClientCAs = caCertPool
		tlsConfig.RootCAs = caCertPool
	}
	tlsConfig.InsecureSkipVerify = insecure
	tlsConfig.MinVersion = tls.VersionTLS12
	return tlsConfig, nil
}

func GetClientAuthType(authType string) (tls.ClientAuthType, error) {
	switch authType {
	case "NoClientCert":
		return tls.NoClientCert, nil
	case "RequestClientCert":
		return tls.RequestClientCert, nil
	case "RequireAnyClientCert":
		return tls.RequireAnyClientCert, nil
	case "VerifyClientCertIfGiven":
		return tls.VerifyClientCertIfGiven, nil
	case "RequireAndVerifyClientCert":
		return tls.RequireAndVerifyClientCert, nil
	default:
		return tls.NoClientCert, fmt.Errorf("Unknown client auth type %s", authType)
	}
}
