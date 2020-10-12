package s3

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
	v4 "github.com/aws/aws-sdk-go/aws/signer/v4"
)

//S3Config ...
type S3Config struct {
	Host string
	// Body           io.ReadSeeker
	Header         map[string]string
	AdminAccessKey string
	AdminSecretKey string
	Region         string
	Bucket         string
	Timeout        time.Duration
}

//S3Error ...
type S3Error struct {
	Code string
}

//Send ...
func (s *S3Config) Send(method string, path string, bodies ...io.ReadSeeker) (*http.Response, error) {
	var req *http.Request
	var err error
	var body io.ReadSeeker
	if len(bodies) != 0 {
		body = bodies[0]
	}
	req, err = http.NewRequest(method, fmt.Sprintf("%s/%s", s.Host, path), body)
	if err != nil {
		return nil, err
	}

	signer := v4.NewSigner(credentials.NewStaticCredentials(s.AdminAccessKey, s.AdminSecretKey, ""))
	signer.Sign(req, body, "s3", s.Region, time.Now().UTC())

	for key, val := range s.Header {
		req.Header.Set(key, val)
	}

	client := &http.Client{
		Timeout: time.Duration(10) * time.Second,
	}
	return client.Do(req)
}
