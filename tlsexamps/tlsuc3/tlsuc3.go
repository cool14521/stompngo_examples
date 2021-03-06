//
// Copyright © 2013-2015 Guy M. Allard
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

/*
Connect and Disconnect from a STOMP broker with a TLS connection, use case 3.

	TLS Use Case 3 - broker *does* authenticate client, client does *not* authenticate broker

	Subcase 3.A - Message broker configuration does *not* require client authentication

	- Expect connection success

	Subcase 3.B - Message broker configuration *does* require client authentication

	- Expect connection success if the broker can authenticate the client certificate

	Example use might be:

		go build
		./tlsuc3 -cliCertFile=/ad3/gma/sslwork/2013/client.crt -cliKeyFile=/ad3/gma/sslwork/2013/client.key

*/
package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net"
	//
	"github.com/gmallard/stompngo"
	"github.com/gmallard/stompngo_examples/sngecomm"
)

var (
	exampid     = "tlsuc3:"
	testConfig  *tls.Config
	cliCertFile string
	cliKeyFile  string
)

func init() {
	flag.StringVar(&cliCertFile, "cliCertFile", "DUMMY_CERT", "Name of client cert file")
	flag.StringVar(&cliKeyFile, "cliKeyFile", "DUMMY_KEY", "Name of client key file")
}

// Connect to a STOMP broker using TLS and disconnect.
func main() {
	fmt.Println(exampid, "starts ...")

	flag.Parse() // Parse flags

	fmt.Println(exampid, "Client Cert File", cliCertFile)
	fmt.Println(exampid, "Client Key File", cliKeyFile)

	// TLS Configuration.
	testConfig = new(tls.Config)
	testConfig.InsecureSkipVerify = true // Do *not* check the broker's certificate

	// Get host and port
	h, p := sngecomm.HostAndPort()
	fmt.Println(exampid, "host", h, "port", p)

	// Be polite, allow SNI (Server Virtual Hosting)
	testConfig.ServerName = h

	// Finish TLS Config initialization, so broker can authenticate client.
	cc, e := tls.LoadX509KeyPair(cliCertFile, cliKeyFile)
	if e != nil {
		log.Fatalln(e) // Handle this ......
	}
	testConfig.Certificates = append(testConfig.Certificates, cc) // Add cert

	// This is OK, but does not seem to be required
	testConfig.BuildNameToCertificate() // Build names map

	// Connect logic: use net.Dial and tls.Client
	t, e := net.Dial("tcp", net.JoinHostPort(h, p))
	if e != nil {
		log.Fatalln(e) // Handle this ......
	}
	fmt.Println(exampid, "dial complete ...")
	n := tls.Client(t, testConfig)
	e = n.Handshake()
	if e != nil {
		log.Fatalln(e) // Handle this ......
	}

	sngecomm.DumpTLSConfig(exampid, testConfig, n)

	fmt.Println(exampid, "handshake complete ...")

	// Connect Headers
	ch := sngecomm.ConnectHeaders()

	// Get a stomp connection.  Parameters are:
	// a) the opened net connection
	// b) the connect Headers
	conn, e := stompngo.Connect(n, ch)
	if e != nil {
		log.Fatalln(e) // Handle this ......
	}
	fmt.Println(exampid, "stomp connect complete ...")

	// *NOTE* your application functionaltiy goes here!

	// Polite Stomp disconnects are not required, but highly recommended.
	// Empty headers here.
	e = conn.Disconnect(stompngo.Headers{})
	if e != nil {
		log.Fatalln(e) // Handle this ......
	}
	fmt.Println(exampid, "stomp disconnect complete ...")

	// Close the net connection.
	e = n.Close()
	if e != nil {
		log.Fatalln(e) // Handle this ......
	}
	fmt.Println(exampid, "network close complete ...")

	fmt.Println(exampid, "ends ...")
}
