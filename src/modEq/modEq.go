package modEq

import (
  "math/big"
  "crypto/rand"
  "crypto/sha256"
  //"math/rand"
  //"big"
  "fmt"
)
//q is taken as prime for group Zq
func Prove(N, G, H, q, g, h, C_e, c_eq, e, e_q, r, r_q *big.Int, lambda int) (alpha1, alpha2, s_e, s_r, s_r_dash *big.Int)  {

  lambda_z  := lambda - 3
  lambda_s  := lambda -2
  meu       := 2*lambda -2


  var temp1e, temp3e , r_e *big.Int
  temp1e.Exp(big.NewInt(2), big.NewInt(int64(lambda_z + lambda_s + meu +1 )), nil)
  temp2e, _ := rand.Int(rand.Reader, temp1e)
  temp3e.Mul(temp1e, big.NewInt(-1))
  r_e.Add(temp2e, temp3e)


  var temp1r, temp3r, r_r *big.Int
  temp1r.Exp(big.NewInt(2), big.NewInt(int64(lambda_z + lambda_s +1)), nil)
  N.Div(N, big.NewInt(4))
  temp1r.Mul(N, temp1r)
  temp2r, _ := rand.Int(rand.Reader, temp1r)
  temp3r.Mul(temp2r, big.NewInt(-1))
  r_r.Add(temp2r, temp3r)


  r_r_dash, _ := rand.Int(rand.Reader, q)


  var temp1a1, temp2a1 *big.Int
  temp1a1.Exp(G, r_e, N)
  temp2a1.Exp(H, r_r, N)
  alpha1.Mul(temp1a1, temp2a1)
  alpha1.Mod(alpha1, N)


  var temp1a2, temp2a2, r_e_mod *big.Int
  r_e_mod.Mod(r_e, q)
  temp1a2.Exp(g, r_e_mod, q)
  temp2a2.Exp(h, r_r_dash, q)
  alpha2.Mul(temp1a2, temp2a2)
  alpha2.Mod(alpha1, N)


  //c := big.NewInt(12345)
  x := []*big.Int {alpha1, alpha2, C_e, c_eq}
  hashh := sha256.New()

  for _,y := range x {
    hashh.Write(y.Bytes())
  }
  hh := fmt.Sprintf("%x", hashh.Sum(nil))

  c, _ := new(big.Int).SetString(hh, 16)

  s_e.Mul(c,e)
  s_e.Sub(r_e, s_e)


  s_r.Mul(c,r)
  s_r.Sub(r_r, s_r)

  s_r_dash.Mul(c,r_q)
  s_r_dash.Sub(r_r_dash, s_r_dash)


  return alpha1, alpha2, s_e, s_r, s_r_dash

}



func VerProof(N, G, H, g, h, q, C_e, c_eq, alpha1, alpha2, s_e, s_r, s_r_dash, c *big.Int) int {

  var temp1a1, temp2a1, temp3a1 *big.Int

  temp1a1.Exp(C_e, c, N)
  temp2a1.Exp(G, s_e, N)
  temp3a1.Exp(H, s_r, N)

  temp1a1.Mul(temp1a1, temp2a1)
  temp1a1.Mod(temp1a1, N)
  temp1a1.Mul(temp1a1, temp3a1)
  temp1a1.Mod(temp1a1, N)

  var temp1a2, temp2a2, temp3a2 *big.Int

  temp1a2.Exp(c_eq, c, q)
  s_e.Mod(s_e, q)
  temp2a2.Exp(g, s_e, q)
  temp3a2.Exp(h, s_r_dash, q)
  temp1a2.Mul(temp1a2, temp2a2)
  temp1a2.Mod(temp1a2, q)
  temp1a2.Mul(temp1a2, temp3a2)
  temp1a2.Mod(temp1a2, q)

  if alpha1.Cmp(temp1a1) == 0 && alpha2.Cmp(temp1a2) == 0 {
    return 1
  }

  return 0


}
