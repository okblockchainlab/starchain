package main

// #include <jni.h>
import "C"

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"starchain/okwallet/okwallet"
)

const (
	GET_ADDRESS_BY_PRIVATE_KEY_CMD = "getaddressbyprivatekey"
	CREATE_RAW_TRANSACTION_CMD     = "createrawtransaction"
	SIGN_RAW_TRANSACTION           = "signrawtransaction"
)

type ResponeResult struct {
	Result  string `json:"result"`
	Message string `json:"message"`
}

func setResult(env *C.JNIEnv, result, msg string) C.jobjectArray {
	rr := ResponeResult{Result: result, Message: msg}
	res, _ := json.Marshal(rr)

	objArr := newStringObjectArray(env, 1)
	setObjectArrayStringElement(env, objArr, 0, string(res))
	return objArr
}

func setErrorResult(env *C.JNIEnv, errMsg string) C.jobjectArray {
	return setResult(env, "error", errMsg)
}

func setSuccessResult(env *C.JNIEnv, succMsg string) C.jobjectArray {
	return setResult(env, "success", succMsg)
}

func getAddressByPrivateKeyExecute(env *C.JNIEnv, args []string) C.jobjectArray {
	if len(args) != 1 {
		return setErrorResult(env, "error: "+GET_ADDRESS_BY_PRIVATE_KEY_CMD+" wrong argument count")
	}

	pub, addr, err := okwallet.GetAddressByPrivateKey(args[0])
	if err != nil {
		return setErrorResult(env, "error: "+err.Error())
	}

	var buf bytes.Buffer
	fmt.Fprintf(&buf, `{"pub":"%s",\n"address":"%s"}`, pub, addr)
	return setSuccessResult(env, buf.String())
}

func createRawTransactionExecute(env *C.JNIEnv, args []string) C.jobjectArray {
	if len(args) != 1 {
		return setErrorResult(env, "error: "+CREATE_RAW_TRANSACTION_CMD+" wrong argument count")
	}

	rawTx, err := okwallet.CreateRawTransaction(args[0])
	if err != nil {
		return setErrorResult(env, "error: "+err.Error())
	}

	return setSuccessResult(env, rawTx)
}

func SignRawTransactionExecute(env *C.JNIEnv, args []string) C.jobjectArray {
	if len(args) != 2 {
		return setErrorResult(env, "error: "+SIGN_RAW_TRANSACTION+" wrong argument count")
	}

	signedTx, err := okwallet.SignRawTransaction(args[0], args[1])
	if err != nil {
		return setErrorResult(env, "error: "+err.Error())
	}

	return setSuccessResult(env, signedTx)
}

//export Java_com_okcoin_vault_jni_starchain_Starchainj_execute
func Java_com_okcoin_vault_jni_starchain_Starchainj_execute(env *C.JNIEnv, _ C.jclass, _ C.jstring, jcommand C.jstring) C.jobjectArray {
	command, err := jstring2string(env, jcommand)
	if err != nil {
		return setErrorResult(env, "error: "+err.Error())
	}

	sepExp, err := regexp.Compile(`\s+`)
	if err != nil {
		return setErrorResult(env, "error: "+err.Error())
	}

	args := sepExp.Split(command, -1)
	if len(args) < 2 {
		return setErrorResult(env, "error: invalid command")
	}

	switch args[0] {
	case GET_ADDRESS_BY_PRIVATE_KEY_CMD:
		return getAddressByPrivateKeyExecute(env, args[1:])
	case CREATE_RAW_TRANSACTION_CMD:
		return createRawTransactionExecute(env, args[1:])
	case SIGN_RAW_TRANSACTION:
		return SignRawTransactionExecute(env, args[1:])
	default:
		return setErrorResult(env, "error: unknown command: "+args[0])
	}
	return setErrorResult(env, "error: unknown command: "+args[0])
}

func main() {}
