package okwallet

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"log"
	"os"

	"starchain/account"
	"starchain/common/config"
	"starchain/crypto"
)

func init() {
	file := []byte(`{
  "Configuration": {
    "MagicCode": 5841263,
    "Ver":1,
    "VerifyList": [
      "47.52.44.156:25887",
      "47.75.44.148:25887",
      "47.75.45.95:25887",
      "47.91.218.249:25887",
      "47.75.4.61:25887",
      "47.91.208.36:25887",
      "47.75.44.103:25887"
    ],
    "BookKeepers": [
      "02375c72e9ae42b90df8c18ed77306fc60dbc6a95be327e04c793d8397c2372b18",
      "0269c4016c2b57ca34ab85e6ba94051639fcb58cfb1490c013252a329dad2280a7",
      "03596675a9e1f00bbd9b8dfc21d04ff5bcc5729e71cf4b215fd7967ea28414994e",
      "031b68157c7211441ca42549038b5213f3533cb8f97f45d30f887d45aeeda2ea06",
      "03c4e4fc261d56ce8360ab854e4ca219d4dc7fc3a18d87eff84dc2a48108827e70",
      "0315b69a1d7249e03641a89cb7f0820e9301b6d19fe16f4b1e5d29cc914c72af71",
      "02f2e79b69c7757c49974d5941165737a41755baf6a4bedcd5a01dae5976015c69"
    ],
    "RestPort": 25884,
    "RestStart":true,
    "JsonPort": 25886,
    "NodePort": 25887,
    "NodeType": "service",
    "Tls": false,
    "MultiCoreNum": 7,
    "AllowIp":"127.0.0.1",
    "AppKey":"appkey",
    "SecretKey":"123456",
    "ChainPath":"",
    "WalletPath":"./config/walle.dat",
    "LogLevel":"info",
    "TransactionFee": {
      "Transfer": 0.000
    }
  }
}
`)
	file = bytes.TrimPrefix(file, []byte("\xef\xbb\xbf"))
	cfg := config.ConfigFile{}
	err := json.Unmarshal(file, &cfg)
	if err != nil {
		log.Fatalf("unmarshal config file error %v", err)
		os.Exit(1)
	}
	config.Parameters = &(cfg.ConfigFile)
}

func init() {
	crypto.SetAlg(config.Parameters.EncryptAlg)
}

func GetAddressByPrivateKey(prvkeyStr string) (string, error) {
	prvkey, err := hex.DecodeString(prvkeyStr)
	if err != nil {
		return "", err
	}

	account, err := account.NewAccountWithPrivatekey(prvkey)
	if err != nil {
		return "", err
	}

	return account.ProgramHash.ToAddress()
}
