package main

import (
	"encoding/hex"
	"fmt"
	"github.com/RaghavSood/bitpeers"
)

type BitPeersDB bitpeers.PeersDB

func main() {
	path := "/Users/raghavsood/Development/go/src/github.com/RaghavSood/bitpeers/peers.dat"
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
	fmt.Printf("Version: 0x%s\n", hexstring([]byte{peersDB.Version}))
}

func hexstring(input []byte) string {
	return hex.EncodeToString(input)
}
