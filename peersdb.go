package bitpeers

import (
	"fmt"
	"io/ioutil"
)

type PeersDB struct {
	Path         string
	MessageBytes []byte
	Version      byte
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
	peersDB.Version = dbreader.readByte()

	return peersDB, nil
}

func readDBBytes(peersDB PeersDB) ([]byte, error) {
	return ioutil.ReadFile(peersDB.Path)
}
