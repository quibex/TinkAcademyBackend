package device

type Device struct {
	SerialNum string `json:"serial_num" validate:"required,numeric"`
	Model     string `json:"model" validate:"required"`
	IP        string `json:"ip" validate:"required,ip4_addr"`
}
