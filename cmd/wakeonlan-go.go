// The wakeonlan-go program is designed as easy WoL instrument
package main

import (
	"flag"
	"log"

	"github.com/DanArmor/wakeonlan-go/pkg/wolrunner"
)

func main() {
	macFlag := flag.String("m", "", "MAC address (48 bit)")
	localFlag := flag.String("l", "", "Local address (with port)")
	destinationFlag := flag.String("d", "", "Destination address (with port)")
	flag.Parse()
	if *macFlag == "" {
		log.Fatal("Please, provide MAC address with -m flag!")
	}
	wolr, err := wolrunner.NewWOLRunner(*localFlag, *destinationFlag)
	if err != nil {
		log.Fatal(err)
	}
	if err := wolr.WakeMAC(*macFlag); err != nil {
		log.Fatal(err)
	}
	log.Printf("WoL packet was send to %s successfully. Local: %s, Dest: %s", *macFlag, wolr.LocalUDP().String(), wolr.DestinationUDP().String())
}
