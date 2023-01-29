package main

import (
	"bytes"
	"encoding/binary"
	"log"
	"net"
	"os"
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

func wakeMAC(packet WOLPacket) {
	LocalAddr, err := net.ResolveUDPAddr("udp", getLocalAddress().String()+":20")
	if err != nil {
		log.Fatal(err)
	}
	d := net.Dialer{LocalAddr: LocalAddr}
	conn, err := d.Dial("udp", "255.255.255.255:9")
	if err != nil {
		log.Fatal("Erro during establishing connection!", err)
	}
	defer conn.Close()
	n, err := conn.Write(packet.Marshal())
	if err != nil {
		log.Fatal("Error during sending data!", err)
	}
	log.Printf("Transfered %d bytes", n)
}

func getLocalAddress() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP
}

func main() {
	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) == 0 {
		log.Fatal("Not enough args!")
	} else if len(argsWithoutProg) > 1 {
		log.Fatal("Too many args!")
	}
	wolPacket := NewWOLPacket(argsWithoutProg[0])
	wakeMAC(wolPacket)
}
