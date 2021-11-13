package constant

import (
	"time"
)

const (
	DefaultTimeout  = 60 * time.Second
	DefaultGasPrice = 0 // 1 gwai

	// For testing
	TestEndpoint        = "http://127.0.0.1:8545"
	TestMindenERC20Addr = "0xe868feADdAA8965b6e64BDD50a14cD41e3D5245D"
	TestPrivKey         = "d1c71e71b06e248c8dbe94d49ef6d6b0d64f5d71b1e33a0f39e14dadb070304a"
)
