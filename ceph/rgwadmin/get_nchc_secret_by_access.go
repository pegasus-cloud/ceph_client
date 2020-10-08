package rgwadmin

import (
	"errors"
	"fmt"
	"regexp"
)

type UserInfo struct {
	GroupUUID string
	GroupID   string
	UserUUID  string
	UserID    string
	Members   []string
	Access    string
	Secret    string
	IsPrivate bool
}

//GetNCHCSecretByAccess ...
func (c *RGWAdminConfig) GetNCHCSecretByAccess(rgwUID, access string) (*UserInfo, error) {
	if cacheD, err := c.getCache(rgwUID + access); err == nil {
		return cacheD.(*UserInfo), nil
	}

	uiRGW, err := c.GetRGWUser(rgwUID)
	if err != nil {
		return nil, err
	}

	const (
		uuid    = "[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}"
		uuidlen = 36
	)

	userInfo := &UserInfo{
		GroupID: uiRGW.DisplayName,
	}

	if regex := regexp.MustCompile(fmt.Sprintf("^%s$", uuid)); regex.MatchString(uiRGW.UserID) {
		userInfo.GroupUUID = uiRGW.UserID
		userInfo.IsPrivate = false
	} else if regex := regexp.MustCompile(fmt.Sprintf("^%s-%s$", uuid, uuid)); regex.MatchString(uiRGW.UserID) {
		userInfo.GroupUUID = uiRGW.UserID[uuidlen+1:]
		userInfo.UserUUID = uiRGW.UserID[:uuidlen]
		userInfo.IsPrivate = true
	} else {
		return nil, fmt.Errorf("Can't find Access Key(%s) using id of RGW (%s)at Ceph RGW", access, rgwUID)
	}

	for _, key := range uiRGW.Keys {
		if uiRGW.UserID == key.ID[:len(uiRGW.UserID)] && key.Access == access {
			userInfo.UserID = key.ID[len(uiRGW.UserID)+1:]
			userInfo.Secret = key.Secret
			break
		}
	}
	for _, key := range uiRGW.SubUsers {
		if uiRGW.UserID == key.ID[:len(uiRGW.UserID)] {
			userInfo.Members = append(userInfo.Members, key.ID[len(uiRGW.UserID)+1:])
		}
	}

	if userInfo.Secret == "" {
		return nil, errors.New("User or subuser not found")
	}

	c.putCache(rgwUID+access, userInfo)
	return userInfo, nil
}
