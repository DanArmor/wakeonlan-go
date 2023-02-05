// The wakeonlan-go program is designed as easy WoL instrument
package main

import (
	"github.com/jessevdk/go-flags"
	"log"

	"github.com/DanArmor/wakeonlan-go/pkg/wolrunner"
)

var opts struct {
	Dest       string `short:"d" description:"Destination address(with port)"`
	Local      string `short:"l" description:"Local address of network interface (with port)"`
	Positional struct {
		MACs []string
	} `positional-args:"yes" required:"yes" description:"MAC addresses of machines to wake up"`
}

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatal("Error: ", err)
	}
	if len(opts.Positional.MACs) == 0 {
		log.Fatal("Please, provide MAC addresses!")
	}
	wolr, err := wolrunner.NewWOLRunner(opts.Local, opts.Dest)
	if err != nil {
		log.Fatal(err)
	}
	for i := range opts.Positional.MACs {
		if err := wolr.WakeMAC(opts.Positional.MACs[i]); err != nil {
			log.Fatal(err)
		} else {
			log.Printf("WoL packet was send to %s successfully. Local: %s, Dest: %s", opts.Positional.MACs[i], wolr.LocalUDP().String(), wolr.DestinationUDP().String())
		}
	}
}
