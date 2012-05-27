//
// Copyright © 2011 Guy M. Allard
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
Send messages to a STOMP 1.0 broker.
*/
package main

import (
	"fmt"
	"github.com/gmallard/stompngo"
	"github.com/gmallard/stompngo_examples/sngecomm"
	"log"
	"net"
)

var exampid = "send_10: "

// Connect to a STOMP 1.0 broker, send some messages and disconnect.
func main() {
	fmt.Println(exampid + "starts ...")

	// Open a net connection
	h, p := sngecomm.HostAndPort10()
	n, e := net.Dial("tcp", net.JoinHostPort(h, p))
	if e != nil {
		log.Fatalln(e) // Handle this ......
	}
	fmt.Println(exampid + "dial complete ...")

	eh := stompngo.Headers{}
	conn, e := stompngo.Connect(n, eh)
	if e != nil {
		log.Fatalln(e) // Handle this ......
	}
	fmt.Println(exampid + "stomp connect complete ...")

	// *NOTE* your application functionaltiy goes here!
	s := stompngo.Headers{"destination", sngecomm.Dest()} // send headers
	m := exampid + " message: "
	for i := 1; i <= sngecomm.Nmsgs(); i++ {
		t := m + fmt.Sprintf("%d", i)
		e := conn.Send(s, t)
		if e != nil {
			log.Fatalln(e) // Handle this ...
		}
		fmt.Println(exampid, "send complete:", t)
	}

	// Disconnect from the Stomp server
	e = conn.Disconnect(eh)
	if e != nil {
		log.Fatalln(e) // Handle this ......
	}
	fmt.Println(exampid + "stomp disconnect complete ...")
	// Close the network connection
	e = n.Close()
	if e != nil {
		log.Fatalln(e) // Handle this ......
	}
	fmt.Println(exampid + "network close complete ...")

	fmt.Println(exampid + "ends ...")
}
