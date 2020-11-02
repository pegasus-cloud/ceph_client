package ceph

import (
	"github.com/pegasus-cloud/ceph_client/ceph/rados"
	"github.com/pegasus-cloud/ceph_client/ceph/rgwadmin"
)

var rgwAdminCfg *rgwadmin.RGWAdminConfig

//InitialRGWAdminGlobalConfig ...
func InitialRGWAdminGlobalConfig(c *rgwadmin.RGWAdminConfig) {
	rgwAdminCfg = c
}

// GetBucketID ...
func GetBucketID(bucket string) (bi *rgwadmin.BucketInfoRGW, err error) {
	return rgwAdminCfg.GetBucketID(bucket)
}

//GetRGWUser ...
func GetRGWUser(userID string) (ui *rgwadmin.UserInfoRGW, err error) {
	return rgwAdminCfg.GetRGWUser(userID)
}

//GetNCHCSecretByAccess ...
func GetNCHCSecretByAccess(rgwUID, access string) (*rgwadmin.UserInfo, error) {
	return rgwAdminCfg.GetNCHCSecretByAccess(rgwUID, access)
}

var radosCeph *rados.RadosCeph

//InitialRadosGlobalConfig ...
func InitialRadosGlobalConfig(c *rados.RadosCeph) {
	radosCeph = c
}

//Rados ...
func Rados() *rados.RadosCeph {
	return radosCeph

}

// Mixed functions

//GetHealth ...
func GetHealth() (string, error) {
	return Rados().GetHealth()
}

// GetNCHCSecret ...
func GetNCHCSecret(access string) (*rgwadmin.UserInfo, error) {
	rgwUID, err := Rados().GetRGWUidByAccess(access)
	if err != nil {
		return nil, err
	}

	uInfo, err := GetNCHCSecretByAccess(rgwUID, access)
	if err != nil {
		return nil, err
	}
	return uInfo, nil
}
