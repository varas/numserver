package repository

import "sync"

// NumberRepository stores unique numbers, extracting in a 2-phase-commit manner to enable transactional support
// * uniqueness is guaranteed against all numbers added
// * ExtractTransaction pulls out only the unique numbers added since the last ExtractTransaction call
type NumberRepository interface {
	AddNumber(number uint32) (unique bool)
	// 2PC extract methods:
	ExtractTransaction() []uint32
	Commit()
	Rollback()
}

// InMemoryRepository stores unique numbers in memory with concurrency support
type InMemoryRepository struct {
	// keeps in memory list of numbers added
	uniques      map[uint32]struct{} // faster access than list
	nonExtracted map[uint32]struct{}
	sync.RWMutex
}

// NewInMemoryRepository stores numbers in memory
func NewInMemoryRepository() NumberRepository {
	return &InMemoryRepository{
		uniques:      make(map[uint32]struct{}),
		nonExtracted: make(map[uint32]struct{}),
	}
}

// AddNumber adds a number if unique returning success
func (r *InMemoryRepository) AddNumber(number uint32) (unique bool) {
	r.RLock()
	_, exists := r.uniques[number]
	if exists {
		r.RUnlock()
		return false
	}
	_, exists = r.nonExtracted[number]
	r.RUnlock()

	if exists {
		return false
	}

	r.Lock()
	r.nonExtracted[number] = struct{}{}
	r.Unlock()

	return true
}

// ExtractTransaction returns unique numbers list delaying data removal to commit
func (r *InMemoryRepository) ExtractTransaction() (uniques []uint32) {
	r.Lock()
	for n := range r.nonExtracted {
		uniques = append(uniques, n)
	}

	return
}

// Commit unlocks emptying the stored numbers
func (r *InMemoryRepository) Commit() {
	// move
	for n := range r.nonExtracted {
		r.uniques[n] = struct{}{}
	}

	r.nonExtracted = make(map[uint32]struct{})
	r.Unlock()
}

// Rollback unlocks without data removal
func (r *InMemoryRepository) Rollback() {
	r.Unlock()
}
