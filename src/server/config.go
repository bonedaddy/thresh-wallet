// thresh-wallet
//
// Copyright 2019 by KeyFuse
//
// GPLv3 License

package server

import (
	"encoding/json"
	"io/ioutil"
)

// Config --
type Config struct {
	DataDir              string `json:"datadir"`
	ChainNet             string `json:"chainnet"`
	Endpoint             string `json:"endpoint"`
	TokenSecret          string `json:"token_secret"`
	SpvProvider          string `json:"spv_provider"`
	DisableVCode         bool   `json:"disable_vcode"`
	VCodeExpired         int    `json:"vcode_expired"`
	WalletSyncIntervalMs int    `json:"wallet_sync_interval_ms"`
}

// DefaultConfig -- returns default server config.
func DefaultConfig() *Config {
	return &Config{
		DataDir:              "./wallet",
		ChainNet:             "testnet",
		Endpoint:             ":9099",
		SpvProvider:          "blockstream",
		TokenSecret:          "thresh-wallet-demo-token-secret",
		VCodeExpired:         5 * 60,
		WalletSyncIntervalMs: 30 * 1000,
	}
}

// UnmarshalJSON -- built-in method for set default value when Unmarshal.
func (c *Config) UnmarshalJSON(b []byte) error {
	type confAlias *Config
	conf := confAlias(DefaultConfig())
	if err := json.Unmarshal(b, conf); err != nil {
		return err
	}
	*c = Config(*conf)
	return nil
}

// LoadConfig -- used to load the config from file.
func LoadConfig(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	conf := &Config{}
	if err := json.Unmarshal([]byte(data), conf); err != nil {
		return nil, err
	}
	return conf, nil
}