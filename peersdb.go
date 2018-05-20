package bitpeers

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
)

type PeersDB struct {
	Path          string      `json:"-"`
	MessageBytes  []byte      `json:"message_bytes"` // 0  : 4
	Version       uint8       `json:"version"`       // 4  : 4
	KeySize       uint8       `json:"keysize"`       // 5  : 5
	NKey          []byte      `json:"nkey"`          // 37 : 32
	NNew          uint32      `json:"nnew"`          // 41 : 4
	NTried        uint32      `json:"ntried"`        // 45 : 4
	NewBuckets    uint32      `json:"new_buckets"`   // 49 : 4
	NewAddrInfo   []CAddrInfo `json:"new_addr_info"`
	TriedAddrInfo []CAddrInfo `json:"tried_addr_info"`
}

type CAddrInfo struct {
	Address     CAddress `json:"address"`
	Source      net.IP   `json:"source"`
	LastSuccess uint64   `json:"last_success"`
	Attempts    uint32   `json:"attempts"`
}

type CAddress struct {
	SerializationVersion []byte   `json:"serialization_version"`
	Time                 uint32   `json:"time"`
	ServiceFlags         []byte   `json:"service_flags"`
	PeerAddress          CService `json:"ip"`
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

func (cAddress *CAddress) MarshalJSON() ([]byte, error) {
	type Alias CAddress
	return json.Marshal(&struct {
		IP                   string `json:"ip"`
		SerializationVersion string `json:"serialization_version"`
		ServiceFlags         string `json:"service_flags"`
		*Alias
	}{
		IP:                   cAddress.PeerAddress.String(),
		SerializationVersion: hexstring(cAddress.SerializationVersion),
		ServiceFlags:         binaryString(cAddress.ServiceFlags),
		Alias:                (*Alias)(cAddress),
	})
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

func binaryString(input []byte) string {
	var binaryString string
	for _, x := range input {
		binaryString += fmt.Sprintf("%08b", x)
	}
	return binaryString
}
