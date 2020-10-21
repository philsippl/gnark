/*
Copyright © 2020 ConsenSys

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package fields

import (
	"math/big"

	"github.com/consensys/gurvy/bls377"
	"github.com/consensys/gurvy/bls377/fp"
	"github.com/consensys/gurvy/bw761/fr"
	"github.com/philsippl/gnark/backend"
	"github.com/philsippl/gnark/frontend"
)

// E2 element in a quadratic extension
type E2 struct {
	A0, A1 frontend.Variable
}

// Neg negates a e2 elmt
func (e *E2) Neg(cs *frontend.ConstraintSystem, e1 *E2) *E2 {
	e.A0 = cs.Sub(0, e1.A0)
	e.A1 = cs.Sub(0, e1.A1)
	return e
}

// Add e2 elmts
func (e *E2) Add(cs *frontend.ConstraintSystem, e1, e2 *E2) *E2 {
	e.A0 = cs.Add(e1.A0, e2.A0)
	e.A1 = cs.Add(e1.A1, e2.A1)
	return e
}

// Sub e2 elmts
func (e *E2) Sub(cs *frontend.ConstraintSystem, e1, e2 *E2) *E2 {
	e.A0 = cs.Sub(e1.A0, e2.A0)
	e.A1 = cs.Sub(e1.A1, e2.A1)
	return e
}

// Mul e2 elmts: 5C
func (e *E2) Mul(cs *frontend.ConstraintSystem, e1, e2 *E2, ext Extension) *E2 {

	one := big.NewInt(1)
	minusOne := big.NewInt(-1)

	// 1C
	l1 := cs.LinearExpression(
		cs.Term(e1.A0, one),
		cs.Term(e1.A1, one),
	)
	l2 := cs.LinearExpression(
		cs.Term(e2.A0, one),
		cs.Term(e2.A1, one),
	)
	u := cs.Mul(l1, l2)

	// 2C
	ac := cs.Mul(e1.A0, e2.A0)
	bd := cs.Mul(e1.A1, e2.A1)

	// 1C
	l3 := cs.LinearExpression(
		cs.Term(u, one),
		cs.Term(ac, minusOne),
		cs.Term(bd, minusOne),
	)
	e.A1 = cs.Mul(l3, 1)

	// 1C
	buSquare := backend.FromInterface(ext.uSquare)
	l4 := cs.LinearExpression(
		cs.Term(ac, one),
		cs.Term(bd, &buSquare),
	)
	e.A0 = cs.Mul(l4, 1)

	return e
}

// MulByFp multiplies an fp2 elmt by an fp elmt
func (e *E2) MulByFp(cs *frontend.ConstraintSystem, e1 *E2, c interface{}) *E2 {
	e.A0 = cs.Mul(e1.A0, c)
	e.A1 = cs.Mul(e1.A1, c)
	return e
}

// MulByIm multiplies an fp2 elmt by the imaginary elmt
// ext.uSquare is the square of the imaginary root
func (e *E2) MulByIm(cs *frontend.ConstraintSystem, e1 *E2, ext Extension) *E2 {
	x := e1.A0
	e.A0 = cs.Mul(e1.A1, ext.uSquare)
	e.A1 = x
	return e
}

// Conjugate conjugation of an e2 elmt
func (e *E2) Conjugate(cs *frontend.ConstraintSystem, e1 *E2) *E2 {
	e.A0 = e1.A0
	e.A1 = cs.Sub(0, e1.A1)
	return e
}

// Inverse inverses an fp2elmt
func (e *E2) Inverse(cs *frontend.ConstraintSystem, e1 *E2, ext Extension) *E2 {

	var a0, a1, t0, t1, t1beta frontend.Variable

	a0 = e1.A0
	a1 = e1.A1

	t0 = cs.Mul(e1.A0, e1.A0)
	t1 = cs.Mul(e1.A1, e1.A1)

	t1beta = cs.Mul(t1, ext.uSquare)
	t0 = cs.Sub(t0, t1beta)
	t1 = cs.Inverse(t0)
	e.A0 = cs.Mul(a0, t1)
	e.A1 = cs.Sub(0, a1)
	e.A1 = cs.Mul(e.A1, t1)

	return e
}

// Assign a value to self (witness assignment)
func (e *E2) Assign(a *bls377.E2) {
	e.A0.Assign(bls377FpTobw761fr(&a.A0))
	e.A1.Assign(bls377FpTobw761fr(&a.A1))
}

// MustBeEqual constraint self to be equal to other into the given constraint system
func (e *E2) MustBeEqual(cs *frontend.ConstraintSystem, other E2) {
	cs.AssertIsEqual(e.A0, other.A0)
	cs.AssertIsEqual(e.A1, other.A1)
}

func bls377FpTobw761fr(a *fp.Element) (r fr.Element) {
	for i, v := range a {
		r[i] = v
	}
	return
}
