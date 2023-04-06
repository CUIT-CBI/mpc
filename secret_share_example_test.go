package mpc_test

import (
	"fmt"

	"github.com/CUIT-CBI/mpc"
)

func ExampleSplitAndCombain() {
	secrets := make([]byte, 256)
	for i := 0; i < 256; i++ {
		secrets[i] = byte(i)
	}

	n := 10
	threshold := 3
	parts := mpc.SplitSecret(secrets, n, threshold)

	recovered_secrets := mpc.CombainSecret(parts[2:5])
	if len(recovered_secrets) != len(secrets) {
		fmt.Println("Recovered secrets length is not correct!")
	}

	for i := 0; i < len(secrets); i++ {
		if secrets[i] != recovered_secrets[i] {
			fmt.Println("Recovered secrets is not correct!")
		}
	}
}
