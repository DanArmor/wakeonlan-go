package wolpacket

import (
	"net"
	"log"
	"encoding/binary"
	"bytes"
)

type MACAddress [6]byte

type WOLPacket struct {
	header  [6]byte
	payload [16]MACAddress
}

func NewMACAddress(mac string) MACAddress {
	destinationAddr, err := net.ParseMAC(mac)
	if err != nil {
		log.Fatal("Argument", mac, "is not a mac address!", err)
	}
	log.Printf("Mac address: %v", destinationAddr)
	var macAddr MACAddress
	for i := range macAddr {
		macAddr[i] = destinationAddr[i]
	}
	return macAddr
}

func NewWOLPacket(mac string) WOLPacket {
	macAddr := NewMACAddress(mac)
	var wolPacket WOLPacket
	for i := range wolPacket.header {
		wolPacket.header[i] = 0xFF
	}
	for i := range wolPacket.payload {
		wolPacket.payload[i] = macAddr
	}
	return wolPacket
}

func (wp *WOLPacket) Marshal() []byte {
	var buf bytes.Buffer
	if err := binary.Write(&buf, binary.BigEndian, wp); err != nil {
		log.Fatal("Error during marshaling! ", err)
	}
	log.Printf("Buffer: %v", buf.Bytes())
	return buf.Bytes()
}
