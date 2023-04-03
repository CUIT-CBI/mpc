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

	fmt.Println(secret)
	parts := SplitSecrets(secret, 10, 4)
	recoveredSecret := CombainParts(parts[2:6])
	fmt.Println(recoveredSecret)
	fmt.Printf("Recovered secret: %s\n", string(recoveredSecret))
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

func CombainParts(parts []SecretUnit) []byte {
	t := parts[0].threshold
	indexes := make([]uint8, t)

	values := make([][]int, t)
	for i := 0; i < len(parts); i++ {
		result := make([]int, len(parts[i].data))
		for j := 0; j < len(parts[i].data); j++ {
			result[j] = int(parts[i].data[j])
		}
		values[i] = result
		indexes[i] = parts[i].index
	}

	coeffs := makeRecoveryCoeffs(indexes)

	coeffs, values = extinction(coeffs, values)

	return remainder(coeffs, values)
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

func makeRecoveryCoeffs(coeffs []byte) [][]int {
	coeffLen := len(coeffs)
	results := make([][]int, coeffLen)
	for i := 0; i < coeffLen; i++ {
		row := make([]int, coeffLen)
		row[0] = 1
		tmp := int(coeffs[i])
		for j := 1; j < coeffLen; j++ {
			row[j] = tmp
			tmp *= int(coeffs[i])
		}
		results[i] = row
	}
	return results
}

func extinction(coeffs [][]int, values [][]int) ([][]int, [][]int) {
	if len(coeffs) == 1 {
		return coeffs, values
	}

	valueLen := len(values[0])
	newCoeffs := make([][]int, len(coeffs)-1)
	newValues := make([][]int, len(coeffs)-1)

	for i := 1; i < len(coeffs); i++ {
		coeffsRow := make([]int, len(coeffs)-1)
		base := coeffs[0][len(coeffs)-1]
		mult := coeffs[i][len(coeffs)-1]
		for j := 0; j < len(coeffs)-1; j++ {
			coeffsRow[j] = coeffs[i][j]*base - coeffs[0][j]*mult
		}
		newCoeffs[i-1] = coeffsRow

		valueRow := make([]int, valueLen)
		for j := 0; j < valueLen; j++ {
			valueRow[j] = values[i][j]*base - values[0][j]*mult
		}
		newValues[i-1] = valueRow
	}

	return extinction(newCoeffs, newValues)
}

func remainder(coeffs [][]int, values [][]int) []byte {
	if len(coeffs) != 1 || len(values) != 1 {
		return []byte{}
	}
	r := coeffs[0][0]
	value := values[0]
	result := make([]byte, len(value))
	for i := 0; i < len(value); i++ {
		v := value[i] / r
		v = v % 257
		if v < 0 {
			v += 257
		}
		result[i] = uint8(v)
	}
	return result
}
