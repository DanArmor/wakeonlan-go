package wolrunner

import (
	"github.com/DanArmor/wakeonlan-go/pkg/wolpacket"
	"net"
)

const (
	broadcastIPv4     = "255.255.255.255"
	defaultWOLPort    = ":9"
	googleDnsIPv4     = "8.8.8.8"
	defaultDnsUdpPort = ":53"
	anyAvailablePort  = ":0"
)

type WOLRunner struct {
	localUDP       *net.UDPAddr
	destinationUDP *net.UDPAddr
}

func NewWOLRunner(localAddr string, destinationAddr string) (WOLRunner, error) {
	if localAddr == "" {
		localAddrIP, err := getLocalAddress()
		if err != nil {
			return WOLRunner{}, err
		}
		localAddr = localAddrIP.String() + anyAvailablePort
	}
	localUDP, err := net.ResolveUDPAddr("udp", localAddr)
	if err != nil {
		return WOLRunner{}, err
	}
	if destinationAddr == "" {
		destinationAddr = broadcastIPv4 + defaultWOLPort
	}
	destinationUDP, err := net.ResolveUDPAddr("udp", destinationAddr)
	if err != nil {
		return WOLRunner{}, err
	}
	return WOLRunner{
		localUDP:       localUDP,
		destinationUDP: destinationUDP,
	}, nil
}

func (wolr *WOLRunner) WakeMAC(mac string) error {
	packet, err := wolpacket.NewWOLPacket(mac)
	if err != nil {
		return err
	}
	return wolr.wakeMAC(packet)
}

func (wolr *WOLRunner) wakeMAC(packet wolpacket.WOLPacket) error {
	conn, err := net.DialUDP("udp", wolr.localUDP, wolr.destinationUDP)
	if err != nil {
		return err
	}
	defer conn.Close()
	bytes, err := packet.Marshal()
	if err != nil {
		return err
	}
	_, err = conn.Write(bytes)
	if err != nil {
		return err
	}
	return nil
}

func getLocalAddress() (net.IP, error) {
	conn, err := net.Dial("udp", googleDnsIPv4+defaultDnsUdpPort)
	if err != nil {
		return net.IP{}, err
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP, nil
}
