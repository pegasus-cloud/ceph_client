package rados

import (
	"github.com/ceph/go-ceph/rados"
)

type RadosCeph struct {
	conn     *rados.Conn
	MonHosts string
	Keyring  string
	Timeout  int
	Region   string
	_        struct{}
}
