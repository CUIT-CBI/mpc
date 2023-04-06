package mpc

import (
	"math/rand"
	"time"
)

type SecretPart struct {
	index     int
	Threshold int
	Value     []byte
}

type Poly func(uint8) uint8

func SplitSecret(secret []byte, n, t int) []SecretPart {
	rand.Seed(time.Now().UnixNano())

	result := make([]SecretPart, n)

	for i := 0; i < n; i++ {
		result[i] = SecretPart{
			index:     i + 1,
			Threshold: t,
			Value:     make([]byte, len(secret)),
		}
	}

	for i := 0; i < len(secret); i++ {
		parts := split(secret[i], uint8(n), uint8(t))
		for j := 0; j < n; j++ {
			result[j].Value[i] = parts[j]
		}
	}

	return result
}

func CombainSecret(parts []SecretPart) []byte {
	t := parts[0].Threshold
	secretLen := len(parts[0].Value)
	coeffs := make([][]int, t)
	values := make([][]int, t)

	for i := 0; i < t; i++ {
		// make coeffs
		coeffsRow := make([]int, t)
		coeffsRow[0] = 1
		index := parts[i].index
		mult := index
		for j := 1; j < t; j++ {
			coeffsRow[j] = mult
			mult *= index
		}
		coeffs[i] = coeffsRow

		// make values
		values[i] = make([]int, secretLen)
		for j := 0; j < secretLen; j++ {
			values[i][j] = int(parts[i].Value[j])
		}
	}

	coeffs, values = recursion(coeffs, values)
	secrets := make([]byte, secretLen)
	for i := 0; i < secretLen; i++ {
		v := values[0][i] / coeffs[0][0]
		v = v % 256
		if v < 0 {
			v += 256
		}
		secrets[i] = uint8(v)
	}

	return secrets
}

func split(M byte, n, t uint8) []byte {
	result := make([]byte, n)
	poly := makePoly(M, t)
	for i := uint8(0); i < n; i++ {
		result[i] = poly(i + 1)
	}
	return result
}

func makePoly(M, t byte) Poly {
	coeffs := make([]uint8, t)
	coeffs[0] = M
	for i := uint8(1); i < t; i++ {
		coeffs[i] = uint8(rand.Intn(256))
	}
	return func(index byte) byte {
		result := int(M)
		mult := int(index)
		for i := uint8(1); i < t; i++ {
			result += mult * int(coeffs[i])
			mult *= int(index)
		}
		return uint8(result % 256)
	}
}

func recursion(coeffs [][]int, values [][]int) ([][]int, [][]int) {
	t := len(coeffs)
	if t == 1 {
		return coeffs, values
	}

	coeffs, values = elimination(coeffs, values, t)
	return recursion(coeffs, values)
}

func elimination(coeffs [][]int, values [][]int, t int) ([][]int, [][]int) {
	newCoeffs := make([][]int, t-1)
	newValues := make([][]int, t-1)

	base := coeffs[0][t-1]
	for i := 0; i < t-1; i++ {
		coeffsRow := make([]int, t-1)
		mult := coeffs[i+1][t-1]
		for j := 0; j < t-1; j++ {
			coeffsRow[j] = coeffs[i+1][j]*base - coeffs[0][j]*mult
		}
		newCoeffs[i] = coeffsRow

		secretLen := len(values[0])
		valuesRow := make([]int, secretLen)
		for j := 0; j < secretLen; j++ {
			valuesRow[j] = values[i+1][j]*base - values[0][j]*mult
		}
		newValues[i] = valuesRow
	}

	return newCoeffs, newValues
}
