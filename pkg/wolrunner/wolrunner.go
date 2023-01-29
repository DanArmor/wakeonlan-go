package wolrunner

import (
	"github.com/DanArmor/wakeonlan-go/pkg/wolpacket"
	"log"
	"net"
)

type WOLRunner struct {
	localAddr net.Addr
}

func NewWOLRunner() WOLRunner{
	return WOLRunner{}
}

func (wolr *WOLRunner) WakeMAC(mac string) {
	packet := wolpacket.NewWOLPacket(mac)
	wolr.wakeMAC(packet)
}

func (wolr *WOLRunner) wakeMAC(packet wolpacket.WOLPacket) {
	LocalAddr, err := net.ResolveUDPAddr("udp", wolr.getLocalAddress().String()+":20")
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

func (wolr *WOLRunner) getLocalAddress() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP
}
