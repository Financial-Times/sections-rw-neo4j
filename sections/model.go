package sections

type Section struct {
	UUID                   string                 `json:"uuid"`
	PrefLabel              string                 `json:"prefLabel"`
	AlternativeIdentifiers alternativeIdentifiers `json:"alternativeIdentifiers"`
	Type                   string                 `json:"type,omitempty"`
}

type alternativeIdentifiers struct {
	TME   []string `json:"TME,omitempty"`
	UUIDS []string `json:"uuids"`
}

const (
	tmeIdentifierLabel = "TMEIdentifier"
	uppIdentifierLabel = "UPPIdentifier"
)

type SectionLink struct {
	ApiUrl string `json:"apiUrl"`
}
