package structs

type HTTPValidateRequest struct {
	Message  Message   `json:"message"`
	OrigSign Signature `json:"origSign"`
}

type HTTPLoadResponse struct {
	ID        uint64    `json:"id"`
	Message   Message   `json:"message"`
	OrigSign  Signature `json:"origSign"`
	ValidSign Signature `json:"validSign"`
}

func WrapHTTPLoadResponses(messages []DBMessage) []HTTPLoadResponse {
	var responses []HTTPLoadResponse

	for i := 0; i < len(messages); i++ {
		responses = append(responses, HTTPLoadResponse{
			ID: messages[i].ID,
			Message: Message{
				Signer:   messages[i].Signer,
				Nonce:    messages[i].Nonce,
				GasOrder: messages[i].GasOrder,
				OnBehalf: messages[i].OnBehalf,
				Deadline: messages[i].Deadline,
				Endpoint: messages[i].Endpoint,
				Gas:      messages[i].Gas,
				Data:     messages[i].Data,
			},
			OrigSign:  messages[i].OrigSign,
			ValidSign: messages[i].ValidSign,
		})
	}

	return responses
}
