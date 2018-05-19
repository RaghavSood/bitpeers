package main

import (
	"encoding/hex"
	"fmt"
	"github.com/RaghavSood/bitpeers"
	"net"
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
		fmt.Printf("SerializationVersion: 0x%s\n", hexstring(peersDB.NewAddrInfo[i].SerializationVersion))
		fmt.Printf("LastSuccess: %d\n", peersDB.NewAddrInfo[i].LastSuccess)
		fmt.Printf("Attempts: %d\n", peersDB.NewAddrInfo[i].Attempts)
		fmt.Printf("IPAddress: 0x%s\n", hexstring(peersDB.NewAddrInfo[i].Address.PeerAddress.IPAddress))
		var parsedIP net.IP
		parsedIP = peersDB.NewAddrInfo[i].Address.PeerAddress.IPAddress
		fmt.Printf("Parsed IP: %s\n", parsedIP)
		fmt.Printf("Port: %d\n", peersDB.NewAddrInfo[i].Address.PeerAddress.Port)
		fmt.Printf("Source: 0x%s\n", hexstring(peersDB.NewAddrInfo[i].Source))
		var parsedSourceIP net.IP
		parsedSourceIP = peersDB.NewAddrInfo[i].Source
		fmt.Printf("Parsed Source IP: %s\n", parsedSourceIP)

		fmt.Printf("Time: %d\n", peersDB.NewAddrInfo[i].Address.Time)
		fmt.Printf("ServiceFlags: 0x%s\n", hexstring(peersDB.NewAddrInfo[i].ServiceFlags))
		fmt.Printf("UnknownBytes: 0x%s\n", hexstring(peersDB.NewAddrInfo[i].UnknownBytes))
		fmt.Println("")
	}
	fmt.Println("Tried Addresses:")

	for i = 0; i < peersDB.NTried; i++ {
		fmt.Printf("SerializationVersion: 0x%s\n", hexstring(peersDB.TriedAddrInfo[i].SerializationVersion))
		fmt.Printf("LastSuccess: %d\n", peersDB.TriedAddrInfo[i].LastSuccess)
		fmt.Printf("Attempts: %d\n", peersDB.TriedAddrInfo[i].Attempts)
		fmt.Printf("IPAddress: 0x%s\n", hexstring(peersDB.TriedAddrInfo[i].Address.PeerAddress.IPAddress))
		var parsedIP net.IP
		parsedIP = peersDB.TriedAddrInfo[i].Address.PeerAddress.IPAddress
		fmt.Printf("Parsed IP: %s\n", parsedIP)
		fmt.Printf("Port: %d\n", peersDB.TriedAddrInfo[i].Address.PeerAddress.Port)
		fmt.Printf("Source: 0x%s\n", hexstring(peersDB.TriedAddrInfo[i].Source))
		var parsedSourceIP net.IP
		parsedSourceIP = peersDB.TriedAddrInfo[i].Source
		fmt.Printf("Parsed Source IP: %s\n", parsedSourceIP)

		fmt.Printf("Time: %d\n", peersDB.TriedAddrInfo[i].Address.Time)
		fmt.Printf("ServiceFlags: 0x%s\n", hexstring(peersDB.TriedAddrInfo[i].ServiceFlags))
		fmt.Printf("UnknownBytes: 0x%s\n", hexstring(peersDB.TriedAddrInfo[i].UnknownBytes))
		fmt.Println("")
	}
}

func hexstring(input []byte) string {
	return hex.EncodeToString(input)
}
