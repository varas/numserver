# NumServer

## Specs

Using any programming language (taking performance into consideration), write a server ("Application") that opens a socket and restricts input to at most 5 concurrent clients.
Clients will connect to the Application and write any number of 9 digit numbers, and then close the connection.
The Application must write a de-duplicated list of these numbers to a log file in no particular order.

Primary Considerations

 - The Application should work correctly as defined below in Requirements.
 - The overall structure of the Application should be simple.
 - The code of the Application should be descriptive and easy to read, and the build method and runtime parameters must be well-described and work.
 - The design should be resilient with regard to data loss.
 - The Application should be optimized for maximum throughput, weighed along with the other Primary Considerations and the Requirements below.
 - The solution should be able to be build and run from the command line.

Include specific instructions on dependencies, build, test, and run instructions.

Requirements

1. The Application must accept input from at most 5 concurrent clients on TCP/IP port 4000.
2. Input lines presented to the Application via its socket must either be composed of exactly nine decimal digits (e.g.: 314159265 or 007007009) immediately followed by a server-native newline sequence; or a termination sequence as detailed below.
3. Numbers presented to the Application must include leading zeros as necessary to ensure they are each 9 decimal digits.
4. The log file, to be named "numbers.log", must be created anew and/or cleared when the Application starts.
5. Only numbers may be written to the log file. Each number must be followed by a server-native newline sequence.
6. No duplicate numbers may be written to the log file.
7. Any data that does not conform to a valid line of input should be discarded and the client connection terminated immediately and without comment.
8. Every 10 seconds, the Application must print a report to standard output:
   * The difference since the last report of the count of new unique numbers that have been received.
   * The difference since the last report of the count of new duplicate numbers that have been received.
   * The total number of unique numbers received for this run of the Application.
   * Example text: Received 50 unique numbers, 2 duplicates. Unique total: 567231
9. If any connected client writes a single line with only the word "terminate" followed by a server-native newline sequence, the Application must disconnect all clients and perform a clean shutdown as quickly as possible.
10. Clearly state all of the assumptions you made in completing the Application.

Notes

 - Ensure your application is executable from the command line.
 - Distribute your code with all the necessary instructions to build it and run it.
 - You may write tests at your own discretion. Tests are useful to ensure your Application passes Primary Consideration A.
 - You are not restricted to the Go or Java language and libraries and frameworks that are considered part of the language. You may use common libraries and frameworks such as the Apache Commons and Google Guava, particularly if their use helps improve Application simplicity and readability.
 - Your Application may not for any part of its operation use or require the use of external systems, for example Apache Kafka or Redis.
 - At your discretion, leading zeroes present in the input may be stripped—or not used—when writing output to the log or console.
 - Robust implementations of the Application typically handle more than 2M numbers per 10-second reporting period on a modern MacBook Pro laptop (e.g.: 16 GiB of RAM and a 2.5 GHz Intel i7 processor).

---

## Decisions taken on this implementation

### On product

- Errors are printed to stderr.
- There are no leading zeroes on the output. As it stores numbers this avoid extra load.
- Log file is flushed on intervals, if log file write fails these numbers will be retried on the next flush interval.
- Numbers handled are supossed to fit in memory, otherwise a disk-fetch policy should be added to check for number uniqueness. An approximate-membership-query approach like bloom-filters would fit here to reduce memory consumption and to avoid disk access.

### On design

- Server package wraps the whole *numserver*. 
- There are no end-to-end tests, as there is no CI env and `make test-stress` acts as nice acceptance test to validate output and evaluate performance.
- The runtime acts as service wiring (kind of service locator pattern) and life-cycle management.
- State (numbers and report counts) are managed in a transactional way, so when written if write fails we don’t remove them from the in-memory storage (avoid data loss).
- State is guaranteed via mutex as usual. For top-notch state management a CRDT approach could used to increase throughput if the amount of supported clients increase on the specs.
- Some communications, like the reporter stats could be also achieved via go-channel. This will decouple some services, but a dependency-injection this seemed clear and straight forward to me.

## Instructions

### Build

`make build`: builds for localhost

`make build-linux`: builds for real servers

### Run

`bin/numserver`

Optional arguments: `-port PORT` and `-file LOG` can be used, cli help is provided on incorrect arguments usage.

> Client is not provided as plain netcat can be used `nc localhost 4000`

### Test

`make test`: Unit & integration tests

`make test-stress`: Stress test, used as acceptance test

