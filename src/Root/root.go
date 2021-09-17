package Root

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"

	hash2prime "../Hashtoprime"
)

func generate_random(bound *big.Int) *big.Int {

	upper := bound
	lower := bound.Neg(bound)

	//rand.Seed(time.Now().UnixNano())
	//n := a + rand.Intn(b-a+1)

	max := upper.Sub(upper, lower)
	max.Add(max, big.NewInt(1))

	random, _ := rand.Int(rand.Reader, max)
	random.Add(random, lower)

	return random
}

func generate_s(a, b, c *big.Int) *big.Int {

	s := new(big.Int).Mul(b, c)
	s.Sub(a, s)

	return s

}

func generate_alpha(G, H, N, x, y *big.Int) *big.Int {

	alpha := new(big.Int).Exp(G, x, N)
	temp := new(big.Int).Exp(H, y, N)
	alpha.Mul(alpha, temp)
	alpha.Mod(alpha, N)

	return alpha

}

func generate_alpha_ver(G, H, C, x, y, z, N *big.Int) *big.Int {

	alpha := new(big.Int).Exp(G, x, N)
	temp := new(big.Int).Exp(H, y, N)
	temp1 := new(big.Int).Exp(C, z, N)
	alpha.Mul(alpha, temp)
	alpha.Mul(alpha, temp1)

	return alpha
}

func Prove(crs, commitment, witness []*big.Int, lamda_s, lamda_z, mu int64) []*big.Int {

	N, G, H := crs[0], crs[1], crs[2]
	e, r, W := witness[0], witness[1], witness[2]
	bound := N
	bound.Div(bound, big.NewInt(4))
	bound.Sub(bound, big.NewInt(1))

	r2 := generate_random(bound)
	r3 := generate_random(bound)

	Cw := generate_alpha(W, H, big.NewInt(1), r2, N)
	Cr := generate_alpha(G, H, r2, r3, N)

	bound1 := big.NewInt(2)
	bound1.Exp(bound1, big.NewInt(lamda_s+lamda_z+mu), nil)
	re := generate_random(bound1)

	bound2 := big.NewInt(2)
	bound2.Exp(bound2, big.NewInt(lamda_s+lamda_z), nil)
	bound2.Mul(bound, bound2)
	rr := generate_random(bound2)
	rr2 := generate_random(bound2)
	rr3 := generate_random(bound2)

	bound3 := big.NewInt(2)
	bound3.Exp(bound3, big.NewInt(lamda_s+lamda_z+mu), nil)
	bound3.Mul(bound, bound3)
	r_beta := generate_random(bound3)
	r_delta := generate_random(bound3)

	alpha1 := generate_alpha(G, H, N, re, rr)
	alpha2 := generate_alpha(G, H, N, rr2, rr3)
	alpha3 := generate_alpha(Cw, H, N, re, r_beta.Neg(r_beta))
	alpha4 := generate_alpha_ver(G, H, Cr, r_beta.Neg(r_beta), r_delta.Neg(r_delta), re, N)

	list := []*big.Int{alpha1, alpha2, alpha3, alpha4, commitment[0], commitment[1]}
	h := sha256.New()
	for _, y := range list {
		h.Write(y.Bytes())
	}
	hash := fmt.Sprintf("%x", h.Sum(nil))
	c := new(big.Int)
	c.SetString(hash, 16)
	//fmt.Println(c)
	c = hash2prime.Fu(c)

	se := generate_s(re, c, e)
	sr := generate_s(rr, c, r)
	sr2 := generate_s(rr2, c, r2)
	sr3 := generate_s(rr3, c, r3)
	s_beta := generate_s(r_beta, c.Mul(c, e), r2)
	s_delta := generate_s(r_delta, c.Mul(c, e), r3)

	pi := []*big.Int{Cw, Cr, alpha1, alpha2, alpha3, alpha4, se, sr, sr2, sr3, s_beta, s_delta}
	return pi

}

func VerProof(crs, commitment, pi []*big.Int, lamda, lambda_s, mu int64) int {

	N, G, H := crs[0], crs[1], crs[2]
	Ce, Acc := commitment[0], commitment[1]
	alpha1, alpha2, alpha3, alpha4 := pi[2], pi[3], pi[4], pi[5]

	list := []*big.Int{alpha1, alpha2, alpha3, alpha4, Ce, Acc}
	h := sha256.New()
	for _, y := range list {
		h.Write(y.Bytes())
	}
	hash := fmt.Sprintf("%x", h.Sum(nil))
	c := new(big.Int)
	c.SetString(hash, 16)
	//fmt.Println(c)
	c = hash2prime.Fu(c)

	alpha_1 := generate_alpha_ver(Ce, G, H, c, pi[6], pi[7], N)
	alpha_2 := generate_alpha_ver(pi[1], G, H, c, pi[8], pi[9], N)
	alpha_3 := generate_alpha_ver(Acc, pi[0], H, c, pi[6], pi[9].Neg(pi[9]), N)
	alpha_4 := generate_alpha_ver(pi[1], G, H, pi[6], pi[9].Neg(pi[9]), pi[10].Neg(pi[10]), N)

	upper := big.NewInt(2)
	upper.Exp(upper, big.NewInt(lamda+lambda_s+mu+1), nil)
	lower := new(big.Int).Neg(upper)

	var se_bool bool
	upper_bound := new(big.Int).Sub(upper, pi[6])
	lower_bound := new(big.Int).Sub(pi[6], lower)
	if upper_bound.Sign() == 1 && lower_bound.Sign() == 1 {
		se_bool = true
	}

	if alpha1 == alpha_1 && alpha2 == alpha_2 && alpha3 == alpha_3 && alpha4 == alpha_4 && se_bool {
		return 1
	}

	return 0

}
