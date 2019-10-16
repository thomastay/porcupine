package porcupine

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"sort"
	"strings"
	"time"
)

const (
	START string = "start"
	END   string = "end"
)

type OpType string

const (
	GET    OpType = "get"
	PUT    OpType = "put"
	APPEND OpType = "append"
)

type KVLogEntry struct {
	Time time.Time
	Id   int
	Type string
	Op   OpType
	Key  string
	Val  string
}

func ParseKvList(entries []KVLogEntry) []Event {
	var events []Event
	id := uint(0)
	procIdMap := make(map[int]uint)
	// the process to ID map maps a process id to the "id" of the
	// operation
	// For instance, a GET request will have a start and end
	// these two requests will have the same id
	// assumes that a process does not issue two requests at a time

	// sort the entries by time
	sort.SliceStable(entries[:], func(lhs, rhs int) bool {
		return entries[lhs].Time.Before(entries[rhs].Time)
	})

	for _, entry := range entries {
		//log.Println(entry)
		var entryType EventKind
		var matchId uint
		if strings.ToLower(entry.Type) == START {
			entryType = CallEvent
			matchId = id
			var ok bool
			if _, ok = procIdMap[entry.Id]; ok {
				log.Panicf("process makes 2 concurrent requests, latest is: %+v", entry)
			}
			procIdMap[entry.Id] = id
			id++
		} else { //END
			var ok bool
			matchId, ok = procIdMap[entry.Id]
			if !ok {
				log.Panicf("unmatched element: %+v", entry)
			}
			delete(procIdMap, entry.Id)
			entryType = ReturnEvent
		}
		var opTypeNum uint8
		switch entry.Op {
		case GET:
			opTypeNum = 0
		case PUT:
			opTypeNum = 1
		case APPEND:
			opTypeNum = 2
		}
		events = append(events,
			Event{entryType,
				kvInput{
					op:    opTypeNum,
					key:   entry.Key,
					value: entry.Val,
				},
				matchId,
			})
	}

	// Add a fake return Event for all the unmapped events
	for _, matchId := range procIdMap {
		events = append(events, Event{ReturnEvent, kvInput{}, matchId})
	}
	//log.Println(events)

	return events
}

func ParseKvLog(filename string) []Event {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		panic("can't open file")
	}

	var entries []KVLogEntry

	err = json.Unmarshal(content, &entries)
	if err != nil {
		log.Panicln(err)
	}

	return ParseKvList(entries)

}

func CheckKvEvents(events []Event) bool {
    return CheckEvents(getKvModel(), events)
}

func CheckKvEntries(entries []KVLogEntry) bool {
    events := ParseKvList(entries)
    return CheckKvEvents(events)
}
