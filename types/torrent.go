package types

//import "time"

type TorrentFile struct {
	Announce    string
	InfoHash    [20]byte //single slice
	PieceHashes [][20]byte
	PieceLength int
	Length      int
	Name        string
}

//response from tracker
type TrackerResponse struct {
	Interval int
	Peers    []Peer
}

//Peer structure
type Peer struct {
	IP   string
	Port uint16
}

// communication between Peers
type PeerMessage struct {
	ID      uint8
	Payload []byte
}

//work to be done for a piece (work queue item)
type PieceWork struct {
	Index  int
	Hash   [20]byte //SHA-1 check
	Length int
}

//piece Result
type PieceResult struct {
	Index int
	Data  []byte
}
