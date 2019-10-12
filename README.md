# Porcupine491

Porcupine491 is a fork of [Porcupine](https://github.com/anishathalye/porcupine), 
specialized for EECS 491, the 
[undergraduate Distributed Systems class at the University of Michigan](https://lamport.eecs.umich.edu)

## What is Porcupine491?

Porcupine491 is a fast linearizability checker for testing the correctness of
distributed systems. Specifically, given a history like this:

Now, suppose we have another history:

```
C0:  |------------- PUT(200) -------------|
C1:    |- GET(): 200 -|
C2:                        |- GET(): 0 -|
```

We can check the history with Porcupine and see that it's not linearizable:

```go
ok := porcupine.CheckEvents(registerModel, events)
// returns false
```

## Usage

Porcupine takes an executable model of a system along with a history, and it
runs a decision procedure to determine if the history is linearizable with
respect to the model. Porcupine supports specifying history in two ways, either
as a list of operations with given call and return times, or as a list of
call/return events in time order.

TODO: Add how to write go code

TODO: Add how to write a log


## Implementation

Porcupine implements the algorithm described in [Faster linearizability
checking via P-compositionality][faster-linearizability-checking], an
optimization of the algorithm described in [Testing for
Linearizability][linearizability-testing].

## License


Copyright (c) 2019 Thomas Tay
Copyright (c) 2017-2018 Anish Athalye. 
Released under the MIT License. See [LICENSE.md][license] for details.

[faster-linearizability-checking]: https://arxiv.org/pdf/1504.00204.pdf
[linearizability-testing]: http://www.cs.ox.ac.uk/people/gavin.lowe/LinearizabiltyTesting/paper.pdf
