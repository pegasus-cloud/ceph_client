package s3

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

//UploadPartObject ...
func (s *S3Config) UploadPartObject(object, uploadID string, partnum int, body io.ReadSeeker) (*http.Response, error) {
	res, err := s.Send(http.MethodPut, fmt.Sprintf("%s/%s?partNumber=%d&uploadId=%s", s.Bucket, object, partnum, uploadID), body)
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
