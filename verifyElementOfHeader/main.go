package main

import (
	// it is ok to use math/rand here: we do not need a cryptographically secure random
	// number generator here and we can run the tests a bit faster

	"encoding/hex"
	"fmt"
	"reflect"
	"time"

	gogotypes "github.com/gogo/protobuf/types"

	"github.com/tendermint/tendermint/crypto/merkle"
	"github.com/tendermint/tendermint/crypto/tmhash"
	"github.com/tendermint/tendermint/libs/bytes"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
	tmversion "github.com/tendermint/tendermint/proto/tendermint/version"
	"github.com/tendermint/tendermint/types"
	"github.com/tendermint/tendermint/version"
)

func First[T, U any](val T, _ U) T {
	return val
}

func isTypedNil(o interface{}) bool {
	rv := reflect.ValueOf(o)
	switch rv.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.Slice:
		return rv.IsNil()
	default:
		return false
	}
}

func isEmpty(o interface{}) bool {
	rv := reflect.ValueOf(o)
	switch rv.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
		return rv.Len() == 0
	default:
		return false
	}
}

func cdcEncode(item interface{}) []byte {
	if item != nil && !isTypedNil(item) && !isEmpty(item) {
		switch item := item.(type) {
		case string:
			i := gogotypes.StringValue{
				Value: item,
			}
			bz, err := i.Marshal()
			if err != nil {
				return nil
			}
			return bz
		case int64:
			i := gogotypes.Int64Value{
				Value: item,
			}
			bz, err := i.Marshal()
			if err != nil {
				return nil
			}
			return bz
		case bytes.HexBytes:
			i := gogotypes.BytesValue{
				Value: item,
			}
			bz, err := i.Marshal()
			if err != nil {
				return nil
			}
			return bz
		default:
			return nil
		}
	}

	return nil
}

func makeBlockID(hash []byte, partSetSize uint32, partSetHash []byte) types.BlockID {
	var (
		h   = make([]byte, tmhash.Size)
		psH = make([]byte, tmhash.Size)
	)
	copy(h, hash)
	copy(psH, partSetHash)
	return types.BlockID{
		Hash: h,
		PartSetHeader: types.PartSetHeader{
			Total: partSetSize,
			Hash:  psH,
		},
	}
}

func hexBytesFromString(s string) bytes.HexBytes {
	b, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return bytes.HexBytes(b)
}

// Hash returns the hash of the header.
// It computes a Merkle tree from the header fields
// ordered as they appear in the Header.
// Returns nil if ValidatorHash is missing,
// since a Header is not valid unless there is
// a ValidatorsHash (corresponding to the validator set).
// func (h *Header) Hash() tmbytes.HexBytes {
// 	if h == nil || len(h.ValidatorsHash) == 0 {
// 		return nil
// 	}
// 	hpb := h.Version.ToProto()
// 	hbz, err := hpb.Marshal()
// 	if err != nil {
// 		return nil
// 	}

// 	pbt, err := gogotypes.StdTimeMarshal(h.Time)
// 	if err != nil {
// 		return nil
// 	}

// 	pbbi := h.LastBlockID.ToProto()
// 	bzbi, err := pbbi.Marshal()
// 	if err != nil {
// 		return nil
// 	}
// 	return merkle.HashFromByteSlices([][]byte{
// 		hbz,
// 		cdcEncode(h.ChainID),
// 		cdcEncode(h.Height),
// 		pbt,
// 		bzbi,
// 		cdcEncode(h.LastCommitHash),
// 		cdcEncode(h.DataHash),
// 		cdcEncode(h.ValidatorsHash),
// 		cdcEncode(h.NextValidatorsHash),
// 		cdcEncode(h.ConsensusHash),
// 		cdcEncode(h.AppHash),
// 		cdcEncode(h.LastResultsHash),
// 		cdcEncode(h.EvidenceHash),
// 		cdcEncode(h.ProposerAddress),
// 	})
// }

func main() {
	blockHeader := &types.Header{
		Version: version.Consensus{
			Block: 11,
		},
		ChainID: "Oraichain",
		Height:  10340037,
		Time:    First(time.Parse(time.RFC3339, "2023-02-18T17:07:42.760101663Z")),
		LastBlockID: types.BlockID{
			Hash: tmbytes.HexBytes("73CDB4036015959B90C1E0D422CCA5BAC94E74FBD92A3BB435FBE35781B33917"),
			PartSetHeader: types.PartSetHeader{
				Total: 1,
				Hash:  tmbytes.HexBytes("02929D90D6E40AF4913FDD97860A53B0348DB1F62A4AC00751C3C8BF6BD64960"),
			},
		},
		LastCommitHash:     tmbytes.HexBytes("14BDEB8BA16902C0CA1035D592ED964FBF73DD8F33CEF21288E786ADB7C5A0F8"),
		DataHash:           tmbytes.HexBytes("677BF175DE9C1EDDD2F26AE4161631390A24486D44BBC71982C39965F58967C4"),
		ValidatorsHash:     tmbytes.HexBytes("1A695B879702E2CBA64500C4717D9A96C951ED2083124F1179B7E7223825EA6D"),
		NextValidatorsHash: tmbytes.HexBytes("1A695B879702E2CBA64500C4717D9A96C951ED2083124F1179B7E7223825EA6D"),
		ConsensusHash:      tmbytes.HexBytes("048091BC7DDC283F77BFBF91D73C44DA58C3DF8A9CBC867405D8B7F3DAADA22F"),
		AppHash:            tmbytes.HexBytes("E2BA58BAE0A12D24920774237D0B2FB97CC4678369FB5CAB90FBCAB79F33F244"),
		LastResultsHash:    tmbytes.HexBytes("E3B0C44298FC1C149AFBF4C8996FB92427AE41E4649B934CA495991B7852B855"),
		EvidenceHash:       tmbytes.HexBytes("E3B0C44298FC1C149AFBF4C8996FB92427AE41E4649B934CA495991B7852B855"),
		ProposerAddress:    tmbytes.HexBytes("0BDC699EF20C95A99B746A8F7F18D35E2AAF0C3D"),
	}

	blockHash := blockHeader.Hash()
	fmt.Println("expected hash: ", blockHash)

	// root ,proofs := merkle.HashFromByteSlices()
	version := tmversion.Consensus{Block: 11}
	versionBytes, err := version.Marshal()
	if err != nil {
		fmt.Println("error: ", err)
	}

	time := First(time.Parse(time.RFC3339, "2023-02-18T17:07:42.760101663Z"))
	timeBytes, err := gogotypes.StdTimeMarshal(time)
	if err != nil {
		fmt.Println("error: ", err)
	}

	blockId := types.BlockID{
		Hash: tmbytes.HexBytes("73CDB4036015959B90C1E0D422CCA5BAC94E74FBD92A3BB435FBE35781B33917"),
		PartSetHeader: types.PartSetHeader{
			Total: 1,
			Hash:  tmbytes.HexBytes("02929D90D6E40AF4913FDD97860A53B0348DB1F62A4AC00751C3C8BF6BD64960"),
		},
	}
	blockProto := blockId.ToProto()
	blockBytes, err := blockProto.Marshal()
	if err != nil {
		fmt.Println("error: ", err)
	}

	newHeader := [][]byte{
		versionBytes,
		cdcEncode("Oraichain"),
		cdcEncode(int64(10340037)),
		timeBytes,
		blockBytes,
		cdcEncode(tmbytes.HexBytes("14BDEB8BA16902C0CA1035D592ED964FBF73DD8F33CEF21288E786ADB7C5A0F8")),
		cdcEncode(tmbytes.HexBytes("677BF175DE9C1EDDD2F26AE4161631390A24486D44BBC71982C39965F58967C4")),
		cdcEncode(tmbytes.HexBytes("1A695B879702E2CBA64500C4717D9A96C951ED2083124F1179B7E7223825EA6D")),
		cdcEncode(tmbytes.HexBytes("1A695B879702E2CBA64500C4717D9A96C951ED2083124F1179B7E7223825EA6D")),
		cdcEncode(tmbytes.HexBytes("048091BC7DDC283F77BFBF91D73C44DA58C3DF8A9CBC867405D8B7F3DAADA22F")),
		cdcEncode(tmbytes.HexBytes("E2BA58BAE0A12D24920774237D0B2FB97CC4678369FB5CAB90FBCAB79F33F244")),
		cdcEncode(tmbytes.HexBytes("E3B0C44298FC1C149AFBF4C8996FB92427AE41E4649B934CA495991B7852B855")),
		cdcEncode(tmbytes.HexBytes("E3B0C44298FC1C149AFBF4C8996FB92427AE41E4649B934CA495991B7852B855")),
		cdcEncode(tmbytes.HexBytes("0BDC699EF20C95A99B746A8F7F18D35E2AAF0C3D")),
	}

	// fmt.Println("newHeader: ", newHeader)

	root, proofs := merkle.ProofsFromByteSlices(newHeader)
	fmt.Println("proofs: ", proofs)
	fmt.Printf("root: %x\n", root)

	fmt.Println("Aunts: ", proofs[6])

	fmt.Println("Leaf: ", cdcEncode(tmbytes.HexBytes("677BF175DE9C1EDDD2F26AE4161631390A24486D44BBC71982C39965F58967C4")))

	fmt.Printf("LeafHash: %x\n", proofs[6].LeafHash)

	fmt.Println("Index: ", proofs[6].Index)

	fmt.Println("Total: ", proofs[6].Total)

	err = proofs[6].Verify(root, cdcEncode(tmbytes.HexBytes("677BF175DE9C1EDDD2F26AE4161631390A24486D44BBC71982C39965F58967C4")))
	fmt.Println("err: ", err)

}
