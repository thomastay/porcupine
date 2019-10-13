package porcupine

// Describes the Key-Value pair model for
// our distributed system
// (Note: Porcupine actually works with arbitrary models)
// See the original repo for more applications

type kvInput struct {
	op    uint8 // 0 => get, 1 => put, 2 => append
	key   string
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
			// uninitialized keys start with ""
			return ""
		},
		Step: func(state, input, output interface{}) (bool, interface{}) {
			inp := input.(kvInput)
			out := output.(kvInput)
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
