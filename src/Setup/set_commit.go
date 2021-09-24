package Setup

import (
	"crypto/rand"
	"crypto/rsa"
	"math/big"

	hash2prime "../Hashtoprime"
)

func Set_setup(lambda, mu int) (n, f, p, q big.Int) {

	pk, _ := rsa.GenerateKey(rand.Reader, lambda)
	var F *big.Int
	N := pk.PublicKey.N

	for {
		F, _ = rand.Int(rand.Reader, N)
		if F != pk.Primes[0] && F != pk.Primes[1] && F != big.NewInt(1) {
			break
		}
	}

	G := new(big.Int).Exp(F, big.NewInt(2), pk.PublicKey.N)

	return *N, *G, *pk.Primes[0], *pk.Primes[1]
}

//Set Commitment function

func Set_commit(ck, U []big.Int) (c, o big.Int) {

	P := make([]big.Int, len(U))
	N, G := ck[0], ck[1]

	for i, u := range U {
		P[i] = hash2prime.Hprime(u)
		G.Exp(&G, &P[i], &N)
	}

	//fmt.Println(P)
	return G, *big.NewInt(0)

}

func Set_ver(ck, U []big.Int, Acc big.Int) int {

	P := make([]big.Int, len(U))
	N, G := ck[0], ck[1]

	for i, u := range U {
		P[i] = hash2prime.Hprime(u)
		G.Exp(&G, &P[i], &N)
	}

	if Acc.Cmp(&G) == 0 {
		return 1
	}
	return 0

}
