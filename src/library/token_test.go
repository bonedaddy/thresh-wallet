// thresh-wallet
//
// Copyright 2019 by KeyFuse
//
// GPLv3 License

package library

import (
	"testing"

	"server"

	"github.com/stretchr/testify/assert"
)

func TestTokenAPIGetVCode(t *testing.T) {
	ts := server.MockServer()
	defer ts.Close()

	mobile := "10086"
	body := APIGetVCode(ts.URL, mobile)
	rsp := &VcodeResponse{}
	unmarshal(body, rsp)
	t.Logf("%+v", rsp)
	assert.Equal(t, 200, rsp.Code)
}

func TestTokenAPIGetToken(t *testing.T) {
	ts := server.MockServer()
	defer ts.Close()

	mobile := "10086"
	body := APIGetToken(ts.URL, mobile, "vcode", mockMasterPubKey)
	rsp := &TokenResponse{}
	unmarshal(body, rsp)
	t.Logf("%+v", rsp)
	assert.Equal(t, 200, rsp.Code)
}