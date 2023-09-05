package structs

type DBMessage struct {
	Signer    Address
	Nonce     Uint256
	GasOrder  Uint256 `db:"gas_order"`
	OnBehalf  Address `db:"on_behalf"`
	Deadline  Uint256
	Endpoint  Address
	Gas       Uint256
	Data      Bytes
	OrigSign  Signature `db:"orig_sign"`
	ValidSign Signature `db:"valid_sign"`
	ID        uint64
}
