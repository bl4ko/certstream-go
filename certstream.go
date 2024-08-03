package certstream

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

const (
	PingPeriod     time.Duration = 15 * time.Second
	DefaultTimeout               = 15
	DefaultSleep                 = 5
)

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

type Source struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func EventStream(skipHeartbeats bool, certstreamServerURL string, timeout int) (chan Message, chan error) {
	if timeout == 0 {
		timeout = DefaultTimeout
	}

	outputStream := make(chan Message)
	errStream := make(chan error)

	go connectAndListen(certstreamServerURL, timeout, skipHeartbeats, outputStream, errStream)

	return outputStream, errStream
}

func connectAndListen(url string, timeout int, skipHeartbeats bool, outputStream chan Message, errStream chan error) {
	for {
		c, _, err := websocket.DefaultDialer.Dial(url, nil)
		if err != nil {
			errStream <- fmt.Errorf("Error connecting to certstream: %w", err)
			time.Sleep(DefaultSleep * time.Second)
			continue
		}
		done := make(chan struct{})
		go sendPingMessages(c, done, errStream)

		if err := readMessages(c, timeout, skipHeartbeats, outputStream); err != nil {
			errStream <- fmt.Errorf("Error reading messages: %w", err)
			close(done)
			c.Close()
			time.Sleep(DefaultSleep * time.Second)
			continue
		}

		close(done)
		c.Close()
	}
}

func sendPingMessages(c *websocket.Conn, done chan struct{}, errStream chan error) {
	ticker := time.NewTicker(PingPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := c.WriteMessage(websocket.PingMessage, nil); err != nil {
				errStream <- fmt.Errorf("Error sending ping message: %w", err)
				return
			}
		case <-done:
			return
		}
	}
}

func readMessages(c *websocket.Conn, timeout int, skipHeartbeats bool, outputStream chan Message) error {
	for {
		if err := c.SetReadDeadline(time.Now().Add(time.Duration(timeout) * time.Second)); err != nil {
			return fmt.Errorf("Error creating wss deadline: %w", err)
		}

		_, rawMessage, err := c.ReadMessage()
		if err != nil {
			return fmt.Errorf("Error reading message: %w", err)
		}

		var message Message
		if err := json.Unmarshal(rawMessage, &message); err != nil {
			return fmt.Errorf("Error unmarshalling certstream message: %w", err)
		}

		if skipHeartbeats && message.MessageType == "heartbeat" {
			continue
		}

		outputStream <- message
	}
}
