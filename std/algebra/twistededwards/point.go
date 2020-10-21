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

package twistededwards

import (
	"math/big"

	"github.com/philsippl/gnark/backend"
	"github.com/philsippl/gnark/frontend"
)

// Point point on a twisted Edwards curve in a Snark cs
type Point struct {
	X, Y frontend.Variable
}

// MustBeOnCurve checks if a point is on the twisted Edwards curve
// ax^2 + y^2 = 1 + d*x^2*y^2
func (p *Point) MustBeOnCurve(cs *frontend.ConstraintSystem, curve EdCurve) {

	one := big.NewInt(1)

	l1 := cs.LinearExpression(cs.Term(p.X, &curve.A))
	l2 := cs.LinearExpression(cs.Term(p.X, one))
	axx := cs.Mul(l1, l2)
	yy := cs.Mul(p.Y, p.Y)
	lhs := cs.Add(axx, yy)

	l1 = cs.LinearExpression(cs.Term(p.X, &curve.D))
	l2 = cs.LinearExpression(cs.Term(p.X, one))
	dxx := cs.Mul(l1, l2)
	dxxyy := cs.Mul(dxx, yy)
	rhs := cs.Add(dxxyy, one)

	cs.AssertIsEqual(lhs, rhs)

}

// AddFixedPoint Adds two points, among which is one fixed point (the base), on a twisted edwards curve (eg jubjub)
// p1, base, ecurve are respectively: the point to add, a known base point, and the parameters of the twisted edwards curve
func (p *Point) AddFixedPoint(cs *frontend.ConstraintSystem, p1 *Point, x, y interface{}, curve EdCurve) *Point {

	// https://eprint.iacr.org/2008/013.pdf

	var dxy, bxa big.Int
	bx := backend.FromInterface(x)
	by := backend.FromInterface(y)
	dxy.Mul(&bx, &by).Mul(&dxy, &curve.D)
	bxa.Mul(&bx, &curve.A).Neg(&bxa)                                   // -ax1
	n1 := cs.LinearExpression(cs.Term(p1.X, &by), cs.Term(p1.Y, &bx))  // x1y2+x2y1
	n2 := cs.LinearExpression(cs.Term(p1.Y, &by), cs.Term(p1.X, &bxa)) // y1y2-ax1x2
	ld := cs.LinearExpression(cs.Term(p1.X, &dxy))                     // dx1x2y2
	_d := cs.Mul(ld, p1.Y)                                             // dx1x2y2y1
	d1 := cs.Add(1, _d)                                                // 1+dx1x2y2y1
	d2 := cs.Sub(1, _d)                                                // 1-dx1x2y2y1

	p.X = cs.Div(n1, d1)
	p.Y = cs.Div(n2, d2)

	return p
}

// AddGeneric Adds two points on a twisted edwards curve (eg jubjub)
// p1, p2, c are respectively: the point to add, a known base point, and the parameters of the twisted edwards curve
func (p *Point) AddGeneric(cs *frontend.ConstraintSystem, p1, p2 *Point, curve EdCurve) *Point {

	// https://eprint.iacr.org/2008/013.pdf
	res := Point{}

	one := big.NewInt(1)

	oneWire := cs.Constant(one)

	beta := cs.Mul(p1.X, p2.Y)
	gamma := cs.Mul(p1.Y, p2.X)
	delta := cs.Mul(p1.Y, p2.Y)
	epsilon := cs.Mul(p1.X, p2.X)
	tau := cs.Mul(delta, epsilon)
	num := cs.LinearExpression(
		cs.Term(beta, one),
		cs.Term(gamma, one),
	)
	den := cs.LinearExpression(
		cs.Term(oneWire, one),
		cs.Term(tau, &curve.D),
	)
	res.X = cs.Div(num, den)
	var minusa big.Int
	minusa.Neg(&curve.A).Mod(&minusa, &curve.Modulus)
	num = cs.LinearExpression(
		cs.Term(delta, one),
		cs.Term(epsilon, &minusa),
	)
	var minusd big.Int
	minusd.Neg(&curve.D).Mod(&minusd, &curve.Modulus)
	den = cs.LinearExpression(
		cs.Term(oneWire, one),
		cs.Term(tau, &minusd),
	)
	res.Y = cs.Div(num, den)

	p.X = res.X
	p.Y = res.Y
	return p
}

// Double doubles a points in SNARK coordinates
func (p *Point) Double(cs *frontend.ConstraintSystem, p1 *Point, curve EdCurve) *Point {
	p.AddGeneric(cs, p1, p1, curve)
	return p
}

// ScalarMulNonFixedBase computes the scalar multiplication of a point on a twisted Edwards curve
// p1: base point (as snark point)
// curve: parameters of the Edwards curve
// scal: scalar as a SNARK constraint
// Standard left to right double and add
func (p *Point) ScalarMulNonFixedBase(cs *frontend.ConstraintSystem, p1 *Point, scalar frontend.Variable, curve EdCurve) *Point {

	// first unpack the scalar
	b := cs.ToBinary(scalar, 256)

	res := Point{
		cs.Constant(0),
		cs.Constant(1),
	}

	for i := len(b) - 1; i >= 0; i-- {
		res.Double(cs, &res, curve)
		tmp := Point{}
		tmp.AddGeneric(cs, &res, p1, curve)
		res.X = cs.Select(b[i], tmp.X, res.X)
		res.Y = cs.Select(b[i], tmp.Y, res.Y)
	}

	p.X = res.X
	p.Y = res.Y
	return p
}

// ScalarMulFixedBase computes the scalar multiplication of a point on a twisted Edwards curve
// x, y: coordinates of the base point
// curve: parameters of the Edwards curve
// scal: scalar as a SNARK constraint
// Standard left to right double and add
func (p *Point) ScalarMulFixedBase(cs *frontend.ConstraintSystem, x, y interface{}, scalar frontend.Variable, curve EdCurve) *Point {

	// first unpack the scalar
	b := cs.ToBinary(scalar, 256)

	res := Point{
		cs.Constant(0),
		cs.Constant(1),
	}

	for i := len(b) - 1; i >= 0; i-- {
		res.Double(cs, &res, curve)
		tmp := Point{}
		tmp.AddFixedPoint(cs, &res, x, y, curve)
		res.X = cs.Select(b[i], tmp.X, res.X)
		res.Y = cs.Select(b[i], tmp.Y, res.Y)
	}

	p.X = res.X
	p.Y = res.Y
	return p
}
