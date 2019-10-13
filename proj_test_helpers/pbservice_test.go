package pbservice

import "porcupine"
import "viewservice"
import "fmt"
import "io"
import "io/ioutil"
import "encoding/json"
import "net"
import "testing"
import "time"
import "log"
import "runtime"
import "math/rand"
import "os"
import "sync"
import "strconv"
import "strings"
import "sync/atomic"

func check(ck *Clerk, key string, value string) {
	v := ck.Get(key)
	if v != value {
		log.Fatalf("Get(%v) -> %v, expected %v", key, v, value)
	}
}

func port(tag string, host int) string {
	s := "/var/tmp/824-"
	s += strconv.Itoa(os.Getuid()) + "/"
	os.Mkdir(s, 0777)
	s += "pb-"
	s += strconv.Itoa(os.Getpid()) + "-"
	s += tag + "-"
	s += strconv.Itoa(host)
	return s
}

func TestThomas(t *testing.T) {
	tag := "basic"
	vshost := port(tag+"v", 1)
	_ = viewservice.StartServer(vshost)
	time.Sleep(time.Second)
	_ = viewservice.MakeClerk("", vshost)
	ck := MakeClerk(vshost, "")
	s1 := StartServer(vshost, port(tag, 1))

	var entries []porcupine.KVLogEntry
	ck.porcupine_put(0, entries, "111", "v1")
	ck.porcupine_get(0, entries, "111")
	ck.porcupine_append(0, entries, "111", "v1")
	ck.porcupine_get(0, entries, "111")

	fileData, _ := json.MarshalIndent(entries, "", "  ")
	_ = ioutil.WriteFile("test.json", fileData, 0644)

	s1.kill()
}

func TestThomas2(t *testing.T) {
	tag := "basic"
	vshost := port(tag+"v", 1)
	_ = viewservice.StartServer(vshost)
	time.Sleep(time.Second)
	_ = viewservice.MakeClerk("", vshost)
	s1 := StartServer(vshost, port(tag, 1))

	results := make(chan []porcupine.KVLogEntry)

	numClients := 3
	for i := 0; i < numClients; i++ {
		go func(id int) {
			ck := MakeClerk(vshost, "")
			var entries []porcupine.KVLogEntry
			entries = ck.porcupine_put(id, entries, "111", "v1")
			entries = ck.porcupine_get(id, entries, "111")
			entries = ck.porcupine_append(id, entries, "111", "v1")
			entries = ck.porcupine_get(id, entries, "111")
			results <- entries
		}(i)
	}

	var entries []porcupine.KVLogEntry
	for i := 0; i < numClients; i++ {
		temp := <-results
		entries = append(entries, temp...)
	}

	fileData, _ := json.MarshalIndent(entries, "", "  ")
	_ = ioutil.WriteFile("test.json", fileData, 0644)

	s1.kill()
}

func TestThomas3(t *testing.T) {
	tag := "csu"
	vshost := port(tag+"v", 1)
	vs := viewservice.StartServer(vshost)
	time.Sleep(time.Second)
	vck := viewservice.MakeClerk("", vshost)

	fmt.Printf("Test: Concurrent Put()s to the same key; unreliable ...\n")

	const nservers = 2
	var sa [nservers]*PBServer
	for i := 0; i < nservers; i++ {
		sa[i] = StartServer(vshost, port(tag, i+1))
		sa[i].setunreliable(true)
	}

	for iters := 0; iters < viewservice.DeadPings*2; iters++ {
		view, _ := vck.Get()
		if view.Primary != "" && view.Backup != "" {
			break
		}
		time.Sleep(viewservice.PingInterval)
	}

	// give p+b time to ack, initialize
	time.Sleep(viewservice.PingInterval * viewservice.DeadPings)

	{
		ck := MakeClerk(vshost, "")
		ck.Put("0", "x")
		ck.Put("1", "x")
	}

	done := int32(0)

	vck.Get()
	const nclients = 3
	const nkeys = 2
	results := make(chan []porcupine.KVLogEntry)

	for xi := 0; xi < nclients; xi++ {
		go func(id int) {
			ck := MakeClerk(vshost, "")
			rr := rand.New(rand.NewSource(int64(os.Getpid() + id)))
			var entries []porcupine.KVLogEntry
			for atomic.LoadInt32(&done) == 0 {
				t := rr.Int() % 3
				k := strconv.Itoa(rr.Int() % nkeys)
				v := strconv.Itoa(rr.Int())
				switch t {
				case 0:
					entries = ck.porcupine_put(id, entries, k, v)
				case 1:
					entries = ck.porcupine_get(id, entries, k)
				case 2:
					entries = ck.porcupine_append(id, entries, k, v)
				}
			}
			results <- entries
		}(xi)
	}
	time.Sleep(5 * time.Second)
	atomic.StoreInt32(&done, 1)

	var entries []porcupine.KVLogEntry
	for i := 0; i < nclients; i++ {
		temp := <-results
		entries = append(entries, temp...)
	}

	fileData, _ := json.MarshalIndent(entries, "", "  ")
	_ = ioutil.WriteFile("../../test_data/c3_unreliable.json", fileData, 0644)

	for i := 0; i < nservers; i++ {
		sa[i].kill()
	}
	time.Sleep(time.Second)
	vs.Kill()
	time.Sleep(time.Second)
}
