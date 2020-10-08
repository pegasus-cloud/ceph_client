package rados

import (
	"errors"
	"time"

	"github.com/bluele/gcache"
	"github.com/ceph/go-ceph/rados"
	"github.com/pegasus-cloud/ceph_client/ceph/utility"
)

const RadosCephCacheTag = "RadosCephCacheTag"

type RadosCeph struct {
	conn        *rados.Conn
	MonHosts    string
	Keyring     string
	Timeout     int
	Region      string
	CacheSize   int
	CacheExpire time.Duration
	_           struct{}
}

func (r *RadosCeph) getCache(key string) (interface{}, error) {
	if c := r.checkCache(); c == nil {
		return nil, errors.New("Cache Not Enable")
	} else {
		return c.Get(key)
	}
}

func (r *RadosCeph) putCache(key string, d interface{}) error {
	if c := r.checkCache(); c == nil {
		return errors.New("Cache Not Enable")
	} else {
		return c.Set(key, d)
	}
}

func (r *RadosCeph) checkCache() gcache.Cache {
	if utility.UseCache(RadosCephCacheTag) == nil {
		utility.NewWithExpire(RadosCephCacheTag, r.CacheSize, r.CacheExpire).BuildWithExpire()

	}

	return utility.UseCache(RadosCephCacheTag)
}
