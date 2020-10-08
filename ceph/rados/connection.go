package rados

import (
	"fmt"

	"github.com/ceph/go-ceph/rados"

	"go.uber.org/zap"
)

//Open ....
func (c *RadosCeph) Open() error {
	if c.conn != nil {
		return nil
	}
	conn, err := rados.NewConnWithUser("admin")
	if err != nil {
		zap.L().Error(fmt.Sprintf("cannot create new ceph connection: %v", err))
		return err
	}

	if err := conn.SetConfigOption(
		"mon_host",
		c.MonHosts,
	); err != nil {
		zap.L().Error(fmt.Sprintf("error when set ceph monitor host: %v", err))
		return err
	}

	if err := conn.SetConfigOption("key", c.Keyring); err != nil {
		zap.L().Error(fmt.Sprintf("error when set ceph admin key: %v", err))
		return err
	}

	if err := conn.SetConfigOption("client_mount_timeout", fmt.Sprintf("%d", c.Timeout)); err != nil {
		zap.L().Error(fmt.Sprintf("error when set ceph client timeout: %v", err))
		return err
	}

	if err := conn.Connect(); err != nil {
		zap.L().Error(fmt.Sprintf("cannot connect to ceph cluster: %v", err))
		return err
	}

	c.conn = conn

	zap.L().Info("Connecting to ceph controller")
	return nil
}
