package structs

type DBMessage struct {
	Signer    Address   `db:"signer"`
	Nonce     Uint256   `db:"nonce"`
	GasOrder  Uint256   `db:"gas_order"`
	OnBehalf  Address   `db:"on_behalf"`
	Deadline  Uint256   `db:"deadline"`
	Endpoint  Address   `db:"endpoint"`
	Gas       Uint256   `db:"gas"`
	Data      Bytes     `db:"data"`
	OrigSign  Signature `db:"orig_sign"`
	ValidSign Signature `db:"valid_sign"`
	ID        uint64    `db:"id"`
}

func WrapDBMessage(message Message, orig_sign Signature, valid_sign Signature) DBMessage {
	return DBMessage{
		Signer:    message.Signer,
		Nonce:     message.Nonce,
		GasOrder:  message.GasOrder,
		OnBehalf:  message.OnBehalf,
		Deadline:  message.Deadline,
		Endpoint:  message.Endpoint,
		Gas:       message.Gas,
		Data:      message.Data,
		OrigSign:  orig_sign,
		ValidSign: valid_sign,
	}
}
