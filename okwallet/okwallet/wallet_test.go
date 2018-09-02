package okwallet

import (
	"testing"
)

func TestGetAddressByPrivateKey(t *testing.T) {
	const prvkey = `de47f09131a701b1012ab205df484ec43dd73daa6d42a0fead004a16307d7fdd`
	const expect_pub = `036f9c4bc2c75334da4c145848769c6128dd7edafa8059f28da5feb60e3d17e3df`
	const expect_address = `ShoPrXXG2nMATE3WDA7B7sgJscCoNXauaY`

	pub, addr, err := GetAddressByPrivateKey(prvkey)
	if err != nil {
		t.Fatal(err)
	}
	if expect_pub != pub {
		t.Fatal("GetAddressByPrivateKey failed. expect pubkey " + expect_pub + " but return " + pub)
	}
	if expect_address != addr {
		t.Fatal("GetAddressByPrivateKey failed. expect address " + expect_address + " but return " + addr)
	}
}

func TestCreateRawTransaction(t *testing.T) {
	const input = `
	{
"assetid":"23d4d5aeb154126332cc9aa5dd0a4ce3bb4cf3d2df507bd9a85ab6b7b9c9bbf4",
"pubkey":"036f9c4bc2c75334da4c145848769c6128dd7edafa8059f28da5feb60e3d17e3df",
"fee_per_target": "1",
"targets": [
  {
    "address":"SfgWjUW6H1NvFUgTBhrdxs6uf6upy2Tc5G",
    "value":"100"
  },
  {
    "address":"SbH2nd8p7veM2LG6uJsu6dDRnmww7B94ai",
    "value":"200"
  }
],
"unspent_info": [
  {
    "txid":"1da2fcb9a6487c7df825dd7324cacc4ac5b5c4d0e22d82594bc5365e0b7b9dcc",
    "index": 0,
    "value": "150"
  },
	{
    "txid":"1da2fcb9a6487c7df825dd7324cacc4ac5b5c4d0e22d82594bc5365e0b7b9dcc",
    "index": 1,
    "value": "200"
  }
],
"input_tx": [
  {
		"txid":"1da2fcb9a6487c7df825dd7324cacc4ac5b5c4d0e22d82594bc5365e0b7b9dcc",
		"outputs": [
			{
				"assetid":"23d4d5aeb154126332cc9aa5dd0a4ce3bb4cf3d2df507bd9a85ab6b7b9c9bbf4",
				"value":150,
				"programhash":"e0f1a129f50fc717fe746c5d041d0054874c8609"
			},
			{
				"assetid":"23d4d5aeb154126332cc9aa5dd0a4ce3bb4cf3d2df507bd9a85ab6b7b9c9bbf4",
				"value":200,
				"programhash":"e0f1a129f50fc717fe746c5d041d0054874c8609"
			}
		]
	}
]
}
	`

	const expect = `{"TxData":"800001001335353737303036373931393437373739343130021da2fcb9a6487c7df825dd7324cacc4ac5b5c4d0e22d82594bc5365e0b7b9dcc00001da2fcb9a6487c7df825dd7324cacc4ac5b5c4d0e22d82594bc5365e0b7b9dcc010003f4bbc9b9b7b65aa8d97b50dfd2f34cbbe34c0adda59acc32631254b1aed5d4230003164e02000000c9b43543291886dfa20961bdea57cf0f0f945dcdf4bbc9b9b7b65aa8d97b50dfd2f34cbbe34c0adda59acc32631254b1aed5d42300e721a2040000009962cf555c2b9b0a142cfcec1787b1b534863096f4bbc9b9b7b65aa8d97b50dfd2f34cbbe34c0adda59acc32631254b1aed5d42300f2052a01000000e0f1a129f50fc717fe746c5d041d0054874c860900","input_tx":[{"Txid":"1da2fcb9a6487c7df825dd7324cacc4ac5b5c4d0e22d82594bc5365e0b7b9dcc","outputs":[{"AssetID":"23d4d5aeb154126332cc9aa5dd0a4ce3bb4cf3d2df507bd9a85ab6b7b9c9bbf4","Value":150,"ProgramHash":"e0f1a129f50fc717fe746c5d041d0054874c8609"},{"AssetID":"23d4d5aeb154126332cc9aa5dd0a4ce3bb4cf3d2df507bd9a85ab6b7b9c9bbf4","Value":200,"ProgramHash":"e0f1a129f50fc717fe746c5d041d0054874c8609"}]}]}`

	rawTx, err := CreateRawTransaction(input)
	if err != nil {
		t.Fatal(err)
	}
	if expect != rawTx {
		t.Fatal("CreateRawTransaction failed. expect result is " + expect + " but return " + rawTx)
	}
}

func TestSignRawTransaction(t *testing.T) {
	const prvkey = `de47f09131a701b1012ab205df484ec43dd73daa6d42a0fead004a16307d7fdd`
	const input = `{"TxData":"800001001335353737303036373931393437373739343130021da2fcb9a6487c7df825dd7324cacc4ac5b5c4d0e22d82594bc5365e0b7b9dcc00001da2fcb9a6487c7df825dd7324cacc4ac5b5c4d0e22d82594bc5365e0b7b9dcc010003f4bbc9b9b7b65aa8d97b50dfd2f34cbbe34c0adda59acc32631254b1aed5d4230003164e02000000c9b43543291886dfa20961bdea57cf0f0f945dcdf4bbc9b9b7b65aa8d97b50dfd2f34cbbe34c0adda59acc32631254b1aed5d42300e721a2040000009962cf555c2b9b0a142cfcec1787b1b534863096f4bbc9b9b7b65aa8d97b50dfd2f34cbbe34c0adda59acc32631254b1aed5d42300f2052a01000000e0f1a129f50fc717fe746c5d041d0054874c860900","input_tx":[{"Txid":"1da2fcb9a6487c7df825dd7324cacc4ac5b5c4d0e22d82594bc5365e0b7b9dcc","outputs":[{"AssetID":"23d4d5aeb154126332cc9aa5dd0a4ce3bb4cf3d2df507bd9a85ab6b7b9c9bbf4","Value":150,"ProgramHash":"e0f1a129f50fc717fe746c5d041d0054874c8609"},{"AssetID":"23d4d5aeb154126332cc9aa5dd0a4ce3bb4cf3d2df507bd9a85ab6b7b9c9bbf4","Value":200,"ProgramHash":"e0f1a129f50fc717fe746c5d041d0054874c8609"}]}]}`
	//使用p256r1方法加密时使用一个随机数，因此每次加密出来的数据会不太一样
	//const expect = `800001001335353737303036373931393437373739343130021da2fcb9a6487c7df825dd7324cacc4ac5b5c4d0e22d82594bc5365e0b7b9dcc00001da2fcb9a6487c7df825dd7324cacc4ac5b5c4d0e22d82594bc5365e0b7b9dcc010003f4bbc9b9b7b65aa8d97b50dfd2f34cbbe34c0adda59acc32631254b1aed5d4230003164e02000000c9b43543291886dfa20961bdea57cf0f0f945dcdf4bbc9b9b7b65aa8d97b50dfd2f34cbbe34c0adda59acc32631254b1aed5d42300e721a2040000009962cf555c2b9b0a142cfcec1787b1b534863096f4bbc9b9b7b65aa8d97b50dfd2f34cbbe34c0adda59acc32631254b1aed5d42300f2052a01000000e0f1a129f50fc717fe746c5d041d0054874c86090141400ddf61d9123873b71c8a197c78c8eabb151a8bce709bd42356598f79d3bcfb0b866c42652156e1bfd0c15f536ecf7d5ed31c24ad9f5d9df78b1dfec8483b3c222321036f9c4bc2c75334da4c145848769c6128dd7edafa8059f28da5feb60e3d17e3dfac`

	_, err := SignRawTransaction(prvkey, input)
	if err != nil {
		t.Fatal(err)
	}
}
