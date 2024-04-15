package structs

type DBMessage struct {
	From      Address   `db:"from_"`
	Nonce     Uint256   `db:"nonce"`
	Order     Uint256   `db:"order_"`
	Start     Uint256   `db:"start"`
	To        Address   `db:"to_"`
	Gas       Uint256   `db:"gas"`
	Data      Bytes     `db:"data"`
	OrigSign  Signature `db:"orig_sign"`
	ValidSign Signature `db:"valid_sign"`
	ID        uint64    `db:"id"`
}

func WrapDBMessage(message Message, orig_sign Signature, valid_sign Signature) DBMessage {
	return DBMessage{
		From:      message.From,
		Nonce:     message.Nonce,
		Order:     message.Order,
		Start:     message.Start,
		To:        message.To,
		Gas:       message.Gas,
		Data:      message.Data,
		OrigSign:  orig_sign,
		ValidSign: valid_sign,
	}
}

func UnwrapDBMessage(message DBMessage) (Message, Signature, Signature) {
	return (Message{
		From:  message.From,
		Nonce: message.Nonce,
		Order: message.Order,
		Start: message.Start,
		To:    message.To,
		Gas:   message.Gas,
		Data:  message.Data,
	}), message.OrigSign, message.ValidSign
}
