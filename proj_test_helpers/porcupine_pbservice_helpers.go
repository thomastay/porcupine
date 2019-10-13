package pbservice

import (
	"porcupine"
	"time"
)

func (ck *Clerk) porcupine_put(id int, entries []porcupine.KVLogEntry, key string, value string) []porcupine.KVLogEntry {
	entries = append(entries, porcupine.KVLogEntry{
		Time: time.Now(),
		Id:   id,
		Type: "start",
		Op:   "put",
		Key:  key,
		Val:  value,
	})
	ck.Put(key, value)
	entries = append(entries, porcupine.KVLogEntry{
		Time: time.Now(),
		Id:   id,
		Type: "end",
		Op:   "put",
		Key:  key,
		Val:  value,
	})
	return entries
}

func (ck *Clerk) porcupine_append(id int, entries []porcupine.KVLogEntry, key string, value string) []porcupine.KVLogEntry {
	entries = append(entries, porcupine.KVLogEntry{
		Time: time.Now(),
		Id:   id,
		Type: "start",
		Op:   "append",
		Key:  key,
		Val:  value,
	})
	ck.Append(key, value)
	entries = append(entries, porcupine.KVLogEntry{
		Time: time.Now(),
		Id:   id,
		Type: "end",
		Op:   "append",
		Key:  key,
		Val:  value,
	})
	return entries
}

func (ck *Clerk) porcupine_get(id int, entries []porcupine.KVLogEntry, key string) []porcupine.KVLogEntry {
	entries = append(entries, porcupine.KVLogEntry{
		Time: time.Now(),
		Id:   id,
		Type: "start",
		Op:   "get",
		Key:  key,
		Val:  "",
	})
	v := ck.Get(key)
	entries = append(entries, porcupine.KVLogEntry{
		Time: time.Now(),
		Id:   id,
		Type: "end",
		Op:   "get",
		Key:  key,
		Val:  v,
	})
	return entries
}
