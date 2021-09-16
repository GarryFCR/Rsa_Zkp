package Root

import (
	"crypto/rand"
	"math"
	"math/big"
)

func generate_random(bound *big.Int) *big.Int {

	upper := bound
	//lower:=bound.Mul(bound,big.NewInt(-1))
	lower := bound.Neg(bound)

	//rand.Seed(time.Now().UnixNano())
	//n := a + rand.Intn(b-a+1)

	max := upper.Sub(upper, lower)
	max.Add(max, big.NewInt(1))

	random, _ := rand.Int(rand.Reader, max)
	random.Add(random, lower)

	return random
}

func Prove(crs, commitment, witness []*big.Int, lamda_s, lamda_z, mu int64) []*big.Int {

	N, G, H := crs[0], crs[1], crs[2]
	e, r, W := witness[0], witness[1], witness[2]
	bound := N
	bound.Div(bound, big.NewInt(4))
	bound.Sub(bound, big.NewInt(1))

	r2 := generate_random(bound)
	r3 := generate_random(bound)

	var Cw, Cr, temp *big.Int
	Cw.Exp(H, r2, N)
	Cw.Mul(Cw, W)
	Cw.Mod(Cw, N)

	Cr.Exp(G, r2, nil)
	temp.Exp(H, r3, nil)
	Cr.Mul(Cr, temp)
	Cr.Mod(Cr, N)

	bound1 := math.Pow(2, float64(lamda_s+lamda_z+mu))
	re := generate_random(big.NewInt(int64(bound1)))

	bound2 := big.NewInt(int64(math.Pow(2, float64(lamda_s+lamda_z))))
	bound2.Mul(bound, bound2)
	rr := generate_random(bound2)
	rr2 := generate_random(bound2)
	rr3 := generate_random(bound2)

	bound3 := big.NewInt(int64(math.Pow(2, float64(lamda_s+lamda_z+mu))))
	bound3.Mul(bound, bound3)
	r_beta := generate_random(bound3)
	r_delta := generate_random(bound3)

	var alpha1, alpha2, alpha3, alpha4, temp1 *big.Int
	alpha1.Exp(G, rr, N)
	temp.Exp(H, re, N)
	alpha1.Mul(alpha1, temp)
	alpha1.Mod(alpha1, N)

	alpha2.Exp(G, rr2, N)
	temp.Exp(H, rr3, N)
	alpha2.Mul(alpha2, temp)
	alpha2.Mod(alpha2, N)

	alpha3.Exp(Cw, re, N)
	temp.Exp(H, r_beta.Neg(r_beta), N)
	alpha3.Mul(alpha3, temp)
	alpha3.Mod(alpha3, N)

	alpha4.Exp(Cr, re, N)
	temp.Exp(H, r_beta.Neg(r_delta), N)
	temp1.Exp(G, r_beta.Neg(r_beta), N)
	alpha4.Mul(alpha4, temp)
	alpha4.Mul(alpha4, temp1)
	alpha4.Mod(alpha4, N)
}
