/*
	SetCommitment is a (non-hiding) commitment scheme for set of strings.
	In this set Commitment scheme we map the set members to set of prime using Hash2prime function
	and then the Set Commitment is build as an RSA Accumulator to the set of those primes derived from
	set of strings using Has2prime
*/

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
		if new(big.Int).Mod(F, pk.Primes[0]).Cmp(big.NewInt(0)) != 0 && new(big.Int).Mod(F, pk.Primes[1]).Cmp(big.NewInt(0)) != 0 && F != big.NewInt(1) {
			break
		}
	}

	G := new(big.Int).Exp(F, big.NewInt(2), pk.PublicKey.N)

	return *N, *G, *pk.Primes[0], *pk.Primes[1]
}

//Set_commit generates the RSA Accumulator for the set of primes derived from set of strings using Hash2prime

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
