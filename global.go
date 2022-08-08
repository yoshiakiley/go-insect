package insect

import (
	clientv3 "go.etcd.io/etcd/client/v3"
)

var GlobalEtcdCLI = &clientv3.Client{}

var GlobalEtcdAddress = ""
var GlobalEtcdTTL int64 = 1

var INSECT_SERVER_NAME = ""
var INSECT_SERVER_URL = ""
var INSECT_SERVER_PORT int = 0
