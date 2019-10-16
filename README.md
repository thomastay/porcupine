# Porcupine491

Porcupine491 is a fork of [Porcupine](https://github.com/anishathalye/porcupine), 
specialized for EECS 491, the 
[undergraduate Distributed Systems class at the University of Michigan](https://lamport.eecs.umich.edu)

## What is Porcupine491?

Porcupine491 is a fast linearizability checker for testing the correctness of
a distributed Key-Value store. 

Porcupine491 supports a Key-value store with the following operations:
  1. PUT (key, val) --> Updates key with value
  2. APPEND (key, val) --> appends value to key. If key is unassigned, performs PUT(key,val)
  3. GET (key) --> returns the value corresponding to that key. Returns "" if key is unassigned. 

Specifically, given a history like this:

```
C0:  |------------- PUT(0, "200") -------------|
C1:    |- GET(0): "200" -|     |- GET(0): "" -|
```

A little explanation of the above diagram. This is a time series diagram of two clients performing 
concurrent operations. The x-axis is time, and C0 and C1 are two clients. At time 0, client 0 performs 
a PUT operation of "200" (which doesn't return). A little after, client 1 performs a GET 
operation and receives "200". Later on, client 1 performs another GET operation and receives "". Lastly, the server
responds to the PUT request from client 0.

From client 0's perspective, the system is fine; he only did one PUT. From client 1's perspective, the system 
is consistent, since she did only GETs. Together, though, there is no system of operations that could
result in such a system.

We can check the history with Porcupine491 and see that it's not linearizable:

```go
ok := porcupine.CheckKvEntries(entries)
// returns false
```

This system of events would be fine, though:

```
C0:  |------------- PUT(0, "200") -------------|
C1:    |- GET(0): "" -|     |- GET(0): "200" -|
```

Now, client 1 receives a "" first, then a "200". This system is linearizable because client 1's request could 
have been completed before client 0's PUT request. Then, client 1's second GET request could have happened
after the PUT from client 0, i.e. this situation might be how the server served the requests:

```
C0:                      |-- PUT(0, "200") --|
C1:    |- GET(0): "" -|                        |- GET(0): "200" -|
```


```go
ok := porcupine.CheckKvEntries(entries)
// returns true
```

These two tests correspond to the test [c2-bad](test_data/c2_bad.json) and [c2-good](test_data/c2_good.json) in the test repo.

Given another history [total order violation](test_data/c3_total_order_violation.json):

```
C0:  |-- PUT(1) ---|  |-- GET(1) ---|
C1:  |-- PUT(2) ---|  |-- PUT(2) ---|
```
Porcupine rejects this history of events too.


## Usage

### Usage as a Library

Porcupine can be used as a testing library, suitable for integration into _gotest_.
Simply have each client log its events in time order, and then concatenate the logs from each client. You do not have to do any special merging of logs from different clients, just a simple append will do. We provide a file called porcupine\_test\_helpers.go, which aims to be a drop in replacement for the regular PUT, GET and APPEND functions.

Clients must log according to a struct called porcupine.KVLogEntry. See the test helpers for implementation details.

### Usage as a JSON parser

To interface with applications that aren't written in Go, Porcupine491 also can read JSON. Simply write JSON (in the format below) to an output file. Then, use the porcupine491 binaries to read the JSON file as such:

```
./porcupine test.json
```

*Outputs*:

It will produce no output for linearizable JSON files. For nonlinearizable output files, it will print an error message to stdout.

## JSON format

Every time a thread performs an operation, it logs its operation in JSON format. 
For instance, here's what the GET(200) operation above would look like:

```
{
  Time: 2019101200,
  Id: 1,
  Type: start,
  Op: get,
  Key: 10
}
{
  Time: 2019101201,
  Id: 1,
  Type: end,
  Op: get,
  Key: 10,
  Val: 200
}
```

Here's what the entire operation above would look like:

```
[
  {
    Time: 2019101200,
    Id: 0,
    Type: start,
    Op: put,
    Key: 10,
    Val: 200
  },
  {
    Time: 2019101205,
    Id: 0,
    Type: end,
    Op: put,
    Key: 10
  },
  {
    Time: 2019101200,
    Id: 1,
    Type: start,
    Op: get,
    Key: 10
  },
  {
    Time: 2019101201,
    Id: 1,
    Type: end,
    Op: get,
    Key: 10,
    Val: 200
  },
  {
    Time: 2019101203,
    Id: 2,
    Type: start,
    Op: get,
    Key: 10
  },
  {
    Time: 201910120504,
    Id: 2,
    Type: end,
    Op: get,
    Key: 0
  },
]
```

## How it works

Porcupine491 takes an executable model of a system along with a history, and it
runs a decision procedure to determine if the history is linearizable with
respect to the model. Porcupine491 supports specifying history in two ways, either
as a list of operations with given call and return times, or as a list of
call/return events in time order.

Porcupine491 implements the algorithm described in [Faster linearizability
checking via P-compositionality][faster-linearizability-checking], an
optimization of the algorithm described in [Testing for
Linearizability][linearizability-testing].

## License


Copyright (c) 2019 Thomas Tay

Copyright (c) 2017-2018 Anish Athalye. 

Released under the MIT License. See [license](LICENSE.md) for details.

[faster-linearizability-checking]: https://arxiv.org/pdf/1504.00204.pdf
[linearizability-testing]: http://www.cs.ox.ac.uk/people/gavin.lowe/LinearizabiltyTesting/paper.pdf
