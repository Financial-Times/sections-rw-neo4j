package sections

type Section struct {
	UUID          string `json:"uuid"`
	CanonicalName string `json:"canonicalName"`
	TmeIdentifier string `json:"tmeIdentifier,omitempty"`
	Type          string `json:"type,omitempty"`
}

type SectionLink struct {
	ApiUrl string `json:"apiUrl"`
}
