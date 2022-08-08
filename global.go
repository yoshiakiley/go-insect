package insect

import (
	clientv3 "go.etcd.io/etcd/client/v3"
)

var GlobalEtcdCLI = &clientv3.Client{}

var GlobalEtcdAddress = ""
var GlobalEtcdTTL int64 = 2

var ServerName = ""
var ServerUrl = ""
var ServerPort int
var EtcdUser = ""
var EtcdPassword = ""
