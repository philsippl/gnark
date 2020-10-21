package circuits

import (
	"github.com/consensys/gurvy"
	"github.com/philsippl/gnark/frontend"
)

type expCircuit struct {
	X, E frontend.Variable
	Y    frontend.Variable `gnark:",public"`
}

func (circuit *expCircuit) Define(curveID gurvy.ID, cs *frontend.ConstraintSystem) error {
	o := cs.Constant(1)
	b := cs.ToBinary(circuit.E, 4)

	var i int
	for i < len(b) {
		o = cs.Mul(o, o)
		mu := cs.Mul(o, circuit.X)
		o = cs.Select(b[len(b)-1-i], mu, o)
		i++
	}
	cs.AssertIsEqual(circuit.Y, o)
	return nil
}

func init() {
	var circuit, good, bad expCircuit
	r1cs, err := frontend.Compile(gurvy.UNKNOWN, &circuit)
	if err != nil {
		panic(err)
	}

	good.X.Assign(2)
	good.E.Assign(12)
	good.Y.Assign(4096)

	bad.X.Assign(2)
	bad.E.Assign(12)
	bad.Y.Assign(4095)

	addEntry("expo", r1cs, &good, &bad)
}
