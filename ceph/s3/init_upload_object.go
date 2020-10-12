package s3

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

//InitUploadObject ...
func (s *S3Config) InitUploadObject(bucket, object string) (*http.Response, error) {
	res, err := s.Send(http.MethodPost, fmt.Sprintf("%s/%s?uploads", s.Bucket, object))
	if err != nil {
		return res, err
	}
	if res.StatusCode != http.StatusOK {
		defer res.Body.Close()
		body := S3Error{}
		bodyBytes, _ := ioutil.ReadAll(res.Body)
		xml.Unmarshal(bodyBytes, &body)

		return res, errors.New(body.Code)
	}
	return res, err
}
