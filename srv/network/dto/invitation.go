package dto

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"path/filepath"
	"strings"
)

type ErrBindJSON error

type Invitation struct {
	SenderIP string
	Filename string
}

func (in *Invitation) BindAndValidate(invi []byte) error {
	err := json.Unmarshal(invi, in)
	if err != nil {
		return err
	}

	if net.ParseIP(in.SenderIP) == nil {
		return errors.New("invalid IP address")
	}

	fileExtension := filepath.Ext(in.Filename)
	if len(fileExtension) == 0 {
		return errors.New("invalid File extension")
	}

	i := string(invi)
	if !strings.Contains(i, "request2send") {
		return errors.New("invalid Request")
	}

	return nil
}

func (in *Invitation) Valid() bool {
	if len(in.Filename) != 0 && len(in.SenderIP) != 0 {
		return true
	}

	return false
}

func (in *Invitation) Format() string {
	return fmt.Sprintf("Invitattion: IP %v in your Local Network wants to send you this package: %v \n", in.SenderIP, in.Filename)
}

func (in *Invitation) Encode() string {
	return fmt.Sprintf("request2send-%v-%v", in.SenderIP, in.Filename)
}

func NewInvitation() *Invitation {
	return &Invitation{}
}
