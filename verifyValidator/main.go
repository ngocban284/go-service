package main

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"fmt"

	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/crypto/merkle"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	pc "github.com/tendermint/tendermint/proto/tendermint/crypto"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	"github.com/tendermint/tendermint/types"
)

func First[T, U any](val T, _ U) T {
	return val
}

// PubKeyToProto takes crypto.PubKey and transforms it to a protobuf Pubkey
func PubKeyToProto(k crypto.PubKey) (pc.PublicKey, error) {
	var kp pc.PublicKey
	switch k := k.(type) {
	case ed25519.PubKey:
		kp = pc.PublicKey{
			Sum: &pc.PublicKey_Ed25519{
				Ed25519: k,
			},
		}
	case secp256k1.PubKey:
		kp = pc.PublicKey{
			Sum: &pc.PublicKey_Secp256K1{
				Secp256K1: k,
			},
		}
	default:
		return kp, fmt.Errorf("toproto: key type %v is not supported", k)
	}
	return kp, nil
}

func Bytes(pubKey string, votingPower int64) []byte {
	pk, err := PubKeyToProto(ed25519.PubKey(First(base64.StdEncoding.DecodeString(pubKey))))
	if err != nil {
		panic(err)
	}

	pbv := tmproto.SimpleValidator{
		PubKey:      &pk,
		VotingPower: votingPower,
	}

	bz, err := pbv.Marshal()
	if err != nil {
		panic(err)
	}
	return bz
}

func main() {

	vals := types.ValidatorSet{
		Validators: []*types.Validator{
			{
				Address:          First(hex.DecodeString("oraivalcons145dpmvvaj5dzm00hzrtnwyv89kz704gdk4xrl3")),
				PubKey:           ed25519.PubKey(First(base64.StdEncoding.DecodeString("jLunee++7+9tO0vVIBG59POGwkbShGiOWbtggTZMjMM="))),
				VotingPower:      200,
				ProposerPriority: 84,
			},
			{
				Address:          First(hex.DecodeString("oraivalcons1r5e75xg7zu2pax6wu4q8u4ql5m8f2t3xe27c6p")),
				PubKey:           ed25519.PubKey(First(base64.StdEncoding.DecodeString("l9PY/oC7El5N7BmHIhn2Rw1n+BBSKxwPKAjnh6JCQbc="))),
				VotingPower:      2,
				ProposerPriority: -12,
			},
			{
				Address:          First(hex.DecodeString("oraivalcons1rergngrr8w30gus3tsmgdxsardksgjjkzs2kgh")),
				PubKey:           ed25519.PubKey(First(base64.StdEncoding.DecodeString("w8cGC01n/3SDiUgCTq8aFQgAp5lsjlOqOIhsm/s4jOs="))),
				VotingPower:      2,
				ProposerPriority: -12,
			},
			{
				Address:          First(hex.DecodeString("oraivalcons1w6jkamejxccjmu5ys63zu8n9ttwlcxv9wsms6g")),
				PubKey:           ed25519.PubKey(First(base64.StdEncoding.DecodeString("iHea1XlBnUpjaE5wDBBD9XJ+I9lQj7YUMr1wJYxiSpk="))),
				VotingPower:      2,
				ProposerPriority: -12,
			},
			{
				Address:          First(hex.DecodeString("oraivalcons1339lws5730vztp7hx0kkcrfhnqfs7gws69zg6m")),
				PubKey:           ed25519.PubKey(First(base64.StdEncoding.DecodeString("Fnu5TVF9wY/Z3lHlt0rTZ6Q6NCnNToKnspzMdEGEJO8="))),
				VotingPower:      2,
				ProposerPriority: -12,
			},
			{
				Address:          First(hex.DecodeString("oraivalcons15jpgu6dauyy9v9eje2nqcnx3y556w62d7jlhlu")),
				PubKey:           ed25519.PubKey(First(base64.StdEncoding.DecodeString("Og7coQxbSm4cMWbgpLqJrNPTbZi3TZBqr2CT9rPrz+E="))),
				VotingPower:      2,
				ProposerPriority: -12,
			},
			{
				Address:          First(hex.DecodeString("oraivalcons1h404jv5mkcdk9fag22l0retun9pc6ryxmspf9m")),
				PubKey:           ed25519.PubKey(First(base64.StdEncoding.DecodeString("zYJafIuidhsS9dIkl0u1empVdoTKShG3LvIG18fwif4="))),
				VotingPower:      2,
				ProposerPriority: -12,
			},
			{
				Address:          First(hex.DecodeString("oraivalcons1uaadxrp3hw834u78uq877p780a99qgdm5lzyqc")),
				PubKey:           ed25519.PubKey(First(base64.StdEncoding.DecodeString("0uDyc0WNK6VW98XHCYCXgyetK863YIyP31pikPp8jiU="))),
				VotingPower:      2,
				ProposerPriority: -12,
			},
		},
	}

	// fmt.Printf("vals : %x", vals.Hash())

	root := vals.Hash()
	fmt.Printf("expected root : %x\n", root)

	// // test := tmproto.SimpleValidator{
	// // 	PubKey:      ed25519.PubKey(First(base64.StdEncoding.DecodeString("jLunee++7+9tO0vVIBG59POGwkbShGiOWbtggTZMjMM="))),
	// // 	VotingPower: 200,
	// // }
	// pubTest, er := PubKeyToProto(ed25519.PubKey(First(base64.StdEncoding.DecodeString("jLunee++7+9tO0vVIBG59POGwkbShGiOWbtggTZMjMM="))))

	// if er != nil {
	// 	fmt.Println("error")
	// }

	// fmt.Printf("pubTest : %x \n", pubTest)

	valBytes := [][]byte{
		Bytes("jLunee++7+9tO0vVIBG59POGwkbShGiOWbtggTZMjMM=", 200),
		Bytes("l9PY/oC7El5N7BmHIhn2Rw1n+BBSKxwPKAjnh6JCQbc=", 2),
		Bytes("w8cGC01n/3SDiUgCTq8aFQgAp5lsjlOqOIhsm/s4jOs=", 2),
		Bytes("iHea1XlBnUpjaE5wDBBD9XJ+I9lQj7YUMr1wJYxiSpk=", 2),
		Bytes("Fnu5TVF9wY/Z3lHlt0rTZ6Q6NCnNToKnspzMdEGEJO8=", 2),
		Bytes("Og7coQxbSm4cMWbgpLqJrNPTbZi3TZBqr2CT9rPrz+E=", 2),
		Bytes("zYJafIuidhsS9dIkl0u1empVdoTKShG3LvIG18fwif4=", 2),
		Bytes("0uDyc0WNK6VW98XHCYCXgyetK863YIyP31pikPp8jiU=", 2),
	}

	newRoot, _ := merkle.ProofsFromByteSlices(valBytes)
	fmt.Printf("newRoot : %x \n", newRoot)

	if bytes.Equal(root, newRoot) {
		fmt.Println("validator hash is ok")
	}

	// fmt.Printf("proofs : %x \n", proofs[0].Proof)
}
