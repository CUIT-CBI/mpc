package main

import (
	"fmt"
	"math/rand"
	"time"
)

type SecretShare struct {
	index uint8
	value uint8
}

type Matrix struct {
	t      int
	matrix [][]int
	values []int
}

func main() {
	rand.Seed(time.Now().UnixNano())
	M := uint8(rand.Intn(255))
	n := uint8(10)
	t := uint8(3)
	coeff := makePoly(uint8(M), uint8(n), uint8(t))
	fmt.Println(coeff)
	shares := SplitShares(coeff, M, n)
	fmt.Println(shares)

	selected := make([]SecretShare, 3)
	selected[0] = shares[2]
	selected[1] = shares[5]
	selected[2] = shares[8]

	recovered := RecoverSecret(selected, int(t))
	fmt.Printf("recovered secret: %d\n", recovered)
}

func makePoly(M uint8, n uint8, t uint8) []uint8 {
	coeffs := make([]uint8, t)
	coeffs[0] = M
	for i := 1; uint8(i) < t; i++ {
		randomCoeff := rand.Intn(int(M))
		coeffs[i] = uint8(randomCoeff)
	}
	return coeffs
}

func SplitShares(coeff []uint8, s uint8, n uint8) []SecretShare {
	shares := make([]SecretShare, n)

	for i := 1; uint8(i) <= n; i++ {
		tmp := int(coeff[0])

		multip := i
		for j := 1; j <= len(coeff)-1; j++ {
			tmp += int(coeff[j]) * multip
			multip *= i
		}

		shares[i-1] = SecretShare{
			index: uint8(i),
			value: uint8(tmp % 256),
		}
	}

	return shares
}

func RecoverSecret(shares []SecretShare, t int) uint8 {
	matrix := Matrix{
		matrix: make([][]int, t),
		values: make([]int, t),
		t:      t,
	}
	for i := 0; i < len(shares); i++ {
		param := make([]int, t)
		param[0] = 1
		tmp := int(shares[i].index)
		for j := 1; j < t; j++ {
			param[j] = tmp
			tmp *= int(shares[i].index)
		}
		matrix.matrix[i] = param
		matrix.values[i] = int(shares[i].value)
	}

	fmt.Println(matrix.matrix)
	fmt.Print(matrix.values)
	return 0
}
