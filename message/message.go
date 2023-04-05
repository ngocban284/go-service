package message

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"time"

	"github.com/mihongtech/tendermint/libs/protoio"
	protoTypes "github.com/tendermint/tendermint/proto/tendermint/types"
	tmTypes "github.com/tendermint/tendermint/types"
)

func First[T, U any](val T, _ U) T {
	return val
}

// FromProto sets a protobuf PartSetHeader to the given pointer
func PartSetHeaderFromProto(ppsh *protoTypes.PartSetHeader) (*tmTypes.PartSetHeader, error) {
	if ppsh == nil {
		return nil, errors.New("nil PartSetHeader")
	}
	psh := new(tmTypes.PartSetHeader)
	psh.Total = ppsh.Total
	psh.Hash = ppsh.Hash

	return psh, psh.ValidateBasic()
}

func BlockIDFromProto(bID *protoTypes.BlockID) (*tmTypes.BlockID, error) {
	if bID == nil {
		return nil, errors.New("nil BlockID")
	}

	blockID := new(tmTypes.BlockID)
	ph, err := PartSetHeaderFromProto(&bID.PartSetHeader)
	if err != nil {
		return nil, err
	}

	blockID.PartSetHeader = *ph
	blockID.Hash = bID.Hash

	return blockID, blockID.ValidateBasic()
}

func CanonicalizeBlockID(bid protoTypes.BlockID) *protoTypes.CanonicalBlockID {
	rbid, err := BlockIDFromProto(&bid)
	if err != nil {
		panic(err)
	}
	var cbid *protoTypes.CanonicalBlockID
	if rbid == nil || rbid.IsZero() {
		cbid = nil
	} else {
		cbid = &protoTypes.CanonicalBlockID{
			Hash:          bid.Hash,
			PartSetHeader: tmTypes.CanonicalizePartSetHeader(bid.PartSetHeader),
		}
	}

	return cbid
}

func CanonicalizeVote(chainID string, vote *protoTypes.Vote) protoTypes.CanonicalVote {
	return protoTypes.CanonicalVote{
		Type:      vote.Type,
		Height:    vote.Height,       // encoded as sfixed64
		Round:     int64(vote.Round), // encoded as sfixed64
		BlockID:   CanonicalizeBlockID(vote.BlockID),
		Timestamp: vote.Timestamp,
		ChainID:   chainID,
	}
}

func VoteSignBytes(chainID string, vote *protoTypes.Vote) []byte {
	pb := CanonicalizeVote(chainID, vote)
	bz, err := protoio.MarshalDelimited(&pb)
	if err != nil {
		panic(err)
	}

	return bz
}

func Message() bool {

	vote := protoTypes.CanonicalVote{
		Type:   2,
		Height: 10320459,
		Round:  0,
		BlockID: &protoTypes.CanonicalBlockID{
			Hash: First(hex.DecodeString("D89A2762A9996953D0396D56478A7A4C4F4ADA8C0631756FCC17E2DD0DD5BB08")),
			PartSetHeader: protoTypes.CanonicalPartSetHeader{
				Total: 1,
				Hash:  First(hex.DecodeString("E987C5881C464D77416F0A52D811FB49F50E6BB592C2A63F921A0B679337A90E")),
			},
		},
		Timestamp: First(time.Parse(time.RFC3339, "2023-02-17T07:06:47.664674294Z")),
		ChainID:   "Oraichain",
	}

	// fmt.Println(vote,"\n")

	bz, err := protoio.MarshalDelimited(&vote)
	if err != nil {
		panic(err)
	}

	// // verify
	publicKey := ed25519.PublicKey(First(base64.StdEncoding.DecodeString("/ShOMJ4joYZBqPVFtD0+skU59lBh84uAyLkmeL6Dpwo=")))

	signature := First(base64.StdEncoding.DecodeString("Oyfq86rjqsiZMPQUWTpKxYm9Ovu/od/XoQksOdq0jw+ITd38m6hcEtU7PpxZ51/DV4CMqJ3uWmyU4rPlKZ9RCQ=="))

	result := ed25519.Verify(publicKey, bz, signature)

	return result

}
