package porcupine

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"testing"
)

type kvInput struct {
	op    uint8 // 0 => get, 1 => put, 2 => append
	key   string
	value string
}

type kvOutput struct {
	value string
}

func getKvModel() Model {
	return Model{
		PartitionEvent: func(history []Event) [][]Event {
			m := make(map[string][]Event)
			match := make(map[uint]string) // id -> key
			for _, v := range history {
				if v.Kind == CallEvent {
					key := v.Value.(kvInput).key
					m[key] = append(m[key], v)
					match[v.Id] = key
				} else {
					key := match[v.Id]
					m[key] = append(m[key], v)
				}
			}
			var ret [][]Event
			for _, v := range m {
				ret = append(ret, v)
			}
			return ret
		},
		Init: func() interface{} {
			// note: we are modeling a single key's value here;
			// we're partitioning by key, so this is okay
			return ""
		},
		Step: func(state, input, output interface{}) (bool, interface{}) {
			inp := input.(kvInput)
			out := output.(kvOutput)
			st := state.(string)
			if inp.op == 0 {
				// get
				return out.value == st, state
			} else if inp.op == 1 {
				// put
				return true, inp.value
			} else {
				// append
				return true, (st + inp.value)
			}
		},
	}
}

func parseKvLog(filename string) []Event {
	file, err := os.Open(filename)
	if err != nil {
		panic("can't open file")
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	invokeGet, _ := regexp.Compile(`{:process (\d+), :type :invoke, :f :get, :key "(.*)", :value nil}`)
	invokePut, _ := regexp.Compile(`{:process (\d+), :type :invoke, :f :put, :key "(.*)", :value "(.*)"}`)
	invokeAppend, _ := regexp.Compile(`{:process (\d+), :type :invoke, :f :append, :key "(.*)", :value "(.*)"}`)
	returnGet, _ := regexp.Compile(`{:process (\d+), :type :ok, :f :get, :key ".*", :value "(.*)"}`)
	returnPut, _ := regexp.Compile(`{:process (\d+), :type :ok, :f :put, :key ".*", :value ".*"}`)
	returnAppend, _ := regexp.Compile(`{:process (\d+), :type :ok, :f :append, :key ".*", :value ".*"}`)

	var events []Event = nil

	id := uint(0)
	procIdMap := make(map[int]uint)
	for {
		lineBytes, isPrefix, err := reader.ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			panic("error while reading file: " + err.Error())
		}
		if isPrefix {
			panic("can't handle isPrefix")
		}
		line := string(lineBytes)

		switch {
		case invokeGet.MatchString(line):
			args := invokeGet.FindStringSubmatch(line)
			proc, _ := strconv.Atoi(args[1])
			events = append(events, Event{CallEvent, kvInput{op: 0, key: args[2]}, id})
			procIdMap[proc] = id
			id++
		case invokePut.MatchString(line):
			args := invokePut.FindStringSubmatch(line)
			proc, _ := strconv.Atoi(args[1])
			events = append(events, Event{CallEvent, kvInput{op: 1, key: args[2], value: args[3]}, id})
			procIdMap[proc] = id
			id++
		case invokeAppend.MatchString(line):
			args := invokeAppend.FindStringSubmatch(line)
			proc, _ := strconv.Atoi(args[1])
			events = append(events, Event{CallEvent, kvInput{op: 2, key: args[2], value: args[3]}, id})
			procIdMap[proc] = id
			id++
		case returnGet.MatchString(line):
			args := returnGet.FindStringSubmatch(line)
			proc, _ := strconv.Atoi(args[1])
			matchId := procIdMap[proc]
			delete(procIdMap, proc)
			events = append(events, Event{ReturnEvent, kvOutput{args[2]}, matchId})
		case returnPut.MatchString(line):
			args := returnPut.FindStringSubmatch(line)
			proc, _ := strconv.Atoi(args[1])
			matchId := procIdMap[proc]
			delete(procIdMap, proc)
			events = append(events, Event{ReturnEvent, kvOutput{}, matchId})
		case returnAppend.MatchString(line):
			args := returnAppend.FindStringSubmatch(line)
			proc, _ := strconv.Atoi(args[1])
			matchId := procIdMap[proc]
			delete(procIdMap, proc)
			events = append(events, Event{ReturnEvent, kvOutput{}, matchId})
		}
	}

	for _, matchId := range procIdMap {
		events = append(events, Event{ReturnEvent, kvOutput{}, matchId})
	}

	return events
}

func checkKv(t *testing.T, logName string, correct bool) {
	t.Parallel()
	kvModel := getKvModel()
	events := parseKvLog(fmt.Sprintf("test_data/kv/%s.txt", logName))
	res := CheckEvents(kvModel, events)
	if res != correct {
		t.Fatalf("expected output %t, got output %t", correct, res)
	}
}

func TestKv1ClientOk(t *testing.T) {
	checkKv(t, "c01-ok", true)
}

func TestKv1ClientBad(t *testing.T) {
	checkKv(t, "c01-bad", false)
}

func TestKv10ClientsOk(t *testing.T) {
	checkKv(t, "c10-ok", true)
}

func TestKv10ClientsBad(t *testing.T) {
	checkKv(t, "c10-bad", false)
}

func TestKv50ClientsOk(t *testing.T) {
	checkKv(t, "c50-ok", true)
}

func TestKv50ClientsBad(t *testing.T) {
	checkKv(t, "c50-bad", false)
}

