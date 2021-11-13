package eth

import (
	"encoding/hex"
	"math"
	"math/big"
	"testing"

	// "github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/shopspring/decimal"
)

func TestGenerateAddrs(t *testing.T) {
	var (
		size = 10
	)
	addrs, err := GenerateAddrs(size)
	if err != nil {
		t.Fatal("err GenerateAddrs", err)
	}

	if g, w := len(addrs), size; g != w {
		t.Errorf("got: %d, want: %d", g, w)
	}

	for _, addr := range addrs {
		if IsZeroAddress(addr) {
			t.Error("address is zero")
		}
	}
}

func TestGenerateRandInts(t *testing.T) {
	var (
		size = 10
		max  = math.MaxUint32
	)

	results := make([]uint32, size)
	GenerateRandInts(size, max, func(num, index int) {
		results[index] = uint32(num)
	})

	for _, result := range results {
		if 0 > int(result) || int(result) > max {
			t.Error("invalid range")
		}
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// From https://raw.githubusercontent.com/miguelmota/ethereum-development-with-go-book/master/code/util/util_test.go
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func TestPublicKeyBytesToAddress(t *testing.T) {
	t.Parallel()
	{
		publicKeyBytes, _ := hex.DecodeString("049a7df67f79246283fdc93af76d4f8cdd62c4886e8cd870944e817dd0b97934fdd7719d0810951e03418205868a5c1b40b192451367f28e0088dd75e15de40c05")
		got := PublicKeyBytesToAddress(publicKeyBytes).Hex()
		expected := "0x96216849c49358B10257cb55b28eA603c874b05E"

		if got != expected {
			t.Errorf("Expected %v, got %v", expected, got)
		}
	}
}

func TestIsValidAddress(t *testing.T) {
	t.Parallel()
	validAddress := "0x323b5d4c32345ced77393b3530b1eed0f346429d"
	invalidAddress := "0xabc"
	invalidAddress2 := "323b5d4c32345ced77393b3530b1eed0f346429d"
	{
		got := IsValidAddress(validAddress)
		expected := true

		if got != expected {
			t.Errorf("Expected %v, got %v", expected, got)
		}
	}

	{
		got := IsValidAddress(invalidAddress)
		expected := false

		if got != expected {
			t.Errorf("Expected %v, got %v", expected, got)
		}
	}

	{
		got := IsValidAddress(invalidAddress2)
		expected := false

		if got != expected {
			t.Errorf("Expected %v, got %v", expected, got)
		}
	}
}

func TestIsZeroAddress(t *testing.T) {
	t.Parallel()
	validAddress := common.HexToAddress("0x323b5d4c32345ced77393b3530b1eed0f346429d")
	zeroAddress := common.HexToAddress("0x0000000000000000000000000000000000000000")

	{
		isZeroAddress := IsZeroAddress(validAddress)

		if isZeroAddress {
			t.Error("Expected to be false")
		}
	}

	{
		isZeroAddress := IsZeroAddress(zeroAddress)

		if !isZeroAddress {
			t.Error("Expected to be true")
		}
	}

	{
		isZeroAddress := IsZeroAddress(validAddress.Hex())

		if isZeroAddress {
			t.Error("Expected to be false")
		}
	}

	{
		isZeroAddress := IsZeroAddress(zeroAddress.Hex())

		if !isZeroAddress {
			t.Error("Expected to be true")
		}
	}
}

func TestToWei(t *testing.T) {
	t.Parallel()
	amount := decimal.NewFromFloat(0.02)
	got := ToWei(amount, 18)
	expected := new(big.Int)
	expected.SetString("20000000000000000", 10)
	if got.Cmp(expected) != 0 {
		t.Errorf("Expected %s, got %s", expected, got)
	}
}

func TestToDecimal(t *testing.T) {
	t.Parallel()
	weiAmount := big.NewInt(0)
	weiAmount.SetString("20000000000000000", 10)
	ethAmount := ToDecimal(weiAmount, 18)
	f64, _ := ethAmount.Float64()
	expected := 0.02
	if f64 != expected {
		t.Errorf("%v does not equal expected %v", ethAmount, expected)
	}
}

func TestCalcGasLimit(t *testing.T) {
	t.Parallel()
	gasPrice := big.NewInt(0)
	gasPrice.SetString("2000000000", 10)
	gasLimit := uint64(21000)
	expected := big.NewInt(0)
	expected.SetString("42000000000000", 10)
	gasCost := CalcGasCost(gasLimit, gasPrice)
	if gasCost.Cmp(expected) != 0 {
		t.Errorf("expected %s, got %s", gasCost, expected)
	}
}

func TestSigRSV(t *testing.T) {
	t.Parallel()

	sig := "0x789a80053e4927d0a898db8e065e948f5cf086e32f9ccaa54c1908e22ac430c62621578113ddbb62d509bf6049b8fb544ab06d36f916685a2eb8e57ffadde02301"
	r, s, v := SigRSV(sig)
	expectedR := "789a80053e4927d0a898db8e065e948f5cf086e32f9ccaa54c1908e22ac430c6"
	expectedS := "2621578113ddbb62d509bf6049b8fb544ab06d36f916685a2eb8e57ffadde023"
	expectedV := uint8(28)
	if hexutil.Encode(r[:])[2:] != expectedR {
		t.FailNow()
	}
	if hexutil.Encode(s[:])[2:] != expectedS {
		t.FailNow()
	}
	if v != expectedV {
		t.FailNow()
	}
}
