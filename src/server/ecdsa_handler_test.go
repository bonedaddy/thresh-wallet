// thresh-wallet
//
// Copyright 2019 by KeyFuse
//
// GPLv3 License

package server

import (
	"testing"

	"proto"

	"github.com/tokublock/tokucore/xcore/bip32"
	"github.com/tokublock/tokucore/xcrypto"
	"github.com/tokublock/tokucore/xcrypto/secp256k1"

	"github.com/stretchr/testify/assert"
)

func TestEcdsaAddressHandler(t *testing.T) {
	ts := MockServer()
	defer ts.Close()

	// Token.
	{
		req := &proto.TokenRequest{
			UID:          mockUID,
			MasterPubKey: mockCliMasterPubKey,
		}
		httpRsp, err := proto.NewRequest().Post(ts.URL+"/api/token", req)
		assert.Nil(t, err)

		resp := &proto.TokenResponse{}
		httpRsp.Json(resp)
		t.Log(resp)
	}

	// New address.
	{
		req := &proto.EcdsaAddressRequest{}
		httpRsp, err := proto.NewRequest().SetHeaders("Authorization", mockToken).Post(ts.URL+"/api/ecdsa/newaddress", req)
		assert.Nil(t, err)
		assert.Equal(t, 200, httpRsp.StatusCode())
	}
}

func TestEcdsaR2S2Handler(t *testing.T) {
	var pos uint32
	var shareR *secp256k1.Scalar

	ts := MockServer()
	defer ts.Close()

	pos = 1
	hash := []byte{0x01, 0x02, 0x03, 0x04}

	// Client.
	climasterkey, err := bip32.NewHDKeyFromString(mockCliMasterPrvKey)
	assert.Nil(t, err)
	clichildkey, err := climasterkey.Derive(pos)
	assert.Nil(t, err)
	cliprv := clichildkey.PrivateKey()
	aliceParty := xcrypto.NewEcdsaParty(cliprv)

	// Phase2.
	encpk1, encpub1, r1 := aliceParty.Phase2(hash)
	_ = encpk1
	_ = encpub1

	// Token.
	{
		req := &proto.TokenRequest{
			UID:          mockUID,
			MasterPubKey: mockCliMasterPubKey,
		}

		_, err := proto.NewRequest().Post(ts.URL+"/api/token", req)
		assert.Nil(t, err)
	}

	// R2.
	{
		req := &proto.EcdsaR2Request{
			Pos:  pos,
			Hash: hash,
			R1:   r1,
		}
		httpRsp, err := proto.NewRequest().SetHeaders("Authorization", mockToken).Post(ts.URL+"/api/ecdsa/r2", req)
		assert.Nil(t, err)
		assert.Equal(t, 200, httpRsp.StatusCode())

		rsp := &proto.EcdsaR2Response{}
		err = httpRsp.Json(rsp)
		assert.Nil(t, err)
		shareR = rsp.ShareR
	}

	// S2.
	{
		req := &proto.EcdsaS2Request{
			Pos:     pos,
			Hash:    hash,
			R1:      r1,
			EncPK1:  encpk1,
			EncPub1: encpub1,
			ShareR:  shareR,
		}
		httpRsp, err := proto.NewRequest().SetHeaders("Authorization", mockToken).Post(ts.URL+"/api/ecdsa/s2", req)
		assert.Nil(t, err)
		assert.Equal(t, 200, httpRsp.StatusCode())
	}

	// S2 error.
	{
		req := &proto.EcdsaS2Request{
			Pos:     2,
			Hash:    hash,
			R1:      r1,
			EncPK1:  encpk1,
			EncPub1: encpub1,
			ShareR:  shareR,
		}
		httpRsp, err := proto.NewRequest().SetHeaders("Authorization", mockToken).Post(ts.URL+"/api/ecdsa/s2", req)
		assert.Nil(t, err)
		assert.Equal(t, 500, httpRsp.StatusCode())
	}
}