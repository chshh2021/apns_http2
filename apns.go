package apns_http2

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/http2"
)

const (
	SERVICE_PRODUCTION = "https://api.push.apple.com/3/device/"
	SERVICE_SANDBOX    = "https://api.development.push.apple.com/3/device/"
)

var (
	ErrCertification = errors.New("invalid certificate")
)

type ErrorResponse struct {
	Reason    string `json:"reason"`
	Timestamp int64  `json:"timestamp"`
}

func (this *ErrorResponse) Error() string {
	return this.Reason
}

type Client struct {
	*http.Client
	topic string
}

func New(pem_file string) (*Client, error) {
	cert, err := tls.LoadX509KeyPair(pem_file, pem_file)
	if err != nil {
		return nil, err
	}

	x509cert, err := x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		return nil, err
	}

	topic := strings.Split(x509cert.Subject.CommonName, ": ")[1]

	tlsConfig := &tls.Config{Certificates: []tls.Certificate{cert}}

	if len(tlsConfig.Certificates) > 0 {
		tlsConfig.BuildNameToCertificate()
	}

	transport := &http.Transport{TLSClientConfig: tlsConfig}
	err = http2.ConfigureTransport(transport)
	if err != nil {
		return nil, err
	}

	return &Client{&http.Client{Transport: transport}, topic}, nil
}

func (this *Client) Send(token string, data string, expire int64) error {
	req, err := http.NewRequest(http.MethodPost, SERVICE_PRODUCTION+token, strings.NewReader(data))
	if err != nil {
		return err
	}

	req.Header.Set("apns-topic", this.topic)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("apns-expiration", fmt.Sprintf("%d", time.Now().Unix()+expire))

	res, err := this.Do(req)
	if err != nil {
		if strings.Contains(err.Error(), "tls: internal error") {
			return ErrCertification
		}
		return err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		return nil
	}

	return parseErrorResponse(res.Body, res.StatusCode)
}
