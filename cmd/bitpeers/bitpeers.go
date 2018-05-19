package main

import (
	"fmt"
	"github.com/RaghavSood/bitpeers"
)

type BitPeersDB bitpeers.PeersDB

func main() {
	path := "/Users/raghavsood/Development/go/src/github.com/RaghavSood/bitpeers/peers.dat"
	peersDb := BitPeersDB(bitpeers.NewPeersDB(path))
	peersDb.dump()
}

func (peersDB BitPeersDB) dump() {
	fmt.Printf("Path: %s\n", peersDB.Path)
}
