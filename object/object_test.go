package object

import "testing"

func TestBooleanHashKey(t *testing.T) {
	bhka := &Boolean{Value: true}
	bhkb := &Boolean{Value: false}

	bhkc := &Boolean{Value: true}
	bhkd := &Boolean{Value: false}

	if bhka.HashKey() != bhkc.HashKey() {
		t.Errorf("bools with same content have different hash keys")
	}

	if bhkb.HashKey() != bhkd.HashKey() {
		t.Errorf("bools with same content have different hash keys")
	}

	if bhka.HashKey() == bhkd.HashKey() {
		t.Errorf("bools with different content have same hash keys")
	}
}

func TestIntegerHashKey(t *testing.T) {
	ihka := &Integer{Value: 5}
	ihkb := &Integer{Value: 5}

	ihkc := &Integer{Value: 10}
	ihkd := &Integer{Value: 10}

	if ihka.HashKey() != ihkb.HashKey() {
		t.Errorf("integers with same content have different hash keys")
	}

	if ihkc.HashKey() != ihkd.HashKey() {
		t.Errorf("integers with same content have different hash keys")
	}

	if ihka.HashKey() == ihkc.HashKey() {
		t.Errorf("integers with different content have same hash keys")
	}
}

func TestStringHashKey(t *testing.T) {
	shka := &String{Value: "Hello World"}
	shkb := &String{Value: "Hello World"}

	shkc := &String{Value: "My name is"}
	shkd := &String{Value: "My name is"}

	if shka.HashKey() != shkb.HashKey() {
		t.Errorf("strings with same content have different hash keys")
	}

	if shkc.HashKey() != shkd.HashKey() {
		t.Errorf("strings with same content have different hash keys")
	}

	if shka.HashKey() == shkc.HashKey() {
		t.Errorf("strings with different content have same hash keys")
	}
}
