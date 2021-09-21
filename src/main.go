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

	c, o := setup.Pedersen_commit(ck, prime, u)
	//fmt.Println(c, o)

	ver := setup.Pedersen_ver(ck, c, u, o)
	if ver == 1 {
		fmt.Println("Commitment established")
	}

	N, G := setup.Set_setup(32, 32)
	fmt.Println(N, G)

	ck1 := []*big.Int{N, G}
	set := []*big.Int{big.NewInt(12342), big.NewInt(12343), big.NewInt(12344), big.NewInt(12345)}

	com, _ := setup.Set_commit(ck1, set)
	fmt.Println(com)

	ver1 := setup.Set_ver(ck1, set, com)
	fmt.Println(ver1)

	//alpha1, alpha2, s_e, s_r, s_r_dash := modEq.Prove(N, g, h, prime, g, h,  )

}
