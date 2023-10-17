package primitives

import "fmt"

type CircularArray struct {
	array      []interface{}
	head, size int
}

func NewCircularArray(capacity int) *CircularArray {
	return &CircularArray{
		array: make([]interface{}, capacity),
	}
}
func (a *CircularArray) validateItemIndex(index int) error {
	if index < 0 || index >= a.size {
		return fmt.Errorf("out of bounds.")
	}
	return nil
}
func (a *CircularArray) validateCursorIndex(index int) error {
	if index < 0 || index > a.size {
		return fmt.Errorf("out of bounds.")
	}
	return nil
}
func (a *CircularArray) location(index int) int {
	if a.head+index < len(a.array) {
		return a.head + index
	}
	return a.head + index - len(a.array)
}
func (a *CircularArray) enlarge() {
	enlarged := make([]interface{}, 2*len(a.array))
	for i := 0; i < a.size; i++ {
		enlarged[i] = a.array[a.location(i)]
	}
	a.array = enlarged
	a.head = 0
}
func (a *CircularArray) get(index int) (interface{}, error) {
	if err := a.validateItemIndex(index); err != nil {
		return nil, err
	}
	return a.array[a.location(index)], nil
}
func (a *CircularArray) set(index int, item interface{}) error {
	if err := a.validateItemIndex(index); err != nil {
		return err
	}
	a.array[a.location(index)] = item
	return nil
}
func (a *CircularArray) move(from, to, length int) error {
	if from < to {
		for i := length - 1; i >= 0; i-- {
			v, err := a.get(from + i)
			if err != nil {
				return err
			}
			a.set(to+i, v)
		}
	} else if from > to {
		for i := 0; i < length; i++ {
			v, err := a.get(from + i)
			if err != nil {
				return err
			}
			a.set(to+i, v)
		}
	}

	return nil
}
func (a *CircularArray) insert(index, amount int) error {
	if err := a.validateCursorIndex(index); err != nil {
		return err
	}
	if amount < 0 {
		return fmt.Errorf("invalid amount.")
	}

	for a.size+amount > len(a.array) {
		a.enlarge()
	}
	if 2*index >= a.size {
		a.size += amount
		a.move(index, index+amount, a.size-amount-index)
	} else {
		a.head -= amount
		a.size += amount
		if a.head < 0 {
			a.head += len(a.array)
		}
		a.move(amount, 0, index)
	}
	for i := 0; i < amount; i++ {
		a.set(index+i, nil)
	}

	return nil
}
func (a *CircularArray) remove(index, amount int) error {
	if err := a.validateCursorIndex(index); err != nil {
		return err
	}
	if amount < 0 {
		return fmt.Errorf("invalid amount.")
	}
	if err := a.validateCursorIndex(index + amount); err != nil {
		return err
	}
	if 2*index >= a.size-amount {
		a.move(index+amount, index, a.size-amount-index)
		for i := 0; i < amount; i++ {
			a.set(a.size-i-1, nil)
		}
		a.size -= amount
	} else {
		a.move(0, amount, index)
		for i := 0; i < amount; i++ {
			a.set(i, nil)
		}
		a.head += amount
		a.size -= amount
		if a.head >= len(a.array) {
			a.head -= len(a.array)
		}
	}

	return nil
}
