package nan0chat

import (
	"encoding/base64"
	"fmt"
	"flag"
	"math/rand"
	"time"
)

// create a new random number generator
var random = rand.New(rand.NewSource(time.Now().Unix()))

// application parameter flags definitions
var EncryptKey = flag.String("key", "", "(Required) Encryption Key encoded in Base64")
var Signature = flag.String("sig", "", "(Required) HMAC Signature encoded in Base64")
var IsServer = flag.Bool("server", false, "Is this a server?")
var Host = flag.String("host", "localhost", "Host name for server")
var Port = flag.Int("port", 6865, "Port number for server (if --server is [true])")
var CustomUsername = flag.String("username", "", "A custom user name")

// The nan0 functions require a specific key type and width, this is a way to make
// that conversion from strings to the required type.
func KeysToNan0Bytes(encKeyShare, authKeyShare string) (encKey, authKey *[32]byte) {
	encKeyBytes, _ := base64.StdEncoding.DecodeString(encKeyShare)
	authkeyBytes, _ := base64.StdEncoding.DecodeString(authKeyShare)

	encKey = &[32]byte{}
	authKey = &[32]byte{}
	copy(encKey[:], encKeyBytes)
	copy(authKey[:], authkeyBytes)

	return
}

// Handles errors generated in the application
func handleErr(e error, exec func(err error) interface{}) interface{} {
	if e != nil {
		fmt.Printf("Error occurred: %v\n", e)
		if exec != nil {
			return exec(e)
		}
	}

	return e
}

func (user *User) SetUserName(name string) {
	user.UserName = name
}

func (user *User) SetUserId(id int64) {
	user.UserId = id
}
