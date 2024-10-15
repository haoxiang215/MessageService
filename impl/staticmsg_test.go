/*
Copyright 2021 Ethan Blanton <eblanton@buffalo.edu>

This file is part of a CSE 486/586 project from the University at
Buffalo.  Distribution of this file or its associated repository
requires the written permission of Ethan Blanton.  Sharing this file
may be a violation of academic integrity, please consult the course
policies for more package.
*/

package impl

import (
	"bytes"
	"net"
	"testing"

	"cse586.messageservice/given/directory"
)

// The following variables represent a correctly-formed message, both
// in its constituent parts and as encoded.

// staticMsgSender is the sender of the message staticMsg
var staticMsgSender string = "gray"

// staticMsgRecipient is the recipient of the message staticMsg
var staticMsgRecipient string = "lynch"

// staticMsgText is the body of the message staticMsg; it is the
// string "Hello, world!"
var staticMsgText = [...]byte{0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x2c,
	0x20, 0x77, 0x6f, 0x72, 0x6c, 0x64, 0x21}

// staticMsg is the contents of the two-byte big-endian size header
// followed by an api.Message protobuf encoding the static message
var staticMsg = [...]byte{
	0x0, 0x1c,
	0xa, 0x4, 0x67, 0x72, 0x61, 0x79, 0x12, 0x5,
	0x6c, 0x79, 0x6e, 0x63, 0x68, 0x1a, 0xd, 0x48,
	0x65, 0x6c, 0x6c, 0x6f, 0x2c, 0x20, 0x77, 0x6f,
	0x72, 0x6c, 0x64, 0x21,
}

// TestReceiveStaticMsg creates a MessageService, then sends the
// static message to it, receives the result, and ensures that it
// matches the original message.  This should test that your
// MessageService implements the same protocol as required by the
// handout.  You could create the analogue to this test (i.e.,
// TestSendStaticMsg), and you probably should.
func TestReceiveStaticMsg(t *testing.T) {
	ms, err := NewMessageService(staticMsgRecipient)
	if err != nil {
		t.Fatalf("Could not create service: %v", err)
	}
	defer ms.Close()

	addr, ok := directory.Lookup(staticMsgRecipient)
	if !ok {
		t.Fatalf("Could not look up service name: %v", err)
	}

	c, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatalf("Could not connect to service: %v", err)
	}
	defer c.Close()

	var wlen int
	for wlen < len(staticMsg) {
		n, err := c.Write(staticMsg[wlen:])
		if err != nil {
			t.Errorf("Failed writing at %d\n", wlen)
			break
		}
		wlen += n
	}

	rmsg := <-ms.Receiver()
	if rmsg.Sender != staticMsgSender ||
		rmsg.Recipient != staticMsgRecipient ||
		bytes.Compare(staticMsgText[:], rmsg.Data) != 0 {
		t.Errorf("Messages differ: %s->%s %#v, %s->%s %#v", staticMsgSender, staticMsgRecipient, staticMsgText[:], rmsg.Sender, rmsg.Recipient, rmsg.Data)
	}
}
