package data

import "testing"

func TestCheckValidation(t *testing.T) {

	prod := &Product{
		Name:  "Coffee",
		Price: 2.59,
		SKU:   "dsik-dsd-trfd",
	}

	err := prod.Validate()

	if err != nil {
		t.Fatal(err)
	}

}
