package circuits

import (
	"github.com/consensys/gurvy"
	"github.com/philsippl/gnark/frontend"
)

type fromBinaryCircuit struct {
	B0, B1, B2, B3 frontend.Variable
	Y              frontend.Variable `gnark:",public"`
}

func (circuit *fromBinaryCircuit) Define(curveID gurvy.ID, cs *frontend.ConstraintSystem) error {
	cs.AssertIsBoolean(circuit.B0)
	cs.AssertIsBoolean(circuit.B1)
	cs.AssertIsBoolean(circuit.B2)
	cs.AssertIsBoolean(circuit.B3)

	r := cs.FromBinary(circuit.B0, circuit.B1, circuit.B2, circuit.B3)

	cs.AssertIsEqual(circuit.Y, r)
	return nil
}

func init() {
	var circuit, good, bad fromBinaryCircuit
	r1cs, err := frontend.Compile(gurvy.UNKNOWN, &circuit)
	if err != nil {
		panic(err)
	}

	good.B0.Assign(1)
	good.B1.Assign(0)
	good.B2.Assign(1)
	good.B3.Assign(1)
	good.Y.Assign(13)

	bad.B0.Assign(1)
	bad.B1.Assign(0)
	bad.B2.Assign(1)
	bad.B3.Assign(1)
	bad.Y.Assign(12)

	addEntry("frombinary", r1cs, &good, &bad)
}
