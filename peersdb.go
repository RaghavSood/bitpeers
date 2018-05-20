package bitpeers

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net"
)

type PeersDB struct {
	Path          string
	MessageBytes  []byte // 0  : 4
	Version       uint8  // 4  : 4
	KeySize       uint8  // 5  : 5
	NKey          []byte // 37 : 32
	NNew          uint32 // 41 : 4
	NTried        uint32 // 45 : 4
	NewBuckets    uint32 // 49 : 4
	NewAddrInfo   []CAddrInfo
	TriedAddrInfo []CAddrInfo
}

type CAddrInfo struct {
	Address     CAddress
	Source      net.IP
	LastSuccess uint64
	Attempts    uint32
}

type CAddress struct {
	SerializationVersion []byte
	Time                 uint32
	ServiceFlags         []byte
	PeerAddress          CService
}

type CService struct {
	IPAddress net.IP
	Port      uint16 // This is serialized as BigEndian
}

func NewPeersDB(path string) (PeersDB, error) {
	peersDB := PeersDB{
		Path: path,
	}

	dbbytes, err := readDBBytes(peersDB)
	if err != nil {
		return peersDB, fmt.Errorf("Couldn't read peer file %s", peersDB.Path)
	}

	dbreader := DBReader{
		Bytes:  dbbytes,
		Cursor: 0,
	}

	peersDB.MessageBytes = dbreader.readBytes(4)
	peersDB.Version = dbreader.readUint8()
	peersDB.KeySize = dbreader.readUint8()
	peersDB.NKey = dbreader.readBytes(32)                  // uint256 type
	peersDB.NNew = dbreader.readUint32()                   // int type
	peersDB.NTried = dbreader.readUint32()                 // int type
	peersDB.NewBuckets = dbreader.readUint32() ^ (1 << 30) // int type

	peersDB.NewAddrInfo = make([]CAddrInfo, peersDB.NNew)
	peersDB.TriedAddrInfo = make([]CAddrInfo, peersDB.NTried)

	var i uint32
	for i = 0; i < peersDB.NNew; i++ {
		peersDB.NewAddrInfo[i] = dbreader.readCAddrInfo()
	}

	for i = 0; i < peersDB.NTried; i++ {
		peersDB.TriedAddrInfo[i] = dbreader.readCAddrInfo()
	}

	return peersDB, nil
}

func (dbreader *DBReader) readCAddrInfo() (cAddrInfo CAddrInfo) {
	cAddrInfo.Address.SerializationVersion = dbreader.readBytes(4)
	cAddrInfo.Address.Time = dbreader.readUint32()
	cAddrInfo.Address.ServiceFlags = reverseBytes(dbreader.readBytes(8))
	cAddrInfo.Address.PeerAddress.IPAddress = dbreader.readBytes(16)
	cAddrInfo.Address.PeerAddress.Port = dbreader.readBigEndianUint16()

	cAddrInfo.Source = dbreader.readBytes(16)
	cAddrInfo.LastSuccess = dbreader.readUint64()
	cAddrInfo.Attempts = dbreader.readUint32()
	return
}

func readDBBytes(peersDB PeersDB) ([]byte, error) {
	return ioutil.ReadFile(peersDB.Path)
}

func (cAddrInfo CAddrInfo) String() string {
	return fmt.Sprintf("%s\nSource: %s\nLastSuccess: %d\nAttempts: %d\n\n", cAddrInfo.Address, cAddrInfo.Source, cAddrInfo.LastSuccess, cAddrInfo.Attempts)
}

func (cAddress CAddress) String() string {
	return fmt.Sprintf("SerializationVersion: %s\nTime: %d\nServiceFlags: 0x%s\nIP: %s", hexstring(cAddress.SerializationVersion), cAddress.Time, hexstring(cAddress.ServiceFlags), cAddress.PeerAddress)
}

func (cService CService) String() string {
	return fmt.Sprintf("%s:%d", cService.IPAddress, cService.Port)
}

func hexstring(input []byte) string {
	return hex.EncodeToString(input)
}

func reverseBytes(input []byte) []byte {
	for i, j := 0, len(input)-1; i < j; i, j = i+1, j-1 {
		input[i], input[j] = input[j], input[i]
	}
	return input
}
