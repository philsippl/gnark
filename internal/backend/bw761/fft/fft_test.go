// Copyright 2020 ConsenSys AG
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by gnark/internal/generators DO NOT EDIT

package fft

import (
	"math/big"
	"strconv"
	"testing"

	"github.com/consensys/gurvy/bw761/fr"
	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

// --------------------------------------------------------------------
// utils

func evaluatePolynomial(pol []fr.Element, val fr.Element) fr.Element {
	var acc, res, tmp fr.Element
	res.Set(&pol[0])
	acc.Set(&val)
	for i := 1; i < len(pol); i++ {
		tmp.Mul(&acc, &pol[i])
		res.Add(&res, &tmp)
		acc.Mul(&acc, &val)
	}
	return res
}

// --------------------------------------------------------------------
// tests

func TestFFT(t *testing.T) {

	var maxSize int
	maxSize = 1 << 10

	domain := NewDomain(maxSize)

	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 5

	properties := gopter.NewProperties(parameters)

	properties.Property("DIF FFT should be consistent with dual basis", prop.ForAll(

		// checks that a random evaluation of a dual function eval(gen**ithpower) is consistent with the FFT result
		func(ithpower int) bool {

			pol := make([]fr.Element, maxSize)
			backupPol := make([]fr.Element, maxSize)

			for i := 0; i < maxSize; i++ {
				pol[i].SetRandom()
			}
			copy(backupPol, pol)

			domain.FFT(pol, DIF)
			BitReverse(pol)

			sample := domain.Generator
			sample.Exp(sample, big.NewInt(int64(ithpower)))

			eval := evaluatePolynomial(backupPol, sample)

			return eval.Equal(&pol[ithpower])

		},
		gen.IntRange(0, maxSize-1),
	))

	properties.Property("DIT FFT should be consistent with dual basis", prop.ForAll(

		// checks that a random evaluation of a dual function eval(gen**ithpower) is consistent with the FFT result
		func(ithpower int) bool {

			pol := make([]fr.Element, maxSize)
			backupPol := make([]fr.Element, maxSize)

			for i := 0; i < maxSize; i++ {
				pol[i].SetRandom()
			}
			copy(backupPol, pol)

			BitReverse(pol)
			domain.FFT(pol, DIT)

			sample := domain.Generator
			sample.Exp(sample, big.NewInt(int64(ithpower)))

			eval := evaluatePolynomial(backupPol, sample)

			return eval.Equal(&pol[ithpower])

		},
		gen.IntRange(0, maxSize-1),
	))

	properties.Property("bitReverse(DIF FFT(DIT FFT (bitReverse))))==id", prop.ForAll(

		func() bool {

			pol := make([]fr.Element, maxSize)
			backupPol := make([]fr.Element, maxSize)

			for i := 0; i < maxSize; i++ {
				pol[i].SetRandom()
			}
			copy(backupPol, pol)

			BitReverse(pol)
			domain.FFT(pol, DIT)
			domain.FFTInverse(pol, DIF)
			BitReverse(pol)

			check := true
			for i := 0; i < len(pol); i++ {
				check = check && pol[i].Equal(&backupPol[i])
			}
			return check
		},
	))

	properties.Property("DIT FFT(DIF FFT)==id", prop.ForAll(

		func() bool {

			pol := make([]fr.Element, maxSize)
			backupPol := make([]fr.Element, maxSize)

			for i := 0; i < maxSize; i++ {
				pol[i].SetRandom()
			}
			copy(backupPol, pol)

			domain.FFTInverse(pol, DIF)
			domain.FFT(pol, DIT)

			check := true
			for i := 0; i < len(pol); i++ {
				check = check && (pol[i] == backupPol[i])
			}
			return check
		},
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))

}

// --------------------------------------------------------------------
// benches
func BenchmarkBitReverse(b *testing.B) {

	var maxSize uint
	maxSize = 1 << 20

	pol := make([]fr.Element, maxSize)
	for i := uint(0); i < maxSize; i++ {
		pol[i].SetRandom()
	}

	for i := 8; i < 20; i++ {
		b.Run("bit reversing 2**"+strconv.Itoa(i)+"bits", func(b *testing.B) {
			_pol := make([]fr.Element, 1<<i)
			copy(_pol, pol)
			b.ResetTimer()
			for j := 0; j < b.N; j++ {
				BitReverse(_pol)
			}
		})
	}

}

func BenchmarkFFT(b *testing.B) {

	var maxSize uint
	maxSize = 1 << 20

	pol := make([]fr.Element, maxSize)
	for i := uint(0); i < maxSize; i++ {
		pol[i].SetRandom()
	}

	for i := 8; i < 20; i++ {
		b.Run("fft 2**"+strconv.Itoa(i)+"bits", func(b *testing.B) {
			sizeDomain := 1 << i
			_pol := make([]fr.Element, sizeDomain)
			copy(_pol, pol)
			domain := NewDomain(sizeDomain)
			b.ResetTimer()
			for j := 0; j < b.N; j++ {
				domain.FFT(_pol, DIT)
			}
		})
	}

}
