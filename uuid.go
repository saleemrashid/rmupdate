package rmupdate

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

type UUID [16]byte

func NewUUID() UUID {
	var uuid UUID
	if _, err := rand.Read(uuid[:]); err != nil {
		panic(err)
	}
	// Set UUID version to 4
	uuid[6] = (uuid[6] & 0x0f) | 0x40
	// Set UUID variant to DCE
	uuid[8] = (uuid[8] & 0x3f) | 0x80
	return uuid
}

func (uuid UUID) String() string {
	return fmt.Sprintf("{%x-%x-%x-%x-%x}",
		uuid[:4],
		uuid[4:6],
		uuid[6:8],
		uuid[8:10],
		uuid[10:])
}

func NewMachineID() string {
	uuid := NewUUID()
	return hex.EncodeToString(uuid[:])
}
