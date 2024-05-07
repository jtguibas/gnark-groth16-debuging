# Gnark Groth16 Verification Issue

## Get Proof

`go test -v -timeout 30s -run ^TestGroth16$ github.com/succinctlabs/sp1-recursion-gnark`

Then paste the proof into `contracts/src/VerifierGroth16.sol` and run `forge test --vvv`.# gnark-groth16-debuging
