package v2

// 임시로 집합 구조체 만들어 놓음
type set[T comparable] struct {
	m map[T]struct{}
}

func newSet[T comparable]() set[T] {
	return set[T]{
		make(map[T]struct{}),
	}
}

func (s *set[T]) Add(k T) {
	s.m[k] = struct{}{}
}

func (s *set[T]) Exists(k T) bool {
	_, ok := s.m[k]
	return ok
}

func (s *set[T]) GetElems() []T {
	var rts []T
	for k := range s.m {
		rts = append(rts, k)
	}
	return rts
}
