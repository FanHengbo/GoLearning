package ftpserver

import (
	"net"
)

type FTPserver struct {
	listener    net.Listener
	clientCount int32
}
