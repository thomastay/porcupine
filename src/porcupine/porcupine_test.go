package porcupine

import (
	"fmt"
	"testing"
)

func checkKv(t *testing.T, logName string, expected bool) {
	// t.Parallel()
	events := ParseKvLog(fmt.Sprintf("../../test_data/%s.json", logName))
	res := CheckKvEvents(events)
	if res != expected {
		t.Fatalf("expected output %t, got output %t", expected, res)
	}
}

func TestTotalOrderViolation(t *testing.T) {
    fmt.Println("Testing for Total Order violation...")
	checkKv(t, "c2_total_order_violation", false)
}

func TestC3Bad(t *testing.T) {
    fmt.Println("Testing for Invalid Read after Write...")
	checkKv(t, "c3_bad", false)
}

func TestC3Good(t *testing.T) {
    fmt.Println("Testing good, short program...")
	checkKv(t, "c3_ok", true)
}

func TestUnreliableLong(t *testing.T) {
    fmt.Println("Testing good, long program 2...")
	checkKv(t, "c3_unreliable_long", true)
}

func TestUnreliableGood(t *testing.T) {
    fmt.Println("Testing good, long program...")
	checkKv(t, "c3_unreliable_ok", true)
}

func TestUnreliableBad(t *testing.T) {
    fmt.Println("Testing for long Invalid Read after Write...")
	checkKv(t, "c3_unreliable_bad", false)
}
