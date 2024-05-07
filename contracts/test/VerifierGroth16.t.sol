// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import {Test, console} from "forge-std/Test.sol";
import {Verifier} from "../src/VerifierGroth16.sol";

contract CounterTest is Test {
    Verifier public verifier;

    function setUp() public {
        verifier = new Verifier();
    }

    function test_Groth16Verifier() public view {
        uint256[8] memory proofs = [
            0x012d00c1cf9761a3e93a61d1b71bd3e15cf9a7fc8b9b18535dd23f50ee4831a8,
            0x211715892e4275635d3853e81761ed90dc1076463b68980c276291efb64a9a89,
            0x17601eff3365916908af969ee77169ca70ccaee499e252d82b33afb676e44bc1,
            0x112447ca6a62acf7b7e5ed2d433842990601a671beb7e100c289be7e1597f25a,
            0x01d0ea30b64846951ab44416e611da201ec62820b045eadf1e53ce8dd1a39712,
            0x27765ebaa2f6e531003e0c37708f3d1a24e73ddc4c9eb8b0a005b498d10d1ab9,
            0x0d538214ca6455f9c9602fbe99612f9e651d7c0d86c8daa32ef472c232978442,
            0x09601ed80712947d1aa8d433578d75d768a0fc23aa13f224dc6e1d8c8e9095e5
        ];

        uint256[2] memory commitments = [
            5656590529146675058289059245332645040837777922060831408756977729008169812504,
            9859017459212054105361600295862626556918694722003898409496634160520790458221
        ];
        uint256[2] memory commitmentPok = [
            19182249199998087483757595626889009287415275414173380679261175107410690224245,
            16769003569912393248849000598660646262097617362318778297391562363009946087105
        ];

        uint256[3] memory inputs = [uint256(1), uint256(2), uint256(3)];
        verifier.verifyProof(proofs, commitments, commitmentPok, inputs);
    }
}
