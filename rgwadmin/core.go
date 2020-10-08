package rgwadmin

import (
	"io"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
	v4 "github.com/aws/aws-sdk-go/aws/signer/v4"
	"github.com/pegasus-cloud/ceph_client/utility"
)

type RGWAdminConfig struct {
	AccessKey string
	SecretKey string
	Host      string
	Region    string
}

func (s *RGWAdminConfig) sendRequestWithAWSV4(method, url string, header map[string]string, body io.Reader) ([]byte, http.Header, int, error) {
	req, _ := http.NewRequest(method, url, body)

	// genSignature
	signer := v4.NewSigner(credentials.NewStaticCredentials(s.AccessKey, s.SecretKey, ""))
	signer.Sign(req, nil, "s3", s.Region, time.Now().UTC())

	signatureHeader := map[string]string{}
	for key, value := range header {
		signatureHeader[key] = value
	}
	for key, value := range req.Header {
		signatureHeader[key] = value[0]
	}

	return utility.SendRequest(method, url, signatureHeader, body)
}
