/*
Copyright 2021 Ethan Blanton <eblanton@buffalo.edu>

This file is part of a CSE 486/586 project from the University at
Buffalo.  Distribution of this file or its associated repository
requires the written permission of Ethan Blanton.  Sharing this file
may be a violation of academic integrity, please consult the course
policies for more package.
*/

package main

import (
	"cse586.messageservice/api"
	"cse586.messageservice/given/detector"
	"cse586.messageservice/impl"
	"fmt"
	"os"
	"sync/atomic"
	"time"
)

// The heartbeat program accepts its own ID and a list of neighbors on
// the command line.  It must create a MessageService on its own ID,
// then wait for the duration given/detector.StartDelay to pass before
// it begins sending (or expecting replies to) any heartbeats.
// Thereafter, it must send a heartbeat message to each of its
// neighbors every given/detector.BeatInterval.  If it fails to
// receive a heartbeat from any host for a duration of
// given/detector.TimeoutDuration, it must print a message "[neighbor]
// failed", where [neighbor] is the ID of the failed neighbor.
//
// The command line arguments are:
// heartbeat id neighbor1 [neighbor2 ...]
//
// If the command is given fewer than 3 total arguments (program name,
// own ID, one neighbor), it should print an error message and exit
// with a nonzero value.

var heartBeatMsgText = [...]byte{0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x2c,
	0x20, 0x77, 0x6f, 0x72, 0x6c, 0x64, 0x21}

// var heartBeatMsgText = [...]byte{}
const (
	maxMonitorNumber = 1024
)

func main() {
	argsNumber := len(os.Args)
	if argsNumber < 3 {
		// fmt.Fprintln(os.Stderr, "lack of parameters, should be used like:heartbeat id neighbor1 [neighbor2 ...]")
		os.Exit(-1)
	}

	var neighbors []string
	var ms api.MessageService
	//var sender string
	var err error
	for i, v := range os.Args {
		if i == 0 {
			continue
		} else if i == 1 {
			//sender = v
			ms, err = impl.NewMessageService(v)
			if err != nil {
				// fmt.Printf("%s failed\n", v)
				os.Exit(-1)
			}

			defer ms.Close()
		} else {
			neighbors = append(neighbors, v)
		}
	}

	// wait for start
	time.Sleep(detector.StartDelay)

	// send every detector.BeatInterval
	go func() {
		for {
			for _, neighbor := range neighbors {
				go func(neighbor string) {
					ms.Send(neighbor, heartBeatMsgText[:])
				}(neighbor)
			}
			time.Sleep(detector.BeatInterval)
		}
	}()

	var lastReceivedTimestamp [maxMonitorNumber]int64
	startTimestamp := int64(time.Now().UnixNano())
	for i := 0; i < maxMonitorNumber; i++ {
		lastReceivedTimestamp[i] = int64(startTimestamp)
	}

	// var lastReceivedTimestampLock sync.Mutex
	// update array lastReceivedTimestamp
	go func() {
		for {
			rmsg := <-ms.Receiver()
			curTimestamp := time.Now().UnixNano()
			for i, v := range neighbors {
				if v == rmsg.Sender {
					atomic.StoreInt64(&lastReceivedTimestamp[i], int64(curTimestamp))
				}
			}
		}
	}()

	// check if timeout
	for {
		curTimestamp := int64(time.Now().UnixNano())
		for i, v := range neighbors {
			gapTimeSec := curTimestamp - atomic.LoadInt64(&lastReceivedTimestamp[i])
			if gapTimeSec >= int64(detector.TimeoutDuration.Nanoseconds()) {
				fmt.Printf("%s failed\n", v)
				atomic.StoreInt64(&lastReceivedTimestamp[i], int64(curTimestamp))
			}
		}
		time.Sleep(detector.TimeoutDuration)
	}
}
