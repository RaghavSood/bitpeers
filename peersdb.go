package bitpeers

import (
	"fmt"
	"io/ioutil"
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
	Source      []byte
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
	IPAddress []byte
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

		peersDB.NewAddrInfo[i].SerializationVersion = dbreader.readBytes(4)
		peersDB.NewAddrInfo[i].Address.Time = dbreader.readUint32()
		peersDB.NewAddrInfo[i].ServiceFlags = dbreader.readBytes(8)

		peersDB.NewAddrInfo[i].Address.PeerAddress.IPAddress = dbreader.readBytes(16)
		peersDB.NewAddrInfo[i].Address.PeerAddress.Port = dbreader.readBigEndianUint16()
		peersDB.NewAddrInfo[i].Source = dbreader.readBytes(16)

		peersDB.NewAddrInfo[i].LastSuccess = dbreader.readUint64()
		peersDB.NewAddrInfo[i].Attempts = dbreader.readUint32()
	}

	for i = 0; i < peersDB.NTried; i++ {

		peersDB.TriedAddrInfo[i].SerializationVersion = dbreader.readBytes(4)
		peersDB.TriedAddrInfo[i].Address.Time = dbreader.readUint32()
		peersDB.TriedAddrInfo[i].ServiceFlags = dbreader.readBytes(8)

		peersDB.TriedAddrInfo[i].Address.PeerAddress.IPAddress = dbreader.readBytes(16)
		peersDB.TriedAddrInfo[i].Address.PeerAddress.Port = dbreader.readBigEndianUint16()
		peersDB.TriedAddrInfo[i].Source = dbreader.readBytes(16)

		peersDB.TriedAddrInfo[i].LastSuccess = dbreader.readUint64()
		peersDB.TriedAddrInfo[i].Attempts = dbreader.readUint32()
	}

	return peersDB, nil
}

func readDBBytes(peersDB PeersDB) ([]byte, error) {
	return ioutil.ReadFile(peersDB.Path)
}
