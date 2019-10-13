package main

import (
	"fmt"
	"os"
	"porcupine"
)

func main() {
    // Parse the JSON file given as a log
	events := porcupine.ParseKvLog(os.Args[1])
    // Check whether the series of events
    //  is linearizable
    // If it is, it follows the UNIX philosophy
    //  and gives no output
	res := porcupine.CheckKvEntries(events)
	if !res {
		fmt.Println(`
        ##################################
        LINEARIZABILITY VIOLATION DETECTED
        ##################################`)
	}
}
