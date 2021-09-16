package Setup

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"math/big"

	hash2prime "../Hashtoprime"
)

func Set_setup(lambda, nu int) (n, f *big.Int) {

	pk, _ := rsa.GenerateKey(rand.Reader, lambda)
	var F *big.Int
	N := pk.PublicKey.N

	for {
		F, _ = rand.Int(rand.Reader, N)
		if F != pk.Primes[0] && F != pk.Primes[1] && F != big.NewInt(1) {
			break
		}
	}

	G := F.Exp(F, F, pk.PublicKey.N)

	return N, G
}

//Set Commitment function

func Set_commit(ck, U []*big.Int) (c, o *big.Int) {

	P := make([]*big.Int, len(U))
	N, G := ck[0], ck[1]

	for i, u := range U {
		P[i] = hash2prime.Hprime(u)
		G.Exp(G, P[i], N)
	}

	fmt.Println(P)
	return G, nil

}

func Set_ver(ck, U []*big.Int, Acc *big.Int) int {

	P := make([]*big.Int, len(U))
	N, G := ck[0], ck[1]

	for i, u := range U {
		P[i] = hash2prime.Hprime(u)
		G.Exp(G, P[i], N)
	}

	if Acc == G {
		return 1
	}
	return 0

}
