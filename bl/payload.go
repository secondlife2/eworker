package bl

type PayloadCollection struct {
	WindowsVersion string    `json:"version"`
	Token          string    `json:"token"`
	Payloads       []Payload `json:"data"`
}

type Payload struct {
	UserName string `json:"username"`
}

func (p *Payload) HealthCheck() string {
	return "ok"
}
