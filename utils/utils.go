package utils

import (
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/exec"
	"sync"
	"time"
)

const (
	KubeletSocket   = "kubelet.sock"
	DevicePluginDir = "/var/lib/kubelet/device-plugins/"
	DpManagerPort   = ":60000"
)

func ThreadBlock() {
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}

func Check(err error, msg string) {
	if err != nil {
		log.Warnf("ErrorOccurs: %v, %v", msg, err)
	}
}

func Dial(unixSocketPath string, timeout time.Duration) (*grpc.ClientConn, error) {
	c, err := grpc.Dial(unixSocketPath, grpc.WithInsecure(), grpc.WithBlock(),
		grpc.WithTimeout(timeout),
		grpc.WithDialer(func(addr string, timeout time.Duration) (net.Conn, error) {
			return net.DialTimeout("unix", addr, timeout)
		}),
	)
	if err != nil {
		return nil, err
	}
	return c, nil
}

//func StartRPC(object interface{}, protocol string, address string) {
//	err := rpc.Register(object)
//	Check(err, "Starting Register Failed")
//	lis, err := net.Listen(protocol, address)
//	Check(err, "rpc listening failed")
//	log.Infof("RPC server starting success at %s", address)
//	for {
//		conn, err := lis.Accept()
//		Check(err, "Error occurs when accepting connections")
//		log.Infof("Received connection from %s", conn.RemoteAddr())
//		//go processFunc.(func(net.Conn))(conn)
//		go func(conn net.Conn) {
//			log.Debugf("Processing info from %s...%s", conn.LocalAddr(), conn.RemoteAddr())
//			jsonrpc.ServeConn(conn)
//		}(conn)
//	}
//}

func StartJsonSvr(protocol string, address string, processFunc interface{}) {
	log.Debugf("Trying to start jsonsvr")
	lis, err := net.Listen(protocol, address)
	Check(err, "rpc listening failed")
	for {
		conn, err := lis.Accept()
		Check(err, "Error occurs when accepting connections")
		log.Infof("Received connection from %s", conn.RemoteAddr())
		go processFunc.(func(*net.Conn))(&conn)
	}
}

func ReadFromCmd(command string) string {
	output, err := exec.Command("/bin/bash", "-c", command).Output()
	Check(err, "Execution failed")
	return string(output)
}

func GetCWD() string {
	dir, err := os.Getwd()
	Check(err, "getting current working directory failed")
	//fmt.Println(dir)
	return dir
}

func TripleOp(condition bool, a, b interface{}) interface{} {
	if condition {
		return a
	} else {
		return b
	}
}
