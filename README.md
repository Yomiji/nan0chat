# Nan0Chat
#### Secure Chat Client using Nan0 Services

This is a terminal-based chat server and client that demonstrates the [Nan0 API](https://github.com/yomiji/nan0)

#### Building
From the project root directory, do *go build* on entrypoint to create the program.
```
go build Nan0Chat nan0chat/entrypoint
```

#### Running
There are a number of command-line flags that are used to configure the application:
```
Usage of Nan0Chat:
  -host string
        Host name for server (default "localhost")
  -key string
        Encryption Key encoded in Base64.
  -port int
        Port number for server (if --server is [true]) (default 6865)
  -server
        Is this a server? [[false]/true]
  -sig string
        HMAC Signature encoded in Base64.
  -username string
        A custom user name
```
* ***server*** is a flag that indicates whether or not a server is to be started, this defaults to false
* ***host*** is the host name for the server, with the default being "localhost"
* ***port*** is the port for the server, with the default being 6865
* ***key*** is the encryption key, a 256bit string encoded in Base64, used for encryption
* ***sig*** is the signature (HMAC), a 256bit string encoded in Base64, used for authentication
* ***username*** is a custom username assigned to the client application (if a client is started)

###### Start a server:
```
./Nan0Chat  --key <encryption key> --sig <signature> --server=true --port=6865
```
Replace \<encryption key> with the encryption key and \<signature> with the signature

###### Start a client:
```
./Nan0Chat  --key <encryption key> --sig <signature> --host=localhost --server=false --port=6865 --username=Bob
```
Replace \<encryption key> with the encryption key and \<signature> with the signature

#### Obtaining an Encryption Key and Signature
To obtain an encryption key and signature in base64, run the following snippet in a console:
```go
package main

import (
	"github.com/yomiji/nan0"
	"encoding/base64"
	"fmt"
)

// A simple function to make the keys generated a sharable string
func shareKeys() (encKeyShare, authKeyShare string) {
	encKeyBytes := nan0.NewEncryptionKey()
	authKeyBytes := nan0.NewHMACKey()

	encKeyShare = base64.StdEncoding.EncodeToString(encKeyBytes[:])
	authKeyShare = base64.StdEncoding.EncodeToString(authKeyBytes[:])

	return
}

func main() {
	// create the keys
	encKey, sig := shareKeys()
	// print the keys to the console
	fmt.Printf("Encryption Key: %v\nSignature: %v", encKey, sig)
}
```
Output:

***Do NOT use these specific keys in your application***
```
Encryption Key: eQypV5GrNg0Vh+E3A90HJq+89KL/52g/IQApnUm7E+I=
Signature: Ll0RHmOxfYJi9Y5X5YaK078ZjCOHo4wBnDGenjtyqAw=
```
***Do NOT use these specific keys in your application***

#### Usage
The server application will start with a message indicating that it is currently running. The server application must be
interrupted to close. In windows command prompt and linux terminal, you can achieve this by pressing the Ctrl+C
combination.

The client application is a simple edit box below a text area. Inside the text area, there will appear all text entered
into the edit box as well as all incoming messages dispatched from the service. The messages will be prepended with the
name of the user who sent the message. Press the escape key to exit the client application.

##### TODO
* Better client disconnection handling
