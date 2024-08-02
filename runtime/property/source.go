package property

import "sync"

type Source interface {
	Name() string
	Source() any
	ContainsProperty(name string) bool
	Property(name string) (any, bool)
	PropertyOrDefault(name string, defaultValue any) any
	PropertyNames() []string
}

type SourceList struct {
	sources []Source
	mu      sync.RWMutex
}

func NewSourceList() *SourceList {
	return &SourceList{
		sources: make([]Source, 0),
		mu:      sync.RWMutex{},
	}
}

func (s *SourceList) Contains(name string) bool {
	defer s.mu.Unlock()
	s.mu.Lock()

	for _, source := range s.sources {
		if source != nil && source.Name() == name {
			return true
		}
	}

	return false
}

func (s *SourceList) Find(name string) (Source, bool) {
	defer s.mu.Unlock()
	s.mu.Lock()

	for _, source := range s.sources {
		if source != nil && source.Name() == name {
			return source, true
		}
	}

	return nil, false
}

func (s *SourceList) AddFirst(source Source) {
	defer s.mu.Unlock()
	s.mu.Lock()

	s.removeIfPresent(source)
	if len(s.sources) == 0 {
		s.sources = append(s.sources, source)
		return
	}

	s.sources = append(s.sources[:1], s.sources[0:]...)
	s.sources[0] = source
}

func (s *SourceList) AddLast(source Source) {
	defer s.mu.Unlock()
	s.mu.Lock()

	s.removeIfPresent(source)
	s.sources = append(s.sources, source)
}

func (s *SourceList) AddAtIndex(index int, source Source) {
	defer s.mu.Unlock()
	s.mu.Lock()

	s.removeIfPresent(source)

	if len(s.sources) == index {
		s.mu.Unlock()
		s.AddLast(source)
		return
	}

	s.sources = append(s.sources[:index+1], s.sources[index:]...)
	s.sources[index] = source
}

func (s *SourceList) Remove(name string) Source {
	source, index := s.findPropertySourceByName(name)

	if index != -1 {
		s.sources = append(s.sources[:index], s.sources[index+1:]...)
	} else {
		return nil
	}

	return source
}

func (s *SourceList) Replace(name string, source Source) {
	_, index := s.findPropertySourceByName(name)

	if index != -1 {
		s.sources[index] = source
	}
}

func (s *SourceList) Count() int {
	return len(s.sources)
}

func (s *SourceList) PrecedenceOf(source Source) int {
	if source == nil {
		return -1
	}

	_, index := s.findPropertySourceByName(source.Name())
	return index
}

func (s *SourceList) ToSlice() []Source {
	sources := make([]Source, len(s.sources))
	copy(sources, s.sources)
	return sources
}

func (s *SourceList) removeIfPresent(source Source) {
	if source == nil {
		return
	}

	_, index := s.findPropertySourceByName(source.Name())

	if index != -1 {
		s.sources = append(s.sources[:index], s.sources[index+1:]...)
	}
}

func (s *SourceList) findPropertySourceByName(name string) (Source, int) {
	for index, source := range s.sources {

		if source.Name() == name {
			return source, index
		}

	}

	return nil, -1
}
