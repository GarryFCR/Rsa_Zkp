package hashEq

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"

	hash2prime "../Hashtoprime"
	generate "../Root"
)

func Prove(crs, commitment, witness []big.Int, lambda_s, lambda_z, mu, eta int64) []big.Int {

	q, g, h := crs[0], crs[1], crs[2]
	ce, cu := commitment[0], commitment[1]
	e, u, rq, ru, j := witness[0], witness[1], witness[2], witness[3], witness[4]

	//fmt.Println(&g, &q, new(big.Int).Exp(big.NewInt(6), big.NewInt(13), big.NewInt(13)))
	//Sampling
	bound := big.NewInt(2)
	bound.Exp(bound, big.NewInt(lambda_s+lambda_z+mu), &q)
	re := generate.Generate_random(*bound)
	re_dash := generate.Generate_random(*bound)
	rj := generate.Generate_random(*bound)

	rr, _ := rand.Int(rand.Reader, &q)
	rr_dash, _ := rand.Int(rand.Reader, &q)

	bound1 := big.NewInt(2)
	bound1.Exp(bound1, big.NewInt(lambda_s+lambda_z+eta), &q)
	r1 := generate.Generate_random(*bound1)
	r2 := generate.Generate_random(*bound1)
	rr1 := generate.Generate_random(*bound1)
	rr2 := generate.Generate_random(*bound1)

	H_u := hash2prime.Fu(u)
	Ch := *new(big.Int).Exp(&g, &H_u, &q)

	//power = 4log2(u+1)^2
	temp := new(big.Int).Add(&u, big.NewInt(2))
	power := new(big.Int).Div(&H_u, temp)
	power.Mul(power, big.NewInt(2))
	Cl := *new(big.Int).Exp(&g, power, &q)

	Cj := generate.Generate_alpha(g, h, q, j, rj)

	//Calculation of Alpha
	alpha1 := generate.Generate_alpha(g, h, q, re, *rr)
	alpha2 := generate.Generate_alpha(g, h, q, rr1, rr2)
	alpha4 := generate.Generate_alpha(ce, h, q, re_dash, *rr_dash)

	//c=H(hash)
	list := []big.Int{alpha1, alpha2, alpha4, ce, cu}
	h1 := sha256.New()
	for _, y := range list {
		h1.Write(y.Bytes())
	}
	hash := fmt.Sprintf("%x", h1.Sum(nil))
	c_hash := new(big.Int)
	c_hash.SetString(hash, 16)
	c := hash2prime.Fu(*c_hash)

	h_inverse := new(big.Int).ModInverse(&h, &q)
	//power1=2*c*ru*log2(u+1)^2
	power1 := new(big.Int).Mul(power, &ru)
	power1.Mul(power1, &c)
	power1.Div(power1, big.NewInt(2))
	alpha3 := generate.Generate_alpha_ver(cu, Cl, *h_inverse, r1, r2, *power1, q)

	//schnors authentication
	se := generate.Generate_s(re, c, e)
	sr := generate.Generate_s(*rr, c, rq)
	srr1 := generate.Generate_s(rr1, c, u)
	srr2 := generate.Generate_s(rr2, c, ru)
	sr1 := generate.Generate_s(r1, c, *new(big.Int).Div(power, big.NewInt(2)))
	sr2 := generate.Generate_s(r2, c, q)
	sr_dash := generate.Generate_s(*rr_dash, c, rj)
	sr_dash.Add(&sr_dash, new(big.Int).Mul(&c, &rq))
	se_dash := generate.Generate_s(re_dash, c, q)

	pi := []big.Int{Cl, Ch, Cj, alpha1, alpha2, alpha3, alpha4, se, sr, srr1, srr2, sr1, sr2, se_dash, sr_dash}

	return pi

}

func VerProof(crs, commitment, pi []big.Int) int {

	ce, cu := commitment[0], commitment[1]
	q, g, h := crs[0], crs[1], crs[2]
	alpha1, alpha2, alpha3, alpha4 := pi[3], pi[4], pi[5], pi[6]
	se, sr, srr1, srr2, sr1, sr2, se_dash, sr_dash := pi[7], pi[8], pi[9], pi[10], pi[11], pi[12], pi[13], pi[14]

	//c
	list := []big.Int{alpha1, alpha2, alpha4, ce, cu}
	h1 := sha256.New()
	for _, y := range list {
		h1.Write(y.Bytes())
	}
	hash := fmt.Sprintf("%x", h1.Sum(nil))
	c_hash := new(big.Int)
	c_hash.SetString(hash, 16)
	c := hash2prime.Fu(*c_hash)

	alpha_1 := generate.Generate_alpha_ver(ce, g, h, c, se, sr, q)
	alpha_2 := generate.Generate_alpha_ver(cu, g, h, c, srr1, srr2, q)
	alpha_3 := generate.Generate_alpha_ver(pi[1], cu, pi[0], c, sr1, sr2, q)
	cjxch := *new(big.Int).Mul(&pi[1], &pi[2])
	cjxch.Mod(&cjxch, &q)
	alpha_4 := generate.Generate_alpha_ver(ce, cjxch, h, se_dash, c, sr_dash, q)

	if alpha1.Cmp(&alpha_1) == 0 && alpha2.Cmp(&alpha_2) == 0 && alpha3.Cmp(&alpha_3) == 0 && alpha4.Cmp(&alpha_4) == 0 {
		return 1
	}

	return 0

}
