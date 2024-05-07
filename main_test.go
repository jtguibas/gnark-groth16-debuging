package main

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"testing"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/consensys/gnark/std/rangecheck"
)

type MyCircuit struct {
	X            frontend.Variable `gnark:",public"`
	Y            frontend.Variable `gnark:",public"`
	Z            frontend.Variable `gnark:",public"`
	DoRangeCheck bool
}

func (circuit *MyCircuit) Define(api frontend.API) error {
	api.AssertIsEqual(circuit.Z, api.Add(circuit.X, circuit.Y))
	if true || circuit.DoRangeCheck {
		rangeChecker := rangecheck.New(api)
		rangeChecker.Check(circuit.X, 8)
	}
	return nil
}

type Groth16ProofData struct {
	Proof  []string `json:"proof"`
	Inputs []string `json:"inputs"`
}

func TestGroth16(t *testing.T) {
	circuit := MyCircuit{DoRangeCheck: false}

	r1cs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)
	if err != nil {
		panic(err)
	}
	pk, vk, err := groth16.Setup(r1cs)
	if err != nil {
		panic(err)
	}

	buf := new(bytes.Buffer)
	err = vk.ExportSolidity(buf)
	if err != nil {
		panic(err)
	}
	content := buf.String()

	contractFile, err := os.Create("contracts/src/VerifierGroth16.sol")
	if err != nil {
		panic(err)
	}
	w := bufio.NewWriter(contractFile)
	_, err = w.Write([]byte(content))
	if err != nil {
		panic(err)
	}
	contractFile.Close()

	assignment := MyCircuit{
		X: 1,
		Y: 2,
		Z: 3,
	}

	witness, _ := frontend.NewWitness(&assignment, ecc.BN254.ScalarField())
	proof, _ := groth16.Prove(r1cs, pk, witness)

	const fpSize = 4 * 8
	buf = new(bytes.Buffer)
	proof.WriteRawTo(buf)
	proofBytes := buf.Bytes()

	proofs := make([]string, 8)
	// Print out the proof
	for i := 0; i < 8; i++ {
		proofs[i] = "0x" + hex.EncodeToString(proofBytes[i*fpSize:(i+1)*fpSize])
	}

	publicWitness, _ := witness.Public()
	publicWitnessBytes, _ := publicWitness.MarshalBinary()
	publicWitnessBytes = publicWitnessBytes[12:] // We cut off the first 12 bytes because they encode length information

	commitmentCountBigInt := new(big.Int).SetBytes(proofBytes[fpSize*8 : fpSize*8+4])
	commitmentCount := int(commitmentCountBigInt.Int64())

	var commitments []*big.Int = make([]*big.Int, 2*commitmentCount)
	var commitmentPok [2]*big.Int

	for i := 0; i < 2*commitmentCount; i++ {
		commitments[i] = new(big.Int).SetBytes(proofBytes[fpSize*8+4+i*fpSize : fpSize*8+4+(i+1)*fpSize])
	}

	commitmentPok[0] = new(big.Int).SetBytes(proofBytes[fpSize*8+4+2*commitmentCount*fpSize : fpSize*8+4+2*commitmentCount*fpSize+fpSize])
	commitmentPok[1] = new(big.Int).SetBytes(proofBytes[fpSize*8+4+2*commitmentCount*fpSize+fpSize : fpSize*8+4+2*commitmentCount*fpSize+2*fpSize])

	fmt.Println("Generating Fixture")

	fmt.Println("uint256[8] memory proofs = [")
	for i := 0; i < 8; i++ {
		fmt.Print(proofs[i])
		if i != 7 {
			fmt.Println(",")
		}
	}
	fmt.Println("];")
	fmt.Println()

	fmt.Println("uint256[2] memory commitments = [")
	for i := 0; i < 2*commitmentCount; i++ {
		fmt.Print(commitments[i])
		if i != 2*commitmentCount-1 {
			fmt.Println(",")
		}
	}
	fmt.Println("];")

	fmt.Println("uint256[2] memory commitmentPok = [")
	for i := 0; i < 2; i++ {
		fmt.Print(commitmentPok[i])
		if i != 1 {
			fmt.Println(",")
		}
	}
	fmt.Println("];")
	fmt.Println()

	fmt.Println("uint256[3] memory inputs = [")
	fmt.Println("uint256(1),")
	fmt.Println("uint256(2),")
	fmt.Println("uint256(3)")
	fmt.Println("];")
}
