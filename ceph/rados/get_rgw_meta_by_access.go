package rados

import (
	"fmt"

	"github.com/pegasus-cloud/ceph_client/ceph/utility"
)

//GetRGWUidByAccess ....
func (c *RadosCeph) GetRGWUidByAccess(access string) (string, error) {

	if cacheData, err := c.getCache(access); err == nil {
		return cacheData.(string), nil
	}

	ioctx, err := c.conn.OpenIOContext(fmt.Sprintf("%s.rgw.meta", c.Region))
	if err != nil {
		return "", err
	}

	ioctx.SetNamespace("users.keys")
	stat, err := ioctx.Stat(access)
	if err != nil {
		return "", err
	}
	uid := make([]byte, stat.Size-4)

	if _, err := ioctx.Read(access, uid, 4); err != nil {
		return "", err
	}

	c.putCache(access, utility.BytesToStr(uid))

	return utility.BytesToStr(uid), nil
}
