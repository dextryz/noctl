module github.com/ffiat/melange

go 1.19

replace github.com/ffiat/nostr => ../nostr

require (
	github.com/ffiat/nostr v0.0.0-00010101000000-000000000000
	github.com/gorilla/websocket v1.5.0
)

require (
	github.com/btcsuite/btcd/btcec/v2 v2.3.2 // indirect
	github.com/btcsuite/btcd/btcutil v1.1.3 // indirect
	github.com/btcsuite/btcd/chaincfg/chainhash v1.0.1 // indirect
	github.com/decred/dcrd/crypto/blake256 v1.0.0 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.0.1 // indirect
)
