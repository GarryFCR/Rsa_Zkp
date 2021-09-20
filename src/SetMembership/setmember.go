package SetMembership

import (
	"crypto/rand"
	"math/big"
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
