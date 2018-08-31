package okwallet

import (
	"testing"
)

func TestGetAddressByPrivateKey(t *testing.T) {
	const prvkey = `de47f09131a701b1012ab205df484ec43dd73daa6d42a0fead004a16307d7fdd`
	const expect = `ShoPrXXG2nMATE3WDA7B7sgJscCoNXauaY`

	addr, err := GetAddressByPrivateKey(prvkey)
	if err != nil {
		t.Fatal(err)
	}
	if expect != addr {
		t.Fatal("GetAddressByPrivateKey failed. expect " + expect + " but return " + addr)
	}
}
