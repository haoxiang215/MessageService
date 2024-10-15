/*
Copyright 2021 Ethan Blanton <eblanton@buffalo.edu>

This file is part of a CSE 486/586 project from the University at
Buffalo.  Distribution of this file or its associated repository
requires the written permission of Ethan Blanton.  Sharing this file
may be a violation of academic integrity, please consult the course
policies for more package.
*/

package detector

import "testing"

// TestConfigInfo sanity checks the configuration
func TestConfigInfo(t *testing.T) {
	if TimeoutDuration < BeatInterval*2 {
		t.Error("Timeout duration must be at least 2 * beat interval")
	}
}
