package stream

type InputEvent struct {
	Type string `json:"type"`

	X float64 `json:"x"`
	Y float64 `json:"y"`

	X2 float64 `json:"x2,omitempty"`
	Y2 float64 `json:"y2,omitempty"`

	Duration int `json:"duration,omitempty"`
}
