package rgwadmin

import (
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
	v4 "github.com/aws/aws-sdk-go/aws/signer/v4"
	"github.com/bluele/gcache"
	"github.com/pegasus-cloud/ceph_client/ceph/utility"
)

const RGWAdminConfigCacheTag = "RGWAdminConfigCacheTag"

type RGWAdminConfig struct {
	AccessKey   string
	SecretKey   string
	Host        string
	Region      string
	CacheSize   int
	CacheExpire time.Duration
	_           struct{}
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

func (r *RGWAdminConfig) getCache(key string) (interface{}, error) {
	if c := r.checkCache(); c == nil {
		return nil, errors.New("Cache Not Enable")
	} else {
		return c.Get(key)
	}

}

func (r *RGWAdminConfig) putCache(key string, d interface{}) error {
	if c := r.checkCache(); c == nil {
		return errors.New("Cache Not Enable")
	} else {
		return c.Set(key, d)
	}
}

func (r *RGWAdminConfig) checkCache() gcache.Cache {
	if utility.UseCache(RGWAdminConfigCacheTag) == nil {
		utility.NewWithExpire(RGWAdminConfigCacheTag, r.CacheSize, r.CacheExpire).BuildWithExpire()

	}

	return utility.UseCache(RGWAdminConfigCacheTag)
}
