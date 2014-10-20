package main

type State struct {
	pools Pools
}

var state = State{Pools{}}

func (s *State) GetPool(id string) *Pool {
	p, ok := s.pools[id]
	if ok {
		return p
	}
	return nil
}

func (s *State) SetPool(id string, p *Pool) {
	if p.Apps == nil {
		p.Apps = &Apps{}
	}
	s.pools[id] = p
}

func (s *State) GetPoolIDs() []string {
	ids := make([]string, len(s.pools))
	i := 0
	for id, _ := range s.pools {
		ids[i] = id
		i++
	}
	return ids
}
