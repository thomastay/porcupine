# Porcupine491

Porcupine491 is a fork of [Porcupine](https://github.com/anishathalye/porcupine), 
specialized for EECS 491, the 
[undergraduate Distributed Systems class at the University of Michigan](https://lamport.eecs.umich.edu)

## What is Porcupine491?

Porcupine491 is a fast linearizability checker for testing the correctness of
a distributed Key-Value store. Specifically, given a history like this:

```
T0:  |------------- PUT(200) -------------|
T1:    |- GET(): 200 -|
T2:                        |- GET(): 0 -|
```

We can check the history with Porcupine491 and see that it's not linearizable:

```go
ok := porcupine.CheckEvents(registerModel, events)
// returns false
```

Given another history:

```
T0:  |------------- PUT(200) -------------|
T1:    |- GET(): 200 -|
T2:                        |- GET(): 0 -|
```

We can check the history with Porcupine491 and see that it's not linearizable:

```go
ok := porcupine.CheckEvents(registerModel, events)
// returns false
```

Porcupine491 supports the following three operations:
  1. PUT
  2. APPEND
  3. GET

## Usage

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

TODO: Add how to write go code

TODO: Add how to write a log


## Implementation

Porcupine491 implements the algorithm described in [Faster linearizability
checking via P-compositionality][faster-linearizability-checking], an
optimization of the algorithm described in [Testing for
Linearizability][linearizability-testing].

## License


Copyright (c) 2019 Thomas Tay
Copyright (c) 2017-2018 Anish Athalye. 
Released under the MIT License. See [LICENSE.md][license] for details.

[faster-linearizability-checking]: https://arxiv.org/pdf/1504.00204.pdf
[linearizability-testing]: http://www.cs.ox.ac.uk/people/gavin.lowe/LinearizabiltyTesting/paper.pdf
