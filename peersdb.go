package bitpeers

type PeersDB struct {
	Path         string
	MessageBytes []byte
	Version      []byte
}

func NewPeersDB(path string) PeersDB {
	peersDB := PeersDB{
		Path: path,
	}

	return peersDB
}
