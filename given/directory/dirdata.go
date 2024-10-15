/*
Copyright 2021 Ethan Blanton <eblanton@buffalo.edu>

This file is part of a CSE 486/586 project from the University at
Buffalo.  Distribution of this file or its associated repository
requires the written permission of Ethan Blanton.  Sharing this file
may be a violation of academic integrity, please consult the course
policies for more package.
*/

package directory

// Note well that this directory MAY BE REPLACED ENTIRELY (both the
// directory data and the directory implementation!) when testing your
// code.  Do not depend on this data having any particular contents,
// or on any particular hosts being availalble to your code.

// A sample directory with appropriate listen addresses for some
// famous figures in distributed systems.
var directory map[string]*dirEntry = map[string]*dirEntry{
	"gray":    {"gray", "localhost:4586", false},
	"lamport": {"lamport", "localhost:5486", false},
	"lynch":   {"lynch", "localhost:1986", false},
	"mills":   {"mills", "localhost:5905", false},
	"postel":  {"postel", "localhost:1943", false},
}
