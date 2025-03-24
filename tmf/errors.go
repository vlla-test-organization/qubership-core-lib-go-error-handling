package tmf

var TypeV1_0 = "NC.TMFErrorResponse.v1.0"

type Error struct {
	Id             string                  `json:"id"`
	Code           string                  `json:"code"`
	Reason         string                  `json:"reason"`
	Message        *string                 `json:"message,omitempty"`
	ReferenceError *string                 `json:"referenceError,omitempty"`
	Status         *string                 `json:"status,omitempty"`
	Source         interface{}             `json:"source,omitempty"`
	Meta           *map[string]interface{} `json:"meta,omitempty"`
}

type Response struct {
	Id             string                  `json:"id"`
	Code           string                  `json:"code"`
	Reason         string                  `json:"reason"`
	Message        string                  `json:"message"`
	ReferenceError *string                 `json:"referenceError,omitempty"`
	Status         *string                 `json:"status,omitempty"`
	Source         interface{}             `json:"source,omitempty"`
	Meta           *map[string]interface{} `json:"meta,omitempty"`
	Errors         *[]Error                `json:"errors,omitempty"`
	Type           string                  `json:"@type"`
	SchemaLocation *string                 `json:"@schemaLocation,omitempty"`
}
