package wolpacket

import (
	"bytes"
	"encoding/binary"
	"net"
)

type MACAddress [6]byte

type WOLPacket struct {
	header  [6]byte
	payload [16]MACAddress
}

func NewMACAddress(mac string) (MACAddress, error) {
	destinationAddr, err := net.ParseMAC(mac)
	if err != nil {
		return MACAddress{}, err
	}
	var macAddr MACAddress
	for i := range macAddr {
		macAddr[i] = destinationAddr[i]
	}
	return macAddr, nil
}

func NewWOLPacket(mac string) (WOLPacket, error) {
	macAddr, err := NewMACAddress(mac)
	if err != nil {
		return WOLPacket{}, err
	}
	var wolPacket WOLPacket
	for i := range wolPacket.header {
		wolPacket.header[i] = 0xFF
	}
	for i := range wolPacket.payload {
		wolPacket.payload[i] = macAddr
	}
	return wolPacket, nil
}

func (wp *WOLPacket) Marshal() ([]byte, error) {
	var buf bytes.Buffer
	if err := binary.Write(&buf, binary.BigEndian, wp); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
