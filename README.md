# Set Membership 
 
 This repository contains a bare-bones implementation of **CP-SNARK for Set Membership** in Golang. The current implementations is based on the following paper:
 * [Zero-Knowledge Proofs for Set Membership: Efficient, Succinct, Modular](https://eprint.iacr.org/2019/1255.pdf) by **Daniel Benarroch,Matteo Campanelli,Dario Fiore,Kobi Gurkan and Dimitris
Kolonelos**.

### CP SNARK for set membership

In this approach we are commiting neccessary values and then proving them using short ZK-NIZK(ZK-snarks). The sub-protocols are as follows:

 * **Setup** : A simple implementation of Pedersen commitment in Zq for commiting set elements and also corresponding prime elements. We also have a set commitment scheme using RSA accumulators which provides a very short commitment for the entire set.
 * **Root** : This sub-protocol proofs in zero knowledge that an element in the set is a member without revealing the element itself. This is done by calculating the witness,using pedersen commitments and schnorr's authentication
 * **ModEq** : A ZK proof that Ce(commitment of prime element e in RSA group) and ce(commitment of prime element e in prime order group) commit to the same value modulo q.
 * **HashEq** A ZK proof that for a given prime e, whose preimage is set member u there exist a j such that e = Hprime(u) = H(u,j) = F(u)+j .(The proof and verification for this sub-protocol unlike the others  was self-made and may be liability without further investigation )

 * The hash-to-prime used is a deterministic map given in the paper,it maps in a collision resistant manner backed by Cramer's conjecture.
 
 
 The Setmembership package contains the integration of all the above mentioned sub-protocols (Proving and Verifying).
check.go gives an example of the flow.
If you want to change the bounds please modify the hard coded values in the root file. depending to what are your likings.

### Improvement required
* Investigating the security of implemented  HashEq
* Efficiency and testing


