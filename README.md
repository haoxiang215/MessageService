Message Service and Failure Detector
===

This assignment requires you to implement an N-way message service and
an all-pairs heartbeat failure detector.
You should have received a detailed handout with requirements and
grading criteria.

Repository Layout
---

The root of this repository contains this README and a Makefile.  The
Makefile can be used to build commands for testing your project or
satisfying the project requirements.  In particular, it takes care of
building the Go implementation of the protocol buffer definition
required for communication between message service endpoints.

Your entire `MessageService` implementation should be in or under the
`impl` package in the `impl` directory.  One file,
`messageservice.go`, is already in that directory to provide a stub
implementation of `impl.NewMessageService()`.  Your heartbeat detector
should be in `cmd/heartbeat/`.

The other directories of this repository, `api` and `given`, are not
submitted to Autograder and _will not be used when evaluating your
project._  You must not change any files in these directories, and
expect those changes to be available when your code is evaluated!

Given Code
---

The `api` directory defines the `MessageService` API, which is
documented in comments and defined in the project handout.  It also
contains the definition of the ProtoBuf message that will be sent
between the processes in your MessageService.

The `given` directory defines some constants used in your heartbeat
detector under the `detector` package (which may be changed for
grading, so you should **NOT** hard-code their values into your
code!), as well as a directory service interface under the `directory`
package.

The directory service interface provides a mapping from process names
to listening addresses, and comes pre-populated with five famous
figures in distributed systems.  You may add more process ID and
address pairs to the mapping if you wish to test with more than five
processes.  The only public entry points to the directory service API
are `Register()` and `Lookup()`.  Except for testing purposes, your
code should call `Register()` only for the ID passed to
`impl.NewMessageService()`, and should call `Lookup()` only on that ID
or IDs passed to `MessageService.Send()`.

The Makefile will build the protobuf implementation and your
`heartbeat` executable when you run `make`, and run the given test
when you run `make test`.  You may wish to add other commands or tests
to the Makefile, and that is fine; it will not be used for building
your code on the Autograder.

The given code is extensively commented, and you should read and try
to understand all of it.  In particular, the directory service
contains a solution for avoiding races through message passing that
will be very useful to understand going forward.

Testing
---

One test is included for you.  It contains a constant byte array
representing a valid message as sent over the socket between two
`MessageService` processes.  If `make test` does not compile and pass,
you will almost certainly receive no points.

_You should write your own tests._  There are no tests for the
heartbeat command, and only one simple test for the `MessageService`
implementation included in this repository.

Submission
---

Run `make submission` to build `messageservice.tar`, which you will
submit to Autograder.  If you build commands for testing and do not
add them to the Makefile, they may be included as part of the tarball,
which will cause submission to fail.  Add them to the Makefile to
exclude them automatically, or make sure to remove their compiled
binaries before running `make submission`.
