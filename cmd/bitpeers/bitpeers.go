package main

import (
	"encoding/hex"
	"fmt"
	"github.com/RaghavSood/bitpeers"
)

type BitPeersDB bitpeers.PeersDB

func main() {
	path := "/Users/raghavsood/Development/go/src/github.com/RaghavSood/bitpeers/peers_newest.dat"
	rawPeersDB, err := bitpeers.NewPeersDB(path)
	if err != nil {
		fmt.Println(err)
	}

	peersDb := BitPeersDB(rawPeersDB)

	peersDb.dump()
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
