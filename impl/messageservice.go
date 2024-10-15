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
	"cse586.messageservice/api"
	"cse586.messageservice/given/directory"
	"encoding/binary"
	"fmt"
	//proto "github.com/golang/protobuf/proto"
	"google.golang.org/protobuf/proto"
	"net"
)

type messageService struct {
	id string
	// message  api.Message
	listener net.Listener
	receiver chan *api.Message
}

// NewMessageService creates an implementation of the MessageService API,
// according to the behavior of api.NewMessageService.
//
// This method must return an error if id is not known to the
// directory service, or if a listening socket cannot be established
// on the address associated with id.  Otherwise, it should return a
// working MessageService implementation.
func NewMessageService(id string) (api.MessageService, error) {
	addr, ok := directory.Lookup(id)
	if !ok {
		return nil, fmt.Errorf("invalid id: %v", id)
	}

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to listen: %v", err)
	}

	ms := &messageService{
		id:       id,
		listener: listener,
		receiver: make(chan *api.Message),
	}

	go ms.listen()

	return ms, nil
}

func (ms *messageService) Receiver() <-chan *api.Message {
	return ms.receiver
}

func (ms *messageService) Close() error {
	err := ms.listener.Close()
	close(ms.receiver)
	return err
}

func BytesToInt(b []byte) int {
	return int(binary.BigEndian.Uint16(b))
}

func Int16ToBytes(i int16) []byte {
	buf := make([]byte, 2)
	binary.BigEndian.PutUint16(buf, uint16(i))
	return buf
}

func (ms *messageService) listen() {
	for {
		conn, err := ms.listener.Accept()
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Temporary() {
				continue
			}
			if err == net.ErrClosed {
				break
			}
			break
		}

		total_length := 0
		go func() {
			defer conn.Close()
			buf := make([]byte, 0)
			for i := 0; ; i++ {
				tmp := make([]byte, api.MaxMessageLen)
				n, err := conn.Read(tmp)
				if err != nil {
					break
				}

				buf = append(buf, tmp[:n]...)
				if i == 0 {
					total_length = BytesToInt(tmp[:2])
					total_length += 2
					// message too long
					if total_length > api.MaxMessageLen {
						break
					}
				}

				total_length -= n
				if total_length <= 0 {
					break
				}
			}

			binData := buf[2:]
			msg := &api.Message{}
			err = proto.Unmarshal(binData, msg)
			//err = msg.XXX_Unmarshal(binData)
			if err != nil {
				return
			}
			ms.receiver <- msg
		}()
	}
}

func (ms *messageService) Send(recipient string, data []byte) error {
	addr, ok := directory.Lookup(recipient)
	if !ok {
		return fmt.Errorf("unknown recipient ID: %s", recipient)
	}
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to connect: %v", err)
	}
	defer conn.Close()
	msg := &api.Message{
		Sender:    ms.id,
		Recipient: recipient,
		Data:      data,
	}

	var datas []byte
	datas, _ = proto.Marshal(msg)
	//msg.XXX_Marshal(datas, true)
	startBinData := Int16ToBytes(int16(len(datas)))
	datas = append(startBinData, datas[:]...)
	data_length := len(datas)
	if data_length > api.MaxMessageLen {
		return &api.MessageTooLong{
			//Msg: fmt.Sprintf("Message is %d bytes", data_length),
		}
	}

	_, err = conn.Write(datas)
	if err != nil {
		return fmt.Errorf("failed to send: %v", err)
	}

	return nil
}
