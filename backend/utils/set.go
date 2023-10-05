package utils

type Set[D comparable] map[D]bool

func NewSet[D comparable](items ...D) Set[D] {
	s := make(map[D]bool)
	for _, item := range items {
		s[item] = true
	}
	return s
}

func (s Set[D]) Add(item D) {
	s[item] = true
}

func (s Set[D]) Remove(item D) {
	delete(s, item)
}

func (s Set[D]) Contains(item D) bool {
	_, ok := s[item]
	return ok
}

func (s Set[D]) Size() int {
	return len(s)
}

func (s Set[D]) Items() []D {
	items := make([]D, 0, len(s))
	for item := range s {
		items = append(items, item)
	}
	return items
}

func (s Set[D]) Union(other Set[D]) Set[D] {
	result := NewSet[D]()
	for item := range s {
		result.Add(item)
	}
	for item := range other {
		result.Add(item)
	}
	return result
}

func (s Set[D]) Intersection(other Set[D]) Set[D] {
	result := NewSet[D]()
	for item := range s {
		if other.Contains(item) {
			result.Add(item)
		}
	}
	return result
}
