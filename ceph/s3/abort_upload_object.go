package s3

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

//AbortUploadObject ...
func (s *S3Config) AbortUploadObject(object, uploadID string) error {
	res, err := s.Send(http.MethodDelete, fmt.Sprintf("%s/%s?uploadId=%s", s.Bucket, object, uploadID))
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusNoContent {
		defer res.Body.Close()
		body := S3Error{}
		bodyBytes, _ := ioutil.ReadAll(res.Body)
		xml.Unmarshal(bodyBytes, &body)

		return errors.New(body.Code)
	}
	return err
}
