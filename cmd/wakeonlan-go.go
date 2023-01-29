// The wakeonlan-go program is designed as easy WoL instrument
package main

import (
	"github.com/DanArmor/wakeonlan-go/pkg/wolrunner"
	"log"
	"os"
)

func main() {
	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) == 0 {
		log.Fatal("Not enough args!")
	} else if len(argsWithoutProg) > 1 {
		log.Fatal("Too many args!")
	}
	wolr, err := wolrunner.NewWOLRunner("", "")
	if err != nil {
		log.Fatal(err)
	}
	if err := wolr.WakeMAC(argsWithoutProg[0]); err != nil {
		log.Fatal(err)
	}
}
