package network

import (
	"bufio"
	"fmt"
	"fuze/srv/network/dto"

	"io"
	"log"
	"net"
	"os"
)

var (
	InvitationReceive = "InvitationReceive"
)

type Receiver struct{}

func NewReceiver() *Receiver {
	return &Receiver{}
}

func (re *Receiver) HasInvitation(conn net.Conn) (invitation *dto.Invitation, err error) {
	var (
		buf = make([]byte, 1024)
		r   = bufio.NewReader(conn)
		w   = bufio.NewWriter(conn)
	)

	for {
		_, err = r.Read(buf)
		if (err != io.EOF) && (err != nil) {
			log.Fatalf("Error reading invitation: %v", err)
			return invitation, err
		}

		invi := dto.NewInvitation()
		if err = invi.BindAndValidate(buf); err != nil {
			continue
		}

		w.Write([]byte(InvitationReceive))
		w.Flush()

		return invi, nil
	}
}

// Retrieve eventually stop the listening process.
func (re *Receiver) Retrieve(filename string, conn net.Conn) error {
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer f.Close()

	var (
		buf       = make([]byte, 1024)
		r         = bufio.NewReader(conn)
		totalSize = 0
	)

	for {
		read_size, err := r.Read(buf)
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("failed reading file with error %v", err)
		}

		// read only the retrieved bytes.
		size, err := f.Write(buf[:read_size])
		if err != nil {
			log.Fatalf("failed write file with error %v", err)
		}

		totalSize += size
	}

	log.Printf("Retrieved successfully %v bytes", totalSize)
	conn.Close()

	return nil
}
