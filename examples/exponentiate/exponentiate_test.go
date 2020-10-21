package main

import (
	"testing"

	"github.com/consensys/gurvy"
	"github.com/philsippl/gnark/backend/groth16"
	"github.com/philsippl/gnark/frontend"
)

func TestExponentiate(t *testing.T) {
	assert := groth16.NewAssert(t)

	var expCircuit ExponentiateCircuit
	// compiles our circuit into a R1CS
	r1cs, err := frontend.Compile(gurvy.BN256, &expCircuit)
	if err != nil {
		t.Fatal(err)
	}

	{
		var witness ExponentiateCircuit
		witness.X.Assign(2)
		witness.E.Assign(12)
		witness.Y.Assign(4095)
		assert.ProverFailed(r1cs, &witness) // y != x**e
	}

	{
		var witness ExponentiateCircuit
		witness.X.Assign(2)
		witness.E.Assign(12)
		witness.Y.Assign(4096)
		assert.ProverSucceeded(r1cs, &witness)
	}

}
