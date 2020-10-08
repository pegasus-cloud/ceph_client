package rados

import (
	"fmt"

	"github.com/pegasus-cloud/ceph_client/utility"
)

//GetRGWUidByAccess ....
func (c *RadosCeph) GetRGWUidByAccess(access string) (string, error) {
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

	return utility.BytesToStr(uid), nil
}
