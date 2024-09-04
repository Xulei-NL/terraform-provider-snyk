package cloudapi

type SastRequest struct {
	Data SastDataRequest `json:"data"`
}

type SastDataRequest struct {
	Attributes SastAttributesRequest `json:"attributes"`
	Type       string                `json:"type"`
	ID         string                `json:"id"`
}

type SastAttributesRequest struct {
	SastEnabled bool `json:"sast_enabled"`
}

type SastResponse struct {
	Data SastDataResponse `json:"data"`
}

type SastDataResponse struct {
	Attributes SastAttributesResponse `json:"attributes"`
	Type       string                 `json:"type"`
	ID         string                 `json:"id"`
}

type SastAttributesResponse struct {
	AutofixEnabled bool `json:"autofix_enabled"`
	SastEnabled    bool `json:"sast_enabled"`
}
