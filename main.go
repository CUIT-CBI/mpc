package main

import (
	"fmt"
	"math/rand"
)

type SecretShareByte struct {
	index uint8
	value uint8
}

type Matrix struct {
	t      int
	matrix [][]int
	values []int
}

type SecretUnit struct {
	threshold uint8
	index     uint8
	data      []byte
}

func main() {
	s := "Hello world!"
	secret := []byte(s)

	parts := SplitSecrets(secret, 10, 4)
	fmt.Println(parts[3:6])
	recoveredSecret := CombainParts(parts[2:6], 4)
	fmt.Println(string(recoveredSecret))
	// rand.Seed(time.Now().UnixNano())
	// M := uint8(rand.Intn(255))
	// n := uint8(10)
	// t := uint8(3)

	// coeff := makePoly(uint8(M), uint8(n), uint8(t))
	// fmt.Println(M)
	// shares := split(coeff, M, n)
	// fmt.Println(shares)

	// selected := make([]SecretShareByte, 3)
	// selected[0] = shares[2]
	// selected[1] = shares[5]
	// selected[2] = shares[8]

	// recovered := RecoverSecret(selected, int(t))
	// fmt.Printf("recovered secret: %d\n", recovered)
}

func SplitSecrets(secret []byte, n uint8, t uint8) []SecretUnit {
	result := make([]SecretUnit, n)
	secretLen := len(secret)
	for i := uint8(0); i < n; i++ {
		result[i] = SecretUnit{
			threshold: t,
			index:     i + 1,
			data:      make([]byte, secretLen),
		}
	}

	for i := 0; i < secretLen; i++ {
		parts := splitByte(secret[i], n, t)
		for j := uint8(0); j < n; j++ {
			result[j].data[i] = parts[j]
		}
	}

	return result
}

func CombainParts(parts []SecretUnit, t uint8) []byte {
	secretLen := len(parts[0].data)
	secret := make([]byte, secretLen)
	for i := uint8(0); i < uint8(secretLen); i++ {

	}
	return secret
}

func splitByte(M uint8, n, t uint8) []uint8 {
	coeffs := makePoly(M, n, t)
	return split(coeffs, M, n)
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

func split(coeff []uint8, M uint8, n uint8) []byte {
	shares := make([]byte, n)

	for i := 1; uint8(i) <= n; i++ {
		tmp := int(coeff[0])

		multip := i
		for j := 1; j <= len(coeff)-1; j++ {
			tmp += int(coeff[j]) * multip
			multip *= i
		}

		shares[i-1] = uint8(tmp % 257)
	}

	return shares
}

func RecoverSecret(coeffs, shares []byte, t int) uint8 {
	matrix := Matrix{
		matrix: make([][]int, t),
		values: make([]int, t),
		t:      t,
	}
	for i := 0; i < len(shares); i++ {
		param := make([]int, t)
		param[0] = 1
		tmp := int(coeffs[i])
		for j := 1; j < t; j++ {
			param[j] = tmp
			tmp *= int(coeffs[i])
		}
		matrix.matrix[i] = param
		matrix.values[i] = int(shares[i])
	}

	matrix = extinction(matrix)

	v := (matrix.values[0] / matrix.matrix[0][0]) % 257
	if v < 0 {
		v = v + 257
	}

	return uint8(v)
}

func extinction(matrix Matrix) Matrix {
	if matrix.t == 1 {
		return matrix
	}

	multiplication := 1

	for i := 0; i < matrix.t; i++ {

		multiplication *= matrix.matrix[i][matrix.t-1]
	}

	for i := 0; i < matrix.t; i++ {
		tmp := multiplication / matrix.matrix[i][matrix.t-1]
		for j := 0; j < matrix.t-1; j++ {
			matrix.matrix[i][j] *= tmp
		}
		matrix.values[i] *= tmp
	}

	newMatrix := Matrix{
		matrix: make([][]int, matrix.t-1),
		values: make([]int, matrix.t-1),
		t:      matrix.t - 1,
	}

	for i := 0; i < newMatrix.t; i++ {
		newMatrix.values[i] = matrix.values[i+1] - matrix.values[0]
		row := make([]int, newMatrix.t)
		for j := 0; j < newMatrix.t; j++ {
			row[j] = matrix.matrix[i+1][j] - matrix.matrix[0][j]
		}
		newMatrix.matrix[i] = row
	}

	return extinction(newMatrix)
}
