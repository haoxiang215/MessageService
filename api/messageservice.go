/*
Copyright 2021 Ethan Blanton <eblanton@buffalo.edu>

This file is part of a CSE 486/586 project from the University at
Buffalo.  Distribution of this file or its associated repository
requires the written permission of Ethan Blanton.  Sharing this file
may be a violation of academic integrity, please consult the course
policies for more package.
*/

package api

// No marshalled message should be larger than MaxMessageLen
const MaxMessageLen = 65535

type MessageTooLong struct {
	Msg string
}

// MessageService provides a messaging service that can send and
// receive api.Message messages.  Each message has a sender, a
// receiver, and a byte buffer containing the message data.
type MessageService interface {
	// Receiver returns a channel that receives messages from this
	// MessageService.  The message struct includes the sender,
	// receiver (which should be this MessageService), and
	// contents of the message.
	Receiver() <-chan *Message

	// Send will create a Message from this MessageService to the
	// named recipient and send it, provided that the recipient
	// can be found.  If the marshalled message would be larger
	// than MaxMessageLen, it should return a MessageTooLong
	// error.  If any other error occurs, it should return an
	// appropriate error.
	Send(recipient string, data []byte) error

	// Close closes this MessageService, releasing the listening
	// socket and closing all incoming and outgoing sockets.
	Close() error
}

// Implementing this function makes MessageTooLong an error type that
// can be returned.  Create and return an error of this type with
// something like:
//
// return &MessageTooLong{Msg: fmt.Sprintf("Message is %d bytes", len)}
func (err *MessageTooLong) Error() string {
	return err.Msg
}
