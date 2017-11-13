package repository

import (
	"testing"

	"sync"

	"github.com/stretchr/testify/assert"
)

func TestInMemoryRepository_AddNumber(t *testing.T) {
	uniqueNumbers := []uint32{11, 22, 33, 44, 55}

	r := NewInMemoryRepository()

	for _, n := range uniqueNumbers {
		unique := r.AddNumber(n)
		assert.True(t, unique)
	}

	assert.False(t, r.AddNumber(11), "repeated number should not return unique")
}

func TestInMemoryRepository_ExtractTransaction(t *testing.T) {
	uniqueNumbers := []uint32{11, 22, 33, 44, 55}
	repeatedNumbers := []uint32{11, 22}

	r := NewInMemoryRepository()

	for _, n := range append(uniqueNumbers, repeatedNumbers...) {
		_ = r.AddNumber(n)
	}

	result := r.ExtractTransaction()
	r.Commit()

	assert.Len(t, result, len(uniqueNumbers))
	for _, n := range uniqueNumbers {
		assert.Contains(t, result, n)
	}

	result2 := r.ExtractTransaction()
	r.Commit()

	assert.Len(t, result2, 0)
}

func TestInMemoryRepository_Rollback(t *testing.T) {
	r := NewInMemoryRepository()

	_ = r.AddNumber(11)

	result := r.ExtractTransaction()
	r.Rollback()

	assert.Len(t, result, 1)

	result2 := r.ExtractTransaction()
	r.Commit()

	assert.Len(t, result2, 1)
}

func TestInMemoryRepository_ExtractTransaction_ExtractsUniquesSinceLastCall(t *testing.T) {
	set1 := []uint32{11, 22, 33, 44, 55}
	set2 := []uint32{11, 66}

	r := NewInMemoryRepository()

	for _, n := range set1 {
		_ = r.AddNumber(n)
	}

	result := r.ExtractTransaction()
	r.Commit()

	for _, n := range set2 {
		_ = r.AddNumber(n)
	}

	result2 := r.ExtractTransaction()
	r.Commit()

	assert.Len(t, result, len(set1))
	assert.Len(t, result2, 1, "intersection between set1 and set2")
}

func TestInMemoryRepository_AddNumberSupportsConcurrency(t *testing.T) {
	numbersA := makeRange(1, 200)
	numbersB := makeRange(100, 200)
	numbersC := makeRange(200, 300)

	repo := NewInMemoryRepository()

	// wait for concurrent adders
	wg := sync.WaitGroup{}
	wg.Add(3)

	// block race start
	ready := make(chan struct{})

	go addNumbers(repo, numbersA, ready, &wg)
	go addNumbers(repo, numbersB, ready, &wg)
	go addNumbers(repo, numbersC, ready, &wg)

	close(ready)
	wg.Wait()

	assert.Len(t, repo.ExtractTransaction(), 300)
}

func addNumbers(r NumberRepository, numbersA []uint32, ready chan struct{}, wg *sync.WaitGroup) {
	<-ready
	for _, n := range numbersA {
		_ = r.AddNumber(n)
	}
	wg.Done()
}

func makeRange(min, max uint32) (list []uint32) {
	for i := min; i <= max; i++ {
		list = append(list, i)
	}
	return
}
