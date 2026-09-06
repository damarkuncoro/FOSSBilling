package spec

// Specification defines the contract for query filtering criteria
type Specification[T any] interface {
	IsSatisfiedBy(item T) bool
}

// AndSpec combines two specifications with logical AND
type AndSpec[T any] struct {
	left  Specification[T]
	right Specification[T]
}

func And[T any](left, right Specification[T]) Specification[T] {
	return &AndSpec[T]{left: left, right: right}
}

func (s *AndSpec[T]) IsSatisfiedBy(item T) bool {
	return s.left.IsSatisfiedBy(item) && s.right.IsSatisfiedBy(item)
}

// OrSpec combines two specifications with logical OR
type OrSpec[T any] struct {
	left  Specification[T]
	right Specification[T]
}

func Or[T any](left, right Specification[T]) Specification[T] {
	return &OrSpec[T]{left: left, right: right}
}

func (s *OrSpec[T]) IsSatisfiedBy(item T) bool {
	return s.left.IsSatisfiedBy(item) || s.right.IsSatisfiedBy(item)
}

// Filter applies a specification to a slice and returns matching elements
func Filter[T any](items []T, spec Specification[T]) []T {
	result := make([]T, 0)
	for _, item := range items {
		if spec.IsSatisfiedBy(item) {
			result = append(result, item)
		}
	}
	return result
}
