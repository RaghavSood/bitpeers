package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/RaghavSood/bitpeers"
	flag "github.com/spf13/pflag"
	"os"
)

type BitPeersDB bitpeers.PeersDB

var peersFilePath string
var formatOption string
var addressOnly bool

func init() {
	flag.StringVar(&peersFilePath, "filepath", "", "the path to peers.dat")
	flag.StringVar(&formatOption, "format", "json", "the output format {json|text}")
	flag.BoolVar(&addressOnly, "addressonly", false, "outputs only addresses if specified")
	flag.Parse()
}

func main() {
	if peersFilePath == "" {
		fmt.Fprintf(os.Stderr, "Invalid peers file %s\n", peersFilePath)
		os.Exit(1)
	}

	if formatOption != "json" && formatOption != "text" {
		fmt.Fprintf(os.Stderr, "Invalid output format %s\n", formatOption)
		os.Exit(1)
	}

	rawPeersDB, err := bitpeers.NewPeersDB(peersFilePath)
	if err != nil {
		fmt.Println(err)
	}

	peersDb := BitPeersDB(rawPeersDB)

	if addressOnly {
		addressArray := make([]string, peersDb.NTried+peersDb.NNew)
		var i uint32
		for i = 0; i < peersDb.NNew; i++ {
			addressArray[i] = peersDb.NewAddrInfo[i].Address.PeerAddress.String()
		}
		for i = 0; i < peersDb.NTried; i++ {
			addressArray[peersDb.NNew+i] = peersDb.TriedAddrInfo[i].Address.PeerAddress.String()
		}

		if formatOption == "text" {
			for _, i := range addressArray {
				fmt.Printf("%s\n", i)
			}
			return
		} else {
			encodedPeers, err := json.Marshal(addressArray)
			if err != nil {
				fmt.Sprintf("Error converting to JSON: %s", err)
				os.Exit(1)
			}
			fmt.Println(string(encodedPeers))
			return
		}
	}

	if formatOption == "text" {
		peersDb.dump()
	} else {
		encodedPeers, err := json.Marshal(peersDb)
		if err != nil {
			fmt.Sprintf("Error converting to JSON: %s", err)
			os.Exit(1)
		}
		fmt.Println(string(encodedPeers))
	}

}

func (peersDB BitPeersDB) dump() {
	fmt.Println("bitpeers")
	fmt.Println("--------")
	fmt.Printf("Path: %s\n", peersDB.Path)
	fmt.Printf("MessageBytes: 0x%s\n", hexstring(peersDB.MessageBytes))
	fmt.Printf("Version: %d\n", peersDB.Version)
	fmt.Printf("KeySize: %d\n", peersDB.KeySize)
	fmt.Printf("NKey: %s\n", hexstring(peersDB.NKey))
	fmt.Printf("NNew: %d\n", peersDB.NNew)
	fmt.Printf("NTried: %d\n", peersDB.NTried)
	fmt.Printf("NewBuckets: %d\n", peersDB.NewBuckets)
	fmt.Println("")
	var i uint32
	for i = 0; i < peersDB.NNew; i++ {
		fmt.Print(peersDB.NewAddrInfo[i])
	}
	fmt.Println("Tried Addresses:")
	for i = 0; i < peersDB.NTried; i++ {
		fmt.Print(peersDB.TriedAddrInfo[i])
	}
}

func hexstring(input []byte) string {
	return hex.EncodeToString(input)
}
