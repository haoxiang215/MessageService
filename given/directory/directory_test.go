package directory

import (
	"testing"
)

var names []string

// Build a name list for all entries in the directory
func init() {
	names = make([]string, 0, len(directory))
	for key := range directory {
		names = append(names, key)
	}
}

// TestRegisterNoID verifies that registering the ID "" results in a
// valid registration.  This depends on the fact that an unregistered
// entry is in the directory.
func TestRegisterNoID(t *testing.T) {
	id, err := Register("")
	if err != nil {
		t.Error(err)
	}
	if id == "" {
		t.Error("Empty ID")
	}
	unregister(id)
}

// TestRegisterSpecificIDs attempts to register every entry in the
// directory, one at a time, and ensures that all are registerable.
// This depends on the fact that every entry in the directory is
// unregistered.
func TestRegisterSpecificIDs(t *testing.T) {
	var names []string
	for _, id := range names {
		val, err := Register(id)
		if err != nil || val != id {
			t.Errorf("Could not register %v", id)
		}
		defer unregister(id)
	}
}

// TestRegisterTwice attempts to register the same ID twice, and
// ensures that this fails.  This depends on the tested ID being
// unregistered before the test runs.
func TestRegisterTwice(t *testing.T) {
	// Check for a test-chosen ID
	val, err := Register(names[0])
	if err != nil || val != names[0] {
		t.Error("Could not register id")
	}
	val, err = Register(names[0])
	if err == nil {
		t.Error("Second registration of registered ID succeeded")
	}
	unregister(names[0])

	// Check for a directory-chosen ID
	id, err := Register("")
	if err != nil {
		t.Error("Could not register id")
	}
	val, err = Register(id)
	if err == nil {
		t.Error("Second registration of registered ID succeeded")
	}
	unregister(id)
}

// TestRegisterUnknownID tries to register an invalid ID and ensures
// that directory.NoSuchID is generated.
func TestRegisterUnknownID(t *testing.T) {
	_, err := Register("asdf")
	if err == nil {
		t.Error("Registration of invalid ID succeeded")
	}
	if _, ok := err.(NoSuchID); !ok {
		t.Error("Registration of invalid ID returned invalid error")
	}
}

// TestAddressLookup ensures that an address lookup succeeds for a known ID.
func TestAddressLookup(t *testing.T) {
	// Use Register to get an ID
	id, err := Register("")
	if err != nil {
		t.Error(err)
	}
	defer unregister(id)
	addr, ok := Lookup(id)
	if !ok || addr == "" {
		t.Error("Bad address")
	}
}
