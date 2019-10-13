package porcupine

import (
    "fmt"
	"testing"
)

func checkKv(t *testing.T, logName string, correct bool) {
	// t.Parallel()
	kvModel := GetKvModel()
	events := ParseKvLog(fmt.Sprintf("../../test_data/%s.txt", logName))
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
