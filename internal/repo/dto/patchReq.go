package entities

type PatchReq struct {
	Interval *int  `json:"interval,omitempty"`
	Timeout  *int  `json:"timeout,omitempty"`
	Active   *bool `json:"active,omitempty"`
}
