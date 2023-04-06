package mpc_test

import (
	"testing"

	"github.com/CUIT-CBI/mpc"
)

func TestSplitSecret(t *testing.T) {
	secret := make([]byte, 256)
	for i := 0; i < 256; i++ {
		secret[i] = byte(i)
	}

	n := 10
	threshold := 3
	parts := mpc.SplitSecret(secret, n, threshold)

	if len(parts) != n {
		t.FailNow()
	}

	for i := 0; i < n; i++ {
		part := parts[i]
		if part.Threshold != threshold {
			t.FailNow()
		}
		if len(part.Value) != len(secret) {
			t.FailNow()
		}
	}
}

func TestCombainSecret(t *testing.T) {
	secrets := make([]byte, 256)
	for i := 0; i < 256; i++ {
		secrets[i] = byte(i)
	}

	n := 10
	threshold := 3
	parts := mpc.SplitSecret(secrets, n, threshold)

	recovered_secrets := mpc.CombainSecret(parts[2:5])
	if len(recovered_secrets) != len(secrets) {
		t.FailNow()
	}

	for i := 0; i < len(secrets); i++ {
		if secrets[i] != recovered_secrets[i] {
			t.FailNow()
		}
	}
}

func BenchmarkTestSplitSecret(b *testing.B) {
	secret := make([]byte, 256)
	for i := 0; i < 256; i++ {
		secret[i] = byte(i)
	}

	n := 10
	threshold := 3

	for i := 0; i < b.N; i++ {
		mpc.SplitSecret(secret, n, threshold)
	}
}

func BenchmarkTestCombainSecret(b *testing.B) {
	secret := make([]byte, 256)
	for i := 0; i < 256; i++ {
		secret[i] = byte(i)
	}

	n := 10
	threshold := 3

	parts := mpc.SplitSecret(secret, n, threshold)
	selected := parts[3:6]

	for i := 0; i < b.N; i++ {
		mpc.CombainSecret(selected)
	}
}
