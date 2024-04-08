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
				From:  messages[i].From,
				Nonce: messages[i].Nonce,
				Order: messages[i].Order,
				Start: messages[i].Start,
				To:    messages[i].To,
				Gas:   messages[i].Gas,
				Data:  messages[i].Data,
			},
			OrigSign:  messages[i].OrigSign,
			ValidSign: messages[i].ValidSign,
		})
	}

	return responses
}
