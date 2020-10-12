package s3

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

//UploadObject ...
func (s *S3Config) UploadObject(object string, fileReader io.ReadSeeker) (*http.Response, error) {
	res, err := s.Send(http.MethodPut, fmt.Sprintf("%s/%s", s.Bucket, object), fileReader)
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
