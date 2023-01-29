package main

import (
	"log"
	"os"
	"github.com/DanArmor/wakeonlan-go/pkg/wolrunner"
)


func main() {
	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) == 0 {
		log.Fatal("Not enough args!")
	} else if len(argsWithoutProg) > 1 {
		log.Fatal("Too many args!")
	}
	wolr := &wolrunner.WOLRuner{}
	wolr.WakeMAC(argsWithoutProg[0])
}
