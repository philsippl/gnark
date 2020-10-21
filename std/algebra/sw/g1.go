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

package sw

import (
	"math/big"

	"github.com/consensys/gurvy/bls377"
	"github.com/consensys/gurvy/bls377/fp"
	"github.com/consensys/gurvy/bw761/fr"
	"github.com/philsippl/gnark/frontend"
)

// G1Jac point in Jacobian coords
type G1Jac struct {
	X, Y, Z frontend.Variable
}

// G1Affine point in affine coords
type G1Affine struct {
	X, Y frontend.Variable
}

// Neg outputs -p
func (p *G1Jac) Neg(cs *frontend.ConstraintSystem, p1 *G1Jac) *G1Jac {
	p.X = p1.X
	p.Y = cs.Sub(0, p1.Y)
	p.Z = p1.Z
	return p
}

// Neg outputs -p
func (p *G1Affine) Neg(cs *frontend.ConstraintSystem, p1 *G1Affine) *G1Affine {
	p.X = p1.X
	p.Y = cs.Sub(0, p1.Y)
	return p
}

// AddAssign adds p1 to p using the affine formulas with division, and return p
func (p *G1Affine) AddAssign(cs *frontend.ConstraintSystem, p1 *G1Affine) *G1Affine {

	// compute lambda = (p1.y-p.y)/(p1.x-p.x)

	one := big.NewInt(1)
	minusOne := big.NewInt(-1)

	l1 := cs.LinearExpression(
		cs.Term(p1.Y, one),
		cs.Term(p.Y, minusOne),
	)
	l2 := cs.LinearExpression(
		cs.Term(p1.X, one),
		cs.Term(p.X, minusOne),
	)
	l := cs.Div(l1, l2)

	// xr = lambda**2-p.x-p1.x
	_x := cs.LinearExpression(
		cs.Term(cs.Mul(l, l), one),
		cs.Term(p.X, minusOne),
		cs.Term(p1.X, minusOne),
	)

	// p.y = lambda(p.x-xr) - p.y
	t1 := cs.Mul(p.X, l)
	t2 := cs.Mul(l, _x)
	l3 := cs.LinearExpression(
		cs.Term(t1, one),
		cs.Term(t2, minusOne),
		cs.Term(p.Y, minusOne),
	)
	p.Y = cs.Mul(l3, 1)

	//p.x = xr
	p.X = cs.Mul(_x, 1)
	return p
}

// AssignToRefactor sets p to p1 and return it
func (p *G1Jac) AssignToRefactor(cs *frontend.ConstraintSystem, p1 *G1Jac) *G1Jac {
	p.X = cs.Constant(p1.X)
	p.Y = cs.Constant(p1.Y)
	p.Z = cs.Constant(p1.Z)
	return p
}

// AssignToRefactor sets p to p1 and return it
func (p *G1Affine) AssignToRefactor(cs *frontend.ConstraintSystem, p1 *G1Affine) *G1Affine {
	p.X = cs.Constant(p1.X)
	p.Y = cs.Constant(p1.Y)
	return p
}

// AddAssign adds 2 point in Jacobian coordinates
// p=p, a=p1
func (p *G1Jac) AddAssign(cs *frontend.ConstraintSystem, p1 *G1Jac) *G1Jac {

	// get some Element from our pool
	var Z1Z1, Z2Z2, U1, U2, S1, S2, H, I, J, r, V frontend.Variable

	Z1Z1 = cs.Mul(p1.Z, p1.Z)

	Z2Z2 = cs.Mul(p.Z, p.Z)

	U1 = cs.Mul(p1.X, Z2Z2)

	U2 = cs.Mul(p.X, Z1Z1)

	S1 = cs.Mul(p1.Y, p.Z)
	S1 = cs.Mul(S1, Z2Z2)

	S2 = cs.Mul(p.Y, p1.Z)
	S2 = cs.Mul(S2, Z1Z1)

	H = cs.Sub(U2, U1)

	I = cs.Add(H, H)
	I = cs.Mul(I, I)

	J = cs.Mul(H, I)

	r = cs.Sub(S2, S1)
	r = cs.Add(r, r)

	V = cs.Mul(U1, I)

	p.X = cs.Mul(r, r)
	p.X = cs.Sub(p.X, J)
	p.X = cs.Sub(p.X, V)
	p.X = cs.Sub(p.X, V)

	p.Y = cs.Sub(V, p.X)
	p.Y = cs.Mul(p.Y, r)

	S1 = cs.Mul(J, S1)
	S1 = cs.Add(S1, S1)

	p.Y = cs.Sub(p.Y, S1)

	p.Z = cs.Add(p.Z, p1.Z)
	p.Z = cs.Mul(p.Z, p.Z)
	p.Z = cs.Sub(p.Z, Z1Z1)
	p.Z = cs.Sub(p.Z, Z2Z2)
	p.Z = cs.Mul(p.Z, H)

	return p
}

// DoubleAssign doubles the receiver point in jacobian coords and returns it
func (p *G1Jac) DoubleAssign(cs *frontend.ConstraintSystem) *G1Jac {
	// get some Element from our pool
	var XX, YY, YYYY, ZZ, S, M, T frontend.Variable

	XX = cs.Mul(p.X, p.X)
	YY = cs.Mul(p.Y, p.Y)
	YYYY = cs.Mul(YY, YY)
	ZZ = cs.Mul(p.Z, p.Z)
	S = cs.Add(p.X, YY)
	S = cs.Mul(S, S)
	S = cs.Sub(S, XX)
	S = cs.Sub(S, YYYY)
	S = cs.Add(S, S)
	M = cs.Mul(XX, 3) // M = 3*XX+a*ZZ^2, here a=0 (we suppose sw has j invariant 0)
	p.Z = cs.Add(p.Z, p.Y)
	p.Z = cs.Mul(p.Z, p.Z)
	p.Z = cs.Sub(p.Z, YY)
	p.Z = cs.Sub(p.Z, ZZ)
	p.X = cs.Mul(M, M)
	T = cs.Add(S, S)
	p.X = cs.Sub(p.X, T)
	p.Y = cs.Sub(S, p.X)
	p.Y = cs.Mul(p.Y, M)
	YYYY = cs.Mul(YYYY, 8)
	p.Y = cs.Sub(p.Y, YYYY)

	return p
}

// Select sets p1 if b=1, p2 if b=0, and returns it. b must be boolean constrained
func (p *G1Affine) Select(cs *frontend.ConstraintSystem, b frontend.Variable, p1, p2 *G1Affine) *G1Affine {

	p.X = cs.Select(b, p1.X, p2.X)
	p.Y = cs.Select(b, p1.Y, p2.Y)

	return p

}

// FromJac sets p to p1 in affine and returns it
func (p *G1Affine) FromJac(cs *frontend.ConstraintSystem, p1 *G1Jac) *G1Affine {
	s := cs.Mul(p1.Z, p1.Z)
	p.X = cs.Div(p1.X, s)
	p.Y = cs.Div(p1.Y, cs.Mul(s, p1.Z))
	return p
}

// Double double a point in affine coords
func (p *G1Affine) Double(cs *frontend.ConstraintSystem, p1 *G1Affine) *G1Affine {

	var t, d, c1, c2, c3 big.Int
	t.SetInt64(3)
	d.SetInt64(2)
	c1.SetInt64(1)
	c2.SetInt64(-2)
	c3.SetInt64(-1)

	// compute lambda = (3*p1.x**2+a)/2*p1.y, here we assume a=0 (j invariant 0 curve)
	x2 := cs.Mul(p1.X, p1.X)
	cs.Mul(p1.X, p1.X)
	l1 := cs.LinearExpression(
		cs.Term(x2, &t),
	)
	l2 := cs.LinearExpression(
		cs.Term(p1.Y, &d),
	)
	l := cs.Div(l1, l2)

	// xr = lambda**2-p.x-p1.x
	_x := cs.LinearExpression(
		cs.Term(cs.Mul(l, l), &c1),
		cs.Term(p1.X, &c2),
	)

	// p.y = lambda(p.x-xr) - p.y
	t1 := cs.Mul(p1.X, l)
	t2 := cs.Mul(l, _x)
	l3 := cs.LinearExpression(
		cs.Term(t1, &c1),
		cs.Term(t2, &c3),
		cs.Term(p1.Y, &c3),
	)
	p.Y = cs.Mul(l3, 1)

	//p.x = xr
	p.X = cs.Mul(_x, 1)
	return p
}

// ScalarMul computes scalar*p1, affect the result to p, and returns it.
// n is the number of bits used for the scalar mul.
// TODO it doesn't work if the scalar if 1, because it ends up doing P-P at the end, involving division by 0
// TODO add a panic if scalar == 1
func (p *G1Affine) ScalarMul(cs *frontend.ConstraintSystem, p1 *G1Affine, s interface{}, n int) *G1Affine {

	scalar := cs.Constant(s)

	var base, res G1Affine
	base.Double(cs, p1)
	res.AssignToRefactor(cs, p1)

	b := cs.ToBinary(scalar, n)

	var tmp G1Affine

	// start from 1 and use right-to-left scalar multiplication to avoid bugs due to incomplete addition law
	// (I don't see how to avoid that)
	for i := 1; i <= n-1; i++ {
		tmp.AssignToRefactor(cs, &res).AddAssign(cs, &base)
		res.Select(cs, b[i], &tmp, &res)
		base.Double(cs, &base)
	}

	// now check the lsb, if it's one, leave the result as is, otherwise substract P
	tmp.Neg(cs, p1).AddAssign(cs, &res)

	p.Select(cs, b[0], &res, &tmp)

	return p

}

func bls377FpTobw761fr(a *fp.Element) (r fr.Element) {
	for i, v := range a {
		r[i] = v
	}
	return
}

// Assign a value to self (witness assignment)
func (p *G1Jac) Assign(p1 *bls377.G1Jac) {
	p.X.Assign(bls377FpTobw761fr(&p1.X))
	p.Y.Assign(bls377FpTobw761fr(&p1.Y))
	p.Z.Assign(bls377FpTobw761fr(&p1.Z))
}

// MustBeEqual constraint self to be equal to other into the given constraint system
func (p *G1Jac) MustBeEqual(cs *frontend.ConstraintSystem, other G1Jac) {
	cs.AssertIsEqual(p.X, other.X)
	cs.AssertIsEqual(p.Y, other.Y)
	cs.AssertIsEqual(p.Z, other.Z)
}

// Assign a value to self (witness assignment)
func (p *G1Affine) Assign(p1 *bls377.G1Affine) {
	p.X.Assign(bls377FpTobw761fr(&p1.X))
	p.Y.Assign(bls377FpTobw761fr(&p1.Y))
}

// MustBeEqual constraint self to be equal to other into the given constraint system
func (p *G1Affine) MustBeEqual(cs *frontend.ConstraintSystem, other G1Affine) {
	cs.AssertIsEqual(p.X, other.X)
	cs.AssertIsEqual(p.Y, other.Y)
}
