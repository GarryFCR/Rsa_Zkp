package SetMembership

import (
	"crypto/rand"
	"math/big"

	hashEq "github.com/GarryFCR/Rsazkp/src/HashEq"
	hash2prime "github.com/GarryFCR/Rsazkp/src/Hashtoprime"
	root "github.com/GarryFCR/Rsazkp/src/Root"
	pedersen "github.com/GarryFCR/Rsazkp/src/Setup"
	modEq "github.com/GarryFCR/Rsazkp/src/modEq"
)

func KeyGen(ck []big.Int) []big.Int {

	N := ck[0]
	var H *big.Int
	for {
		H, _ = rand.Int(rand.Reader, &N)
		if new(big.Int).GCD(nil, nil, H, &N).Cmp(big.NewInt(1)) == 0 {
			break
		}
	}
	H_ := *new(big.Int).Exp(H, big.NewInt(2), &N)

	crs := []big.Int{ck[0], ck[1], H_, ck[2], ck[3], ck[4]}
	return crs

}

func Prove(crs, U []big.Int, Cu, cu, u, ru big.Int) (Ce, ce big.Int, proof_root, proof_mod, proof_hash []big.Int) {

	e := hash2prime.Hprime(u)
	fu := hash2prime.Fu(u)
	j := new(big.Int).Sub(&e, &fu)
	ce, rq := pedersen.Pedersen_commit(crs[3:], crs[3], e)
	Ce, r := pedersen.Pedersen_commit(crs[:3], crs[0], e)

	Primes := make([]big.Int, len(U))
	G := crs[1]
	for i, u_dash := range U {
		Primes[i] = hash2prime.Hprime(u_dash)
		if u_dash.Cmp(&u) != 0 {

			G.Exp(&G, &Primes[i], &crs[0])
		}

	}

	W := G

	commit := []big.Int{Ce, Cu}
	root_witness := []big.Int{e, r, W}
	pi_root := root.Prove(crs[:3], commit, root_witness, int64(256), int64(256), int64(512))

	commit1 := []big.Int{Ce, ce}
	mod_witness := []big.Int{e, e, r, rq}
	pi_mod := modEq.Prove(crs, commit1, mod_witness, int64(64))

	hash_commitment := []big.Int{ce, cu}
	hash_witness := []big.Int{e, u, rq, ru, *j}
	pi_hash := hashEq.Prove(crs[3:], hash_commitment, hash_witness, int64(512), int64(512), int64(512), int64(512))

	//fmt.Printf("Size:%d\n", unsafe.Sizeof(pi_root)+unsafe.Sizeof(pi_hash)+unsafe.Sizeof(pi_mod)+unsafe.Sizeof(Ce)+unsafe.Sizeof(ce)+unsafe.Sizeof(crs))
	return Ce, ce, pi_root, pi_mod, pi_hash

}

func VerProof(crs []big.Int, Cu, cu big.Int, Ce, ce big.Int, pi_root, pi_mod, pi_hash []big.Int) int {

	commit := []big.Int{Ce, Cu}
	root_bool := root.VerProof(crs[:3], commit, pi_root, int64(512), int64(256), int64(512))

	commit1 := []big.Int{Ce, ce}
	modEq_bool := modEq.VerProof(crs, commit1, pi_mod)

	hash_commitment := []big.Int{ce, cu}
	hash_bool := hashEq.VerProof(crs[3:], hash_commitment, pi_hash)

	//fmt.Println(root_bool, modEq_bool, hash_bool)
	if root_bool == 1 && modEq_bool == 1 && hash_bool == 1 {
		return 1
	}

	return 0

}
