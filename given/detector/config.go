/*
Copyright 2021 Ethan Blanton <eblanton@buffalo.edu>

This file is part of a CSE 486/586 project from the University at
Buffalo.  Distribution of this file or its associated repository
requires the written permission of Ethan Blanton.  Sharing this file
may be a violation of academic integrity, please consult the course
policies for more package.
*/

// The detector package contains configuration information for the
// heartbeat failure detector.
package detector

import "time"

// The heartbeat program should wait startDelay seconds after creating
// its MessageService and before attempting to send any heartbeats.
var StartDelay = 5 * time.Second

// The heartbeat program should assume that a neighbor has failed if
// heartbeats have gone unacknowledged for timeoutDuration.
var TimeoutDuration = 3 * time.Second

// The heartbeat program should send heartbeats every beatInterval.
var BeatInterval = 1 * time.Second
