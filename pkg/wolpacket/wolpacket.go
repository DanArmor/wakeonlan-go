// Package wolpacket provides functionality for creating
// MAC address struct out of string and WoL Magic Packet out of string
package wolpacket

import (
	"bytes"
	"encoding/binary"
	"net"
)

// MACAddress provides MAC as array of bytes
type MACAddress [6]byte

// WOLPacket representing Magic WoL Packet
type WOLPacket struct {
	header  [6]byte        // Synchronization bytes 0xFF 6 times
	payload [16]MACAddress // Payload of Magic Packet - MAC repeated 16 times
}

// NewMACAddress creating MACAddress out of string
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

// NewWOLPacket creating Magic WoL Packet out of string
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

// Marshal converting Magic WoL packet into slice of bytes in BigEndian form
func (wp *WOLPacket) Marshal() ([]byte, error) {
	var buf bytes.Buffer
	if err := binary.Write(&buf, binary.BigEndian, wp); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
