// The wakeonlan-go program is designed as easy WoL instrument
package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"git.mills.io/prologic/bitcask"
	"github.com/jessevdk/go-flags"
	"github.com/kirsle/configdir"

	"github.com/DanArmor/wakeonlan-go/pkg/wolpacket"
	"github.com/DanArmor/wakeonlan-go/pkg/wolrunner"
	"github.com/jedib0t/go-pretty/v6/table"
)

var opts struct {
	Dest       string `short:"d" description:"Destination address(port is optional)"`
	Local      string `short:"l" description:"Local address of network interface (port is optional)"`
	Rec        string `short:"r" description:"If used - action is not performed, but is recorded as alias to given name. If you use alias inside other alias - it will be deep copy, so you can delete used alias in the future."`
	Show       bool   `short:"s" description:"Show list of aliases"`
	Remove     string `long:"rm" description:"Remove alias with a given name"`
	Positional struct {
		MACs []string
	} `positional-args:"yes" required:"yes" description:"MAC addresses of machines to wake up"`
}

type WOLRecord struct {
	Local string
	Dest  string
	Macs  []string
}

func (wolr *WOLRecord) GetDest() string {
	if wolr.Dest == "" {
		return "Default"
	} else {
		return wolr.Dest
	}
}

func (wolr *WOLRecord) GetLocal() string {
	if wolr.Local == "" {
		return "Default"
	} else {
		return wolr.Local
	}
}

func main() {
	// Setup config file
	configPath := configdir.LocalConfig("wakeonlan-go")
	if err := configdir.MakePath(configPath); err != nil {
		fmt.Print("Error:", err)
		return
	}
	configFile := filepath.Join(configPath, "config.db")

	db, err := bitcask.Open(configFile)
	if err != nil {
		fmt.Print("Error:", err)
		return
	}
	defer db.Close()

	if _, err := flags.Parse(&opts); err != nil {
		fmt.Print("Error:", err)
		return
	}

	if opts.Show {
		ch := db.Keys()
		index := 1
		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"#", "Aliase", "Local", "Dest", "MACs"})
		for k := range ch {
			v, err := db.Get(k)
			if err != nil {
				fmt.Print("Error:", err)
				return
			}
			buf := &bytes.Buffer{}
			buf.Write(v)
			var record WOLRecord
			decoder := gob.NewDecoder(buf)
			if err := decoder.Decode(&record); err != nil {
				fmt.Print("Error:", err)
				return
			}
			t.AppendRow(table.Row{index, string(k), record.GetLocal(), record.GetDest(), strings.Join(record.Macs, "\n")})
			index++
		}
		t.SetStyle(table.StyleColoredYellowWhiteOnBlack)
		t.Render()
		return
	}

	if opts.Remove != "" {
		ch := db.Keys()
		for k := range ch {
			if string(k) == opts.Remove {
				db.Delete(k)
				fmt.Print("Alias is deleted\n")
				return
			}
		}
		fmt.Print("No such alias\n")
		return
	}

	if len(opts.Positional.MACs) == 0 {
		fmt.Print("Please, provide MAC addresses or aliases!")
		return
	}

	if opts.Rec != "" {
		// Make a record
		// Let's find other aliases in list of macs
		var addrs []string
		for _, v := range opts.Positional.MACs {
			addr, err := wolpacket.NewMACAddress(v)
			if err != nil {
				buf := &bytes.Buffer{}
				recordData, err := db.Get(addr[:])
				if err != nil {
					fmt.Print("Error:", err)
					return
				}
				buf.Write(recordData)
				var record WOLRecord
				decoder := gob.NewDecoder(buf)
				if err := decoder.Decode(&record); err != nil {
					fmt.Print("Error:", err)
					return
				}
				addrs = append(addrs, record.Macs...)
			} else {
				addrs = append(addrs, v)
			}
		}
		buf := &bytes.Buffer{}
		record := WOLRecord{
			Local: opts.Local,
			Dest:  opts.Dest,
			Macs:  addrs,
		}
		encoder := gob.NewEncoder(buf)
		if err := encoder.Encode(&record); err != nil {
			fmt.Print("Error:", err)
			return
		}
		db.Put([]byte(opts.Rec), buf.Bytes())
	} else {
		// Send WoL
		addrs := make([]WOLRecord, 1)
		addrs[0].Local = opts.Local
		addrs[0].Dest = opts.Dest
		for _, v := range opts.Positional.MACs {
			_, err := wolpacket.NewMACAddress(v)
			if err != nil {
				buf := &bytes.Buffer{}
				recordData, err := db.Get([]byte(v))
				if err != nil {
					panic(err)
				}
				buf.Write(recordData)
				var record WOLRecord
				decoder := gob.NewDecoder(buf)
				if err := decoder.Decode(&record); err != nil {
					panic(err)
				}
				addrs = append(addrs, record)
			} else {
				addrs[0].Macs = append(addrs[0].Macs, v)
			}
		}
		index := 0
		if len(addrs[0].Macs) == 0 {
			index = 1
		}
		count := 1
		for index < len(addrs) {
			wolr, err := wolrunner.NewWOLRunner(addrs[index].Local, addrs[index].Dest)
			if err != nil {
				fmt.Print("Error:", err)
				return
			}
			for i := range addrs[index].Macs {
				if err := wolr.WakeMAC(addrs[index].Macs[i]); err != nil {
					fmt.Print("Error:", err)
					return
				} else {
					fmt.Printf("#%d WoL packet was send to %s successfully. Local: %s, Dest: %s\n", count, addrs[index].Macs[i], wolr.LocalUDP().String(), wolr.DestinationUDP().String())
				}
			}
			index++
			count++
		}
	}
}
