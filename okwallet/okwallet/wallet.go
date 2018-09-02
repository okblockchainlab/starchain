package okwallet

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"os"
	"strconv"

	"starchain/account"
	"starchain/common"
	"starchain/common/config"
	"starchain/core/contract"
	"starchain/core/signature"
	"starchain/core/transaction"
	"starchain/crypto"
	stc_errors "starchain/errors"
	"starchain/util"
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
}`)

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

func GetAddressByPrivateKey(prvkeyStr string) (string, string, error) {
	prvkey, err := hex.DecodeString(prvkeyStr)
	if err != nil {
		return "", "", err
	}

	account, err := account.NewAccountWithPrivatekey(prvkey)
	if err != nil {
		return "", "", err
	}
	pubkey, err := account.PubKey().EncodePoint(true)
	if err != nil {
		return "", "", err
	}

	addrStr, err := account.ProgramHash.ToAddress()

	return common.BytesToHexString(pubkey), addrStr, err
}

type UTXOUnspentInfo struct {
	Txid  string
	Index uint32
	Value string
}

type OKTXOutput struct {
	AssetID     string
	Value       int64
	ProgramHash string
}

type InputTransaction struct {
	Txid           string
	UnSpendOutputs []OKTXOutput `json:"outputs"`
}

type RawTxInfo struct {
	AssetID      string
	Pubkey       string
	FeePerTarget string `json:"fee_per_target"`
	Targets      []util.BatchOut
	UnspentInfo  []UTXOUnspentInfo  `json:"unspent_info"`
	InputTx      []InputTransaction `json:"input_tx"`
}

type RawTransaction struct {
	TxData  string
	InputTx []InputTransaction `json:"input_tx"`
}

//OKMockLedgerStore is used for
//transaction.TxStore.GetTransaction called in transaction.GetReference
type OKMockLedgerStore struct {
	InputTx map[common.Uint256][]*transaction.TxOutput
}

func (l *OKMockLedgerStore) GetTransaction(hash common.Uint256) (*transaction.Transaction, error) {
	o, ok := l.InputTx[hash]
	if !ok {
		return nil, errors.New("can not find transaction")
	}

	return &transaction.Transaction{Outputs: o}, nil
}

func (l *OKMockLedgerStore) GetQuantityIssued(AssetId common.Uint256) (common.Fixed64, error) {
	return common.Fixed64(0), errors.New("not implemented")
}

func makeOKMockLedgerStore(InputTx []InputTransaction) (*OKMockLedgerStore, error) {
	result := OKMockLedgerStore{make(map[common.Uint256][]*transaction.TxOutput)}

	for _, it := range InputTx {
		txidBytes, err := hex.DecodeString(it.Txid)
		if err != nil {
			return nil, err
		}
		refID, err := common.Uint256ParseFromBytes(txidBytes)
		if err != nil {
			return nil, err
		}

		var outputs []*transaction.TxOutput
		for _, uo := range it.UnSpendOutputs {
			assetid, err := string2AssetID(uo.AssetID)
			if err != nil {
				return nil, err
			}

			value := common.Fixed64(uo.Value)

			phashBytes, err := hex.DecodeString(uo.ProgramHash)
			if err != nil {
				return nil, err
			}
			programhash, err := common.Uint160ParseFromBytes(phashBytes)
			if err != nil {
				return nil, err
			}

			outputs = append(outputs, &transaction.TxOutput{AssetID: assetid, Value: value, ProgramHash: programhash})
		}

		result.InputTx[refID] = outputs
	}

	return &result, nil
}

//OKSigner is used for signature.SignBySigner
type OKSigner []byte

func (s OKSigner) PrivKey() []byte {
	return []byte(s)
}

func (s OKSigner) PubKey() *crypto.PubKey {
	return crypto.NewPubKey([]byte(s))
}

func string2AssetID(s string) (common.Uint256, error) {
	tmp, err := common.HexStringToBytesReverse(s)
	if err != nil {
		return common.Uint256{}, err
	}

	var assetID common.Uint256
	if err := assetID.Deserialize(bytes.NewReader(tmp)); err != nil {
		return common.Uint256{}, err
	}

	return assetID, nil
}

func programHashFromPubkey(pubkeyStr string) (common.Uint160, error) {
	pubkeyBytes, err := hex.DecodeString(pubkeyStr)
	if err != nil {
		return common.Uint160{}, err
	}

	pubkey, err := crypto.DecodePoint(pubkeyBytes)
	if err != nil {
		return common.Uint160{}, err
	}

	signatureRedeemScript, err := contract.CreateSignatureRedeemScript(pubkey)
	if err != nil {
		return common.Uint160{}, stc_errors.NewDetailErr(err, stc_errors.ErrNoCode, "CreateSignatureRedeemScript failed")
	}
	programHash, err := common.ToCodeHash(signatureRedeemScript)
	if err != nil {
		return common.Uint160{}, stc_errors.NewDetailErr(err, stc_errors.ErrNoCode, "ToCodeHash failed")
	}

	return programHash, nil
}

func makeTransferTransaction(pubkey string, assetID common.Uint256, unspentInfo []UTXOUnspentInfo, feePerTargetStr string, batchOut ...util.BatchOut) (*transaction.Transaction, error) {
	//TODO: check if being transferred asset is System Token(STC)
	outputNum := len(batchOut)
	if outputNum == 0 {
		return nil, errors.New("nil outputs")
	}

	// get main account which is used to receive changes
	programHash, err := programHashFromPubkey(pubkey)
	if err != nil {
		return nil, err
	}

	feePerTarget, err := common.StringToFixed64(feePerTargetStr)
	if err != nil {
		return nil, err
	}

	var expected common.Fixed64
	input := []*transaction.UTXOTxInput{}
	output := []*transaction.TxOutput{}
	// construct transaction outputs
	for _, o := range batchOut {
		outputValue, err := common.StringToFixed64(o.Value)
		if err != nil {
			return nil, err
		}
		if outputValue <= feePerTarget {
			return nil, errors.New("token is not enough for transaction fee")
		}
		expected += outputValue
		address, err := common.ToScriptHash(o.Address)
		if err != nil {
			return nil, errors.New("invalid address")
		}
		tmp := &transaction.TxOutput{
			AssetID:     assetID,
			Value:       outputValue - feePerTarget,
			ProgramHash: address,
		}
		output = append(output, tmp)
	}

	// construct transaction inputs and changes
	for _, info := range unspentInfo {
		value, err := common.StringToFixed64(info.Value)
		if err != nil {
			return nil, err
		}

		txidBytes, err := hex.DecodeString(info.Txid)
		if err != nil {
			return nil, err
		}
		refID, err := common.Uint256ParseFromBytes(txidBytes)
		if err != nil {
			return nil, err
		}
		inp := &transaction.UTXOTxInput{
			ReferTxID:          refID,
			ReferTxOutputIndex: uint16(info.Index),
		}

		input = append(input, inp)

		if value > expected {
			changes := &transaction.TxOutput{
				AssetID:     assetID,
				Value:       value - expected,
				ProgramHash: programHash,
			}
			// if any, the changes output of transaction will be the last one
			output = append(output, changes)
			expected = 0
			break
		} else if value == expected {
			expected = 0
			break
		} else if value < expected {
			expected = expected - value
		}
	}

	if expected > 0 {
		return nil, errors.New("token is not enough")
	}

	// construct transaction
	txn, err := transaction.NewTransferAssetTransaction(input, output)
	if err != nil {
		return nil, err
	}
	txAttr := transaction.NewTxAttribute(transaction.Nonce, []byte(strconv.FormatInt(rand.Int63(), 10)))
	txn.Attributes = make([]*transaction.TxAttribute, 0)
	txn.Attributes = append(txn.Attributes, &txAttr)

	return txn, nil
}

func sign(signer signature.Signer, context *contract.ContractContext) error {

	sig_contract, err := contract.CreateSignatureContract(signer.PubKey())
	if err != nil {
		return err
	}

	for _, hash := range context.ProgramHashes {
		if hash != sig_contract.ProgramHash {
			return errors.New("ProgramHash is invalid")
		}

		switch {
		case sig_contract.IsStandard():
			sig, err := signature.SignBySigner(context.Data, signer)
			if err != nil {
				return err
			}
			if err := context.AddContract(sig_contract, signer.PubKey(), sig); err != nil {
				return err
			}
		default:
			return errors.New("unknown contract type")
		}
	}

	return nil
}

//reference: SendToManyAddress
func CreateRawTransaction(infoStr string) (string, error) {
	var createInfo RawTxInfo
	if err := json.Unmarshal([]byte(infoStr), &createInfo); err != nil {
		return "", err
	}

	assetID, err := string2AssetID(createInfo.AssetID)
	if err != nil {
		return "", err
	}

	txn, err := makeTransferTransaction(createInfo.Pubkey, assetID, createInfo.UnspentInfo, createInfo.FeePerTarget, createInfo.Targets...)
	if err != nil {
		return "", err
	}

	var buffer bytes.Buffer
	txn.Serialize(&buffer)
	txData := common.BytesToHexString(buffer.Bytes())

	rawtx := RawTransaction{
		TxData:  txData,
		InputTx: createInfo.InputTx,
	}

	result, err := json.Marshal(&rawtx)
	return string(result), err
}

//reference: begin of sendRawTransaction, and end of MakeTransferTransaction
func SignRawTransaction(prvkeyStr, rawTxStr string) (string, error) {
	prvkey, err := hex.DecodeString(prvkeyStr)
	if err != nil {
		return "", err
	}

	var rawtx RawTransaction
	if err := json.Unmarshal([]byte(rawTxStr), &rawtx); err != nil {
		return "", err
	}

	//copy from sendRawTransaction
	rawtxn, _ := common.HexStringToBytes(rawtx.TxData)
	var txn transaction.Transaction
	if err := txn.Deserialize(bytes.NewReader(rawtxn)); err != nil {
		return "", errors.New("invalid raw transaction")
	}

	signer := OKSigner(prvkey)

	oldTxStore := transaction.TxStore
	defer func() { transaction.TxStore = oldTxStore }()
	transaction.TxStore, err = makeOKMockLedgerStore(rawtx.InputTx)
	if err != nil {
		return "", err
	}

	ctx := contract.NewContractContext(&txn)
	if err := sign(signer, ctx); err != nil {
		return "", err
	}
	txn.SetPrograms(ctx.GetPrograms())

	var buffer bytes.Buffer
	txn.Serialize(&buffer)
	return common.BytesToHexString(buffer.Bytes()), nil
}
