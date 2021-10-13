package network

import (
	"fuze/srv/network/dto"
	"net"
)

// FuzeTransmission is the network interface: LAN  WAN (internet).
type FileRetriever interface {
	// Checking if the server retrieve any inviation from localnetwork.
	// Note: udp protocol will not trigger this step, therefore udp' user can proceed to retrieve.
	HasInvitation(conn net.Conn) (invitation *dto.Invitation, err error)

	// Retrieve data and save with format into the given path.
	// The conn param will be closed after the retrieving process.
	Retrieve(filename string, conn net.Conn) error
}

type FileSender interface{}
