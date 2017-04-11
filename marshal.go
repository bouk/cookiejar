package cookiejar

import (
	"encoding/json"
)

func (j *Jar) Each(f func(Entry)) {
	j.mu.Lock()
	defer j.mu.Unlock()

	for _, submap := range j.entries {
		for _, Entry := range submap {
			f(Entry)
		}
	}
}

func (j *Jar) MarshalJSON() ([]byte, error) {
	var entries []Entry

	j.Each(func(e Entry) {
		entries = append(entries, e)
	})

	return json.Marshal(entries)
}

func (j *Jar) Insert(entries ...Entry) {
	j.mu.Lock()
	defer j.mu.Unlock()

	for _, e := range entries {
		host, err := canonicalHost(e.Domain)
		if err != nil {
			continue
		}

		key := jarKey(host, j.psList)
		submap, ok := j.entries[key]
		if !ok {
			submap = make(map[string]Entry)
			j.entries[host] = submap
		}
		id := e.id()
		submap[id] = e
	}

}

func (j *Jar) UnmarshalJSON(data []byte) error {
	var entries []Entry
	if err := json.Unmarshal(data, &entries); err != nil {
		return err
	}

	j.Insert(entries...)

	return nil
}

var (
	_ json.Marshaler   = &Jar{}
	_ json.Unmarshaler = &Jar{}
)
