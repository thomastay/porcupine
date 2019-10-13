package main

import (
    "fmt"
    "os"
    "porcupine"
)

func main() {
	// t.Parallel()
	kvModel := porcupine.GetKvModel()
    events := porcupine.ParseKvLog(os.Args[1])
	res := porcupine.CheckEvents(kvModel, events)
    if !res {
        fmt.Println(`
        ##################################
        LINEARIZABILITY VIOLATION DETECTED
        ##################################`)
    }
}
