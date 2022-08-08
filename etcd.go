package insect

import (
	"context"
	"fmt"
	"net"
	"time"

	uuid "github.com/satori/go.uuid"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func initServerIP() {
	if INSECT_SERVER_URL != "" {
		return
	}
	address, _ := net.InterfaceAddrs()
	for _, addr := range address {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				INSECT_SERVER_URL = ipnet.IP.String()
			}
		}
	}
}

func initEtcdCli() {
	cfg := clientv3.Config{
		Endpoints:   []string{GlobalEtcdAddress},
		DialTimeout: 1 * time.Second,
	}
	cli, err := clientv3.New(cfg)
	if err != nil {
		GlobalEtcdCLI = nil
	}
	GlobalEtcdCLI = cli
}

func putWithLease() {
	server := fmt.Sprintf("/prom/local/gateway/prod/%s", INSECT_SERVER_NAME)
	key := fmt.Sprintf("%s_%s", server, uuid.NewV4().String())
	value := INSECT_SERVER_URL
	if INSECT_SERVER_PORT != 0 {
		value = fmt.Sprintf("%s:%d", INSECT_SERVER_URL, INSECT_SERVER_PORT)
	}

	lease := clientv3.NewLease(GlobalEtcdCLI)
	leaseResp, err := lease.Grant(context.TODO(), GlobalEtcdTTL)
	if err != nil {
		fmt.Println(err)
		return
	}
	leaseID := leaseResp.ID
	_, err = GlobalEtcdCLI.KV.Put(context.TODO(), key, value, clientv3.WithLease(leaseID))
	if err != nil {
		fmt.Println(err)
		return
	}
}

func EtcdProxy() {
	if INSECT_SERVER_NAME != "" {
		initEtcdCli()
		initServerIP()
		for {
			putWithLease()
			time.Sleep(time.Duration(GlobalEtcdTTL) * time.Second)
		}
	}
}
