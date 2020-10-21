package main

import (
	"github.com/consensys/gurvy"
	"github.com/philsippl/gnark/frontend"
	"github.com/philsippl/gnark/io"
	"github.com/philsippl/gnark/std/hash/mimc"
)

func main() {
	var circuit MiMCCircuit

	// compiles our circuit into a R1CS
	r1cs, err := frontend.Compile(gurvy.BN256, &circuit)
	if err != nil {
		panic(err)
	}

	// save the R1CS to disk
	if err = io.WriteFile("circuit.r1cs", r1cs); err != nil {
		panic(err)
	}

	// good solution
	var witness MiMCCircuit
	witness.PreImage.Assign(35)
	witness.Hash.Assign("19226210204356004706765360050059680583735587569269469539941275797408975356275")
	assignment, _ := frontend.ParseWitness(&witness)

	if err = io.WriteWitness("input.json", assignment); err != nil {
		panic(err)
	}
}

// MiMCCircuit defines a pre-image knowledge proof
// mimc(secret preImage) = public hash
type MiMCCircuit struct {
	// struct tag on a variable is optional
	// default uses variable name and secret visibility.
	PreImage frontend.Variable
	Hash     frontend.Variable `gnark:",public"`
}

// Define declares the circuit's constraints
// Hash = mimc(PreImage)
func (circuit *MiMCCircuit) Define(curveID gurvy.ID, cs *frontend.ConstraintSystem) error {
	// hash function
	mimc, _ := mimc.NewMiMC("seed", curveID)

	// specify constraints
	// mimc(preImage) == hash
	cs.AssertIsEqual(circuit.Hash, mimc.Hash(cs, circuit.PreImage))

	return nil
}
