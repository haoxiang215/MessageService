GO := go
PROTOC := protoc

# Build the following commands.  This assumes that each command
# CMDNAME is in the directory cmd/CMDNAME, and can be built by
# changing to that directory and running go build.
COMMANDS := heartbeat

# This rule turns COMMANDS into executable filenames, do not change.
# You don't need to understand this.
CMDFILES := $(shell for word in $(COMMANDS); do echo cmd/$$word/$$word; done)

# Running the command `make` with no arguments should build the
# commands specified in $(COMMANDS).  This is for your convenience,
# you can also build them with `go build`.
#
# The body of this rule is a shell script that loops through every
# command defined in COMMANDS and builds it with go build.  Shell
# scripts embedded in Makefiles have somewhat strange parsing rules
# due to the way that Make works; see `info make` for more
# information.
all: api/message.pb.go go.sum
	for cmd in $(COMMANDS); do (cd cmd/$$cmd; go build); done

# Build a submission tarball.
submission:
	tar cf messageservice.tar \
	    $(shell for word in $(CMDFILES); do echo "--exclude $$word"; done) \
	    --exclude '._*' --exclude '.DS_Store' \
	    Makefile impl cmd

# Run the tests in impl.  If you have other tests you wish to run, you
# may run them here, as well.
test: all
	go test cse586.messageservice/impl

# Run tests on the given code.  You should not need to do this, and
# the outcome of these tests should not affect your grade (assuming
# you have not edited the code that is not submitted to Autograder).
giventest:
	go test cse586.messageservice/given/directory
	go test cse586.messageservice/given/detector

go.sum: api/message.pb.go
	go get cse586.messageservice/api

# This will clean things up a bit.  You can remove other files or
# perform other actions here as you like.
clean:
	rm -f $(CMDFILES) messageservice.tar api/message.pb.go

# Build a protobuf implementation from a protocol description
%.pb.go: %.proto
	$(PROTOC) --go_out=. --go_opt=paths=source_relative $<

.PHONY: all clean submission test
