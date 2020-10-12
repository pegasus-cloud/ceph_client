package s3

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

//UploadPartObject ...
func (s *S3Config) UploadPartObject(object, uploadID string, partnum int) (*http.Response, error) {
	res, err := s.Send(http.MethodPut, fmt.Sprintf("%s/%s?partNumber=%d&uploadId=%s", s.Bucket, object, partnum, uploadID))
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
