package ulid

import (
	"crypto/rand"
	"encoding/binary"
	"time"
)

// A ULID is a 16 byte Universally Unique Lexicographically Sortable Identifier
type ULID [16]byte

const (
	// Crockford"s Base32
	// https://en.wikipedia.org/wiki/Base32
	alphabet          = "0123456789ABCDEFGHJKMNPQRSTVWXYZ"
	alphabetSize      = int64(len(alphabet))
	encodedTimeLength = 10
	encodedRandLength = 16
)

var (
	rander = rand.Reader // random function
	Nil    ULID
)

func New() ULID {
	return Must(NewRandom())
}

// Must returns ulid if err is nil and panics otherwise.
func Must(ulid ULID, err error) ULID {
	if err != nil {
		panic(err)
	}
	return ulid
}

// NewRandom returns a ULID (binary implementation) or panics.
//
// The strength of the ULIDs is based on the strength of the crypto/rand
// package.
func NewRandom() (ULID, error) {
	var (
		ulid ULID
	)
	err := encodeRandom(&ulid)
	if err != nil {
		return Nil, err
	}
	encodeTime(&ulid, time.Now())
	return ulid, err
}

func encodeTime(ulid *ULID, t time.Time) {
	var v [8]byte
	timestamp := uint64(t.UnixNano() / int64(time.Millisecond))
	binary.LittleEndian.PutUint64(v[:], timestamp)
	// Truncates at the 6th byte as designed in the original spec (48 bytes).
	copy(ulid[:], v[:6])
}

func encodeRandom(ulid *ULID) (err error) {
	_, err = rand.Read(ulid[6:])
	return
}

// String returns the string form of ulid (26 characters, non-standard base 32)
func (ulid ULID) String() string {
	var (
		buf [26]byte
		v   [8]byte
	)
	copy(v[:], ulid[:6])
	timestamp := int64(binary.LittleEndian.Uint64(v[:]))
	for x := encodedTimeLength - 1; x >= 0; x-- {
		mod := timestamp % alphabetSize
		buf[x] = alphabet[mod]
		timestamp = (timestamp - mod) / alphabetSize
	}
	for x := encodedTimeLength; x < len(ulid); x++ {
		buf[x] = alphabet[int64(ulid[x])%alphabetSize]
	}
	return string(buf[:])
}
