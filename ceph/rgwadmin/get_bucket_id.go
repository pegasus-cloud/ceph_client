package rgwadmin

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type (
	BucketInfoRGW struct {
		ID         string `json:"id"`
		BucketName string `json:"bucket"`
		Owner      string `json:"owner"`
	}
)

func (s *RGWAdminConfig) GetBucketID(bucket string) (bi *BucketInfoRGW, err error) {
	bodyBytes, _, status, err := s.sendRequestWithAWSV4(
		"GET",
		fmt.Sprintf("http://%s/admin/bucket?bucket=%s", s.Host, bucket),
		nil,
		nil,
	)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		body := struct {
			Code string
		}{}
		json.Unmarshal(bodyBytes, &body)

		return nil, errors.New(body.Code)
	}

	body := &BucketInfoRGW{}
	json.Unmarshal(bodyBytes, &body)

	return body, err
}
