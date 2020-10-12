package rados

import (
	"fmt"

	"github.com/ceph/go-ceph/rados"
)

//Open ....
func (c *RadosCeph) Open() error {
	if c.conn != nil {
		return nil
	}
	conn, err := rados.NewConnWithUser("admin")
	if err != nil {
		return fmt.Errorf("cannot create new ceph connection: %v", err)
	}

	if err := conn.SetConfigOption(
		"mon_host",
		c.MonHosts,
	); err != nil {
		return fmt.Errorf("error when set ceph monitor host: %v", err)
	}

	if err := conn.SetConfigOption("key", c.Keyring); err != nil {
		return fmt.Errorf("error when set ceph admin key: %v", err)
	}

	if err := conn.SetConfigOption("client_mount_timeout", fmt.Sprintf("%d", c.Timeout)); err != nil {
		return fmt.Errorf("error when set ceph client timeout: %v", err)
	}

	if err := conn.Connect(); err != nil {
		return fmt.Errorf("cannot connect to ceph cluster: %v", err)
	}

	c.conn = conn
	return nil
}
