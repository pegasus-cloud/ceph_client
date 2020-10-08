package rgwadmin

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type (
	//UserInfoRGW ....
	UserInfoRGW struct {
		UserID      string `json:"user_id" mapstructure:"user_id"`
		DisplayName string `json:"display_name" mapstructure:"display_name"`
		SubUsers    []struct {
			ID string `json:"id"`
		} `json:"subusers"`
		Keys []S3Key `json:"keys"`
	}
	//S3Key ...
	S3Key struct {
		ID     string `json:"user"`
		Access string `json:"access_key" mapstructure:"access_key"`
		Secret string `json:"secret_key" mapstructure:"secret_key"`
	}
)

//GetRGWUser ...
func (s *RGWAdminConfig) GetRGWUser(userID string) (ui *UserInfoRGW, err error) {
	if cacheD, err := s.getCache(userID); err == nil {
		return cacheD.(*UserInfoRGW), nil
	}
	bodyBytes, _, status, err := s.sendRequestWithAWSV4(
		"GET",
		fmt.Sprintf("http://%s/admin/user?uid=%s", s.Host, userID),
		nil,
		nil,
	)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		errbody := struct {
			Code string
		}{}

		json.Unmarshal(bodyBytes, &errbody)

		return nil, errors.New(errbody.Code)
	}
	bS := &UserInfoRGW{}
	if e := json.Unmarshal(bodyBytes, bS); e != nil {
		return nil, e
	}
	s.putCache(userID, bS)
	return bS, err
}
