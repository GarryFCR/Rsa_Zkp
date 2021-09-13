package main

import (
	"fmt"
	"math/big"

	setup "./Setup"
)

func main() {
	prime, g, h := setup.Pedersen_setup(32, 32)
	//fmt.Println(prime, g, h)

	ck := []*big.Int{prime, g, h}
	u := big.NewInt(12345)

	c, o := setup.Pedersen_commit(ck, u)
	//fmt.Println(c, o)

	ver := setup.Pedersen_ver(ck, c, u, o)
	if ver == 1 {
		fmt.Println("Commitment established")
	}
}
