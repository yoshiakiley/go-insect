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
	if ServerUrl != "" {
		return
	}
	address, _ := net.InterfaceAddrs()
	for _, addr := range address {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ServerUrl = ipnet.IP.String()
			}
		}
	}
}

func initEtcdCli() {
	cfg := clientv3.Config{
		Endpoints:   []string{GlobalEtcdAddress},
		DialTimeout: 5 * time.Second,
		Username:    EtcdUser,
		Password:    EtcdPassword,
	}
	cli, err := clientv3.New(cfg)
	if err != nil {
		fmt.Println(err)
		GlobalEtcdCLI = nil
	}
	GlobalEtcdCLI = cli
}

func putWithLease(serverUuid string) {
	server := fmt.Sprintf("/prom/local/gateway/prod/%s", ServerName)
	key := fmt.Sprintf("%s_%s", server, serverUuid)
	value := ServerUrl
	if ServerPort != 0 {
		value = fmt.Sprintf("%s:%d", ServerUrl, ServerPort)
	}
	if GlobalEtcdCLI == nil {
		fmt.Println("GlobalEtcdCLI is nil")
		return
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
	if ServerName != "" {
		ticker := time.NewTicker(time.Duration(GlobalEtcdTTL-1) * time.Second)
		serverUuid := uuid.NewV4().String()
		defer ticker.Stop()

		initEtcdCli()
		initServerIP()

		for {

			putWithLease(serverUuid)
			<-ticker.C
		}
	}
}
