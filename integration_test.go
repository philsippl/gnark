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

package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/consensys/gurvy"
	"github.com/philsippl/gnark/frontend"
	"github.com/philsippl/gnark/internal/backend/circuits"
	"github.com/philsippl/gnark/io"
)

func TestIntegration(t *testing.T) {
	// create temporary dir for integration test
	parentDir := "./integration_test"
	os.RemoveAll(parentDir)
	defer os.RemoveAll(parentDir)
	if err := os.MkdirAll(parentDir, 0700); err != nil {
		t.Fatal(err)
	}

	// spv: setup, prove, verify
	spv := func(curveID gurvy.ID, name string, _good, _bad frontend.Circuit) {
		t.Logf("%s circuit (%s)", name, curveID.String())
		// path for files
		fCircuit := filepath.Join(parentDir, name+".r1cs")
		fPk := filepath.Join(parentDir, name+".pk")
		fVk := filepath.Join(parentDir, name+".vk")
		fProof := filepath.Join(parentDir, name+".proof")
		fInputGood := filepath.Join(parentDir, name+".good.input")
		fInputBad := filepath.Join(parentDir, name+".bad.input")

		// 2: input files to disk
		good, err := frontend.ParseWitness(_good)
		if err != nil {
			panic("invalid good assignment:" + err.Error())
		}
		bad, err := frontend.ParseWitness(_bad)
		if err != nil {
			panic("invalid bad assignment:" + err.Error())
		}
		if err := io.WriteWitness(fInputGood, good); err != nil {
			t.Fatal(err)
		}
		if err := io.WriteWitness(fInputBad, bad); err != nil {
			t.Fatal(err)
		}

		// 3: run setup
		{
			cmd := exec.Command("go", "run", "main.go", "setup", fCircuit, "--pk", fPk, "--vk", fVk)
			out, err := cmd.CombinedOutput()
			if err != nil {
				t.Log(string(out))
				t.Fatal(err)
			}
		}

		pv := func(fInput string, expectedVerifyResult bool) {
			// 4: run prove
			{
				cmd := exec.Command("go", "run", "main.go", "prove", fCircuit, "--pk", fPk, "--input", fInput, "--proof", fProof)
				out, err := cmd.CombinedOutput()
				if expectedVerifyResult && err != nil {
					t.Log(string(out))
					t.Fatal(err)
				}
			}

			// note: here we ain't testing much when the prover failed. verify will not find a proof file, and that's it.

			// 4: run verify
			{
				cmd := exec.Command("go", "run", "main.go", "verify", fProof, "--vk", fVk, "--input", fInput)
				out, err := cmd.CombinedOutput()
				if expectedVerifyResult && err != nil {
					t.Log(string(out))
					t.Fatal(err)
				} else if !expectedVerifyResult && err == nil {
					t.Log(string(out))
					t.Fatal("verify should have failed but apparently succeeded")
				}
			}
		}

		pv(fInputGood, true)
		pv(fInputBad, false)
	}

	curves := []gurvy.ID{gurvy.BLS377, gurvy.BLS381, gurvy.BN256, gurvy.BW761}

	for name, circuit := range circuits.Circuits {
		if testing.Short() {
			if name != "lut01" && name != "frombinary" {
				continue
			}
		}
		for _, curve := range curves {
			// serialize to disk
			fCircuit := filepath.Join(parentDir, name+".r1cs")
			typedR1CS := circuit.R1CS.ToR1CS(curve)
			if err := io.WriteFile(fCircuit, typedR1CS); err != nil {
				t.Fatal(err)
			}
			spv(curve, name, circuit.Good, circuit.Bad)
		}
	}
}
