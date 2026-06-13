package devices

type Device struct {
	ID             string `json:"id"`
	Status         string `json:"status"`
	Manufacturer   string `json:"manufacturer"`
	Model          string `json:"model"`
	AndroidVersion string `json:"androidVersion"`
	SDK            string `json:"sdk"`
}
