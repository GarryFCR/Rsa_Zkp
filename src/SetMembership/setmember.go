package SetMembership

import (
	"crypto/rand"
	"math/big"

	hashEq "../HashEq"
	hash2prime "../Hashtoprime"
	root "../Root"
	pedersen "../Setup"
	modEq "../modEq"
)

func KeyGen(ck []*big.Int) []*big.Int {

	N := ck[0]
	var H *big.Int
	for {
		H, _ = rand.Int(rand.Reader, N)
		if new(big.Int).Mod(N, H) != big.NewInt(0) && H != big.NewInt(1) {
			break
		}
	}
	H.Exp(H, big.NewInt(2), N)

	crs := []*big.Int{ck[0], ck[1], H, ck[2], ck[3], ck[4]}
	return crs

}

func Prove(crs []*big.Int, Acc, Cu, cu, U, u, ru *big.Int) (Ce, ce *big.Int, proof_root, proof_mod, proof_hash []*big.Int) {

	e := hash2prime.Hprime(u)
	j := new(big.Int).Sub(e, hash2prime.Fu(u))
	ce, rq := pedersen.Pedersen_commit(crs[3:], crs[2], e)
	Ce, r := pedersen.Pedersen_commit(crs[1:3], crs[0], e)
	inverse := new(big.Int).ModInverse(e, crs[0])
	W := new(big.Int).Exp(Acc, inverse, crs[0])

	commit := []*big.Int{Ce, Cu}
	root_witness := []*big.Int{e, r, W}
	pi_root := root.Prove(crs[:3], commit, root_witness, int64(32), int64(32), int64(32))

	mod_witness := []*big.Int{e, e, r, rq}
	pi_mod := modEq.Prove(crs, commit, mod_witness, int64(64))

	hash_commitment := []*big.Int{ce, cu}
	hash_witness := []*big.Int{e, u, rq, ru, j}
	pi_hash := hashEq.Prove(crs[4:], hash_witness, hash_commitment, int64(32), int64(32), int64(32), int64(32))

	return Ce, ce, pi_root, pi_mod, pi_hash

}

func VerProof(crs []*big.Int, Cu, cu *big.Int, Ce, ce *big.Int, pi_root, pi_mod, pi_hash []*big.Int) int {

	commit := []*big.Int{Ce, Cu}
	root_bool := root.VerProof(crs[:3], commit, pi_root, int64(32), int64(32), int64(32))

	modEq_bool := modEq.VerProof(crs, commit, pi_mod)

	hash_commitment := []*big.Int{ce, cu}
	hash_bool := hashEq.VerProof(crs[4:], hash_commitment, pi_hash)

	if root_bool == 1 && modEq_bool == 1 && hash_bool == 1 {
		return 1
	}

	return 0

}
