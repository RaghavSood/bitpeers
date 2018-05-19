package bitpeers

import (
	"bytes"
	"encoding/binary"
	"log"
)

const length_UINT8 = 1
const length_UINT16 = 2
const length_UINT32 = 4
const length_UINT64 = 8

type DBReader struct {
	Bytes  []byte
	Cursor uint64
}

func (r *DBReader) readUint8() uint8 {
	val := uint8(0)
	buf := bytes.NewBuffer(r.Bytes[r.Cursor : r.Cursor+length_UINT8])
	err := binary.Read(buf, binary.LittleEndian, &val)
	if err != nil {
		log.Fatalf("Decode failed: %s", err)
	}
	r.Cursor += length_UINT8
	return val
}

func (r *DBReader) readUint16() uint16 {
	val := uint16(0)
	buf := bytes.NewBuffer(r.Bytes[r.Cursor : r.Cursor+length_UINT16])
	err := binary.Read(buf, binary.LittleEndian, &val)
	if err != nil {
		log.Fatalf("Decode failed: %s", err)
	}
	r.Cursor += length_UINT16
	return val
}

func (r *DBReader) readBigEndianUint16() uint16 {
	val := uint16(0)
	buf := bytes.NewBuffer(r.Bytes[r.Cursor : r.Cursor+length_UINT16])
	err := binary.Read(buf, binary.BigEndian, &val)
	if err != nil {
		log.Fatalf("Decode failed: %s", err)
	}
	r.Cursor += length_UINT16
	return val
}

func (r *DBReader) readUint32() uint32 {
	val := uint32(0)
	buf := bytes.NewBuffer(r.Bytes[r.Cursor : r.Cursor+length_UINT32])
	err := binary.Read(buf, binary.LittleEndian, &val)
	if err != nil {
		log.Fatalf("Decode failed: %s", err)
	}
	r.Cursor += length_UINT32
	return val
}

func (r *DBReader) readUint64() uint64 {
	val := uint64(0)
	buf := bytes.NewBuffer(r.Bytes[r.Cursor : r.Cursor+length_UINT64])
	err := binary.Read(buf, binary.LittleEndian, &val)
	if err != nil {
		log.Fatalf("Decode failed: %s", err)
	}
	r.Cursor += length_UINT64
	return val
}

func (r *DBReader) readByte() byte {
	byteVal := r.Bytes[r.Cursor]
	r.Cursor += 1
	return byteVal
}

func (r *DBReader) readBytes(length uint64) []byte {
	byteVals := r.Bytes[r.Cursor : r.Cursor+length]
	r.Cursor += length
	return byteVals
}

// Allows you to view data without moving the cursor and from a start point.
// Useful for cases to lookahead on data, such as checking if a tx is a
// segwit tx

func (r *DBReader) peekBytesFrom(start uint64, length uint64) []byte {
	byteVals := r.Bytes[start : start+length]
	return byteVals
}

// Allows you to view data without moving the cursor. Useful for cases to
// lookahead on data, such as checking if a tx is a segwit tx

func (r *DBReader) peekBytes(length uint64) []byte {
	byteVals := r.Bytes[r.Cursor : r.Cursor+length]
	return byteVals
}
