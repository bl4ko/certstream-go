# Certstream-go

Simple go client library for interacting with the cerstream logs, inspired by [CaliDog](https://github.com/CaliDog/certstream-go).

# Usage

```go
package main

import (
	"fmt"

	"github.com/bl4ko/certstream-go"
)

func main() {
	certstreamServerURL := "wss://certstream.calidog.io"
	timeout := 15
	stream, errStream := certstream.EventStream(true, certstreamServerURL, timeout)
	for {
		select {
		case message := <-stream:
			fmt.Printf("Received stream: %+v\n\n", message)
		case err := <-errStream:
			fmt.Printf("Received error: %s\n\n", err)
		}
	}
}
```

# Example certstream.Message

Certstream-go returns data in the following format:

```go
type Message struct {
	MessageType string `json:"message_type"`
	Data        struct {
		CertIndex int    `json:"cert_index"`
		CertLink  string `json:"cert_link"`
		LeafCert  struct {
			AllDomains []string `json:"all_domains"`
			Extensions struct {
				AuthorityInfoAccess    string `json:"authorityInfoAccess"`
				AuthorityKeyIdentifier string `json:"authorityKeyIdentifier"`
				BasicConstraints       string `json:"basicConstraints"`
				CertificatePolicies    string `json:"certificatePolicies"`
				CtlPoisonByte          bool   `json:"ctlPoisonByte"`
				ExtendedKeyUsage       string `json:"extendedKeyUsage"`
				KeyUsage               string `json:"keyUsage"`
				SubjectAltName         string `json:"subjectAltName"`
				SubjectKeyIdentifier   string `json:"subjectKeyIdentifier"`
			} `json:"extensions"`
			Fingerprint        string `json:"fingerprint"`
			NotAfter           int    `json:"not_after"`
			NotBefore          int    `json:"not_before"`
			SerialNumber       string `json:"serial_number"`
			SignatureAlgorithm string `json:"signature_algorithm"`
			Subject            Name   `json:"subject"`
			Issuer             Name   `json:"issuer"`
			IsCA               bool   `json:"is_ca"`
		} `json:"leaf_cert"`
		Seen       float64 `json:"seen"`
		Source     Source  `json:"source"`
		UpdateType string  `json:"update_type"`
	} `json:"data"`
}

type Name struct {
	C            string `json:"C"`
	CN           string `json:"CN"`
	L            string `json:"L"`
	O            string `json:"O"`
	OU           string `json:"OU"`
	ST           string `json:"ST"`
	Aggregated   string `json:"aggregated"`
	EmailAddress string `json:"email_address"`
}
```
