/*
Copyright 2021 Ethan Blanton <eblanton@buffalo.edu>

This file is part of a CSE 486/586 project from the University at
Buffalo.  Distribution of this file or its associated repository
requires the written permission of Ethan Blanton.  Sharing this file
may be a violation of academic integrity, please consult the course
policies for more package.
*/

// The directory package provides a directory service for mapping
// logical process names to socket addresses.  It exposes two public
// functions, Register and Lookup.  When a MessageService
// impelmentation starts up, it should call the Register function to
// register itself with the directory service.  If it does not know
// its own name, the directory service can assign it one at that time.
// It can also look up its own listen address via the Lookup function
// using its own name, or the listening address of another process by
// looking up that process's name.
//
// Note that this directory service is designed to allow either
// multiple MessageService instances within a single Go process OR
// multiple MessageService instances in multiple Go processes.  Note
// in particular that while the Register function will prevent
// multiple MessageServices in the same Go process from registering
// the same name, it CANNOT prevent this for multiple Go processes.
//
// You must make NO ASSUMPTIONS about the names that can be queried
// using this directory service or the addresses that those queries
// return.  This implementation may be switched out for a completely
// different implementation, or a different database, during grading.
//
// The sample directory provided to you can be found in dirdata.go.
package directory

import (
	"errors"
	"fmt"
)

// NoSuchID is the error returned when a directory registration or
// lookup attempts to reference an invalid ID.
type NoSuchID string

// Register with the directory service.  If id is the empty string,
// the directory service will return a suitable ID.  If id is
// non-empty, attempt to register with that ID.
//
// If the registered ID is already in use, or if there are no
// available unregistered addresses, this returns an error.
func Register(id string) (string, error) {
	// This channel will be used to receive the result of the
	// registration request from the directory service goroutine.
	c := make(chan dirResult)
	requests <- &dirRequest{id, dir_REGISTER, c}

	// The directory server will eventually service the request
	// sent above.  When it's done, it will send us back a
	// dirResult, which we can send (almost) directly back to the
	// caller.  We have to check for err first, because if entry
	// is nil the following return would crash.
	result := <-c
	if result.err != nil {
		return "", result.err
	}
	return result.entry.id, result.err
}

// Lookup searches the directory for the given ID and returns its
// address if known, or !ok if it is unknown.  A known address will be
// returned whether or not it has been registered.
func Lookup(id string) (string, bool) {
	c := make(chan dirResult)
	requests <- &dirRequest{id, dir_LOOKUP, c}
	result := <-c
	if result.entry == nil {
		return "", false
	}
	return result.entry.address, true
}

func (err NoSuchID) Error() string {
	return fmt.Sprintf("Unknown ID '%s'", string(err))
}

// dirEntry represents a directory entry.
//
// id and address are immutable fields, inUse is mutable.
type dirEntry struct {
	id      string // id is the global ID of this entry
	address string // address is the listen address for this ID
	inUse   bool   // inUse is true if this ID has been registered
}

// dirAction represents an action that can be taken by the directory service.
type dirAction int

const (
	// dir_REGISTER requests registration with the directory service.
	dir_REGISTER dirAction = iota
	// dir_UNREGISTER is used only internally for testing, and
	// requests that an ID is unregistered.
	dir_UNREGISTER
	// dir_LOOKUP requests a lookup on an ID.
	dir_LOOKUP
)

// dirRequest is a request to the directory service.
type dirRequest struct {
	// id is the process ID on which the action should be
	// preformed.
	id string
	// action is the action being requested.
	action dirAction
	// c is where the result of the request will be returned by
	// the directory service goroutine.
	c chan<- dirResult
}

// dirResult is a result returned by the directory service.
type dirResult struct {
	// entry is the directory entry satisfying the corresponding
	// request.  If the result is an error, it may be nil.  The
	// receiving function most not edit the entry or query its
	// mutable fields.
	entry *dirEntry
	// err is non-nil if the request could not be satisfied.
	err error
}

// requests is the gateway between Register and Lookup (and
// unregister) and the directory service.
var requests chan *dirRequest

// Unregister from the directory service.  This does no error checking
// and is used only for testing.
func unregister(id string) {
	requests <- &dirRequest{id, dir_UNREGISTER, nil}
}

// dirService is the directory service goroutine.  This listens on the
// requests channel for incoming dirRequest structs, attempts to
// service the request, and writes response on the request's response
// channel (created in the public directory function that creates the
// request).
//
// This is a goroutine that listens on a channel and sends its
// responses on channels so that it can be the single execution
// context that manipulates the map named directory.  The structs in
// the directory map are not safe for concurrent editing, so putting
// all of the map manipulations in a single goroutine ensures that the
// directory is always consistent.  In essence, dirService "owns" the
// directory map and its contents.
func dirService() {
	// ranging over a channel will "drain" the channel; that is,
	// iterate over every message sent on the channel until the
	// channel is closed.
	for req := range requests {
		// Every path through this switch needs the entry that
		// has been queried, so we fetch it once at the top.
		entry, found := directory[req.id]
		switch req.action {
		case dir_LOOKUP:
			// Look up a given ID in the directory, and
			// return it.
			//
			// A response can be sent back to the
			// requester by writing to the channel
			// included in the dirRequest struct received
			// over the request channel.
			req.c <- dirResult{entry, nil}
		case dir_REGISTER:
			// Register a name if it exists in the map.
			// If the requested name is "", choose and
			// return an unregistered name.
			if !found && req.id == "" {
				for _, entry = range directory {
					if !entry.inUse {
						break
					}
				}
				if entry.inUse {
					req.c <- dirResult{nil, errors.New("No available IDs")}
				}
			}
			if entry == nil {
				// This is the only defined error on
				// the directory service; if a process
				// tries to register an ID that
				// doesn't exist, we return NoSuchID.
				req.c <- dirResult{nil, NoSuchID(req.id)}
				continue
			}
			if entry.inUse {
				req.c <- dirResult{nil, errors.New("Already registered")}
			} else {
				entry.inUse = true
				req.c <- dirResult{entry, nil}
			}
		case dir_UNREGISTER:
			// This action is used only in testing.  It
			// unregisters a previously-registered ID.
			//
			// This could be exposed, but the potential
			// for error is large.
			if found {
				entry.inUse = false
			}
		}
	}
}

// All init functions are called before the package is used.  This one
// createsthe requests channel and starts the directory service
// goroutine.
func init() {
	requests = make(chan *dirRequest)
	go dirService()
}
