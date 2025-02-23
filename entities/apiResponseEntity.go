package entities

type APIResponse struct {
	Status   bool   `json:"success"`
	Message  string `json:"message"`
	Response string `json:"response"`
}
