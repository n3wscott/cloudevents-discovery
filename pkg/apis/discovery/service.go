package discovery

type Service struct {
	ID                 string            `json:"id"`                           // "id": "[a globally unique UUID]",
	URL                string            `json:"url"`                          // "url": "[unique URL to this service]",
	Name               string            `json:"name"`                         // "name": "[unique name for this services]",
	Description        string            `json:"description,omitempty"`        // "description": "[human string]", ?
	DocsURL            string            `json:"docsurl,omitempty"`            //"docsurl": "[URL reference for human documentation]", ?
	SpecVersions       []string          `json:"specversions"`                 // "specversions": [ "[ce-specversion value]" + ],
	SubscriptionURL    string            `json:"subscriptionurl"`              // "subscriptionurl": "[URL to which the Subscribe request will be sent]",
	SubscriptionConfig map[string]string `json:"subscriptionconfig,omitempty"` // "subscriptionconfig": { ?  "[key]": "[value]", * }
	AuthScope          string            `json:"authscope,omitempty"`          //	"authscope": "[string]", ?
	Protocols          []string          `json:"protocols"`                    // "protocols": [ "[string]" + ],
	Types              []ServiceType     `json:"types,omitempty"`              //"types": [ ?
}

type ServiceType struct {
	Type              string                 `json:"type"`                        // "type": "[ce-type value]",
	Description       string                 `json:"description,omitempty"`       // "description": "[human string]", ?
	DataContentType   string                 `json:"datacontenttype,omitempty"`   // "datacontenttype": "[ce-datacontenttype value]", ?
	DataSchema        string                 `json:"dataschema,omitempty"`        // "dataschema": "[ce-dataschema URI]", ?
	DataSchemaType    string                 `json:"dataschematype,omitempty"`    // "dataschematype": "[string per RFC 2046]", ?
	DataSchemaContent string                 `json:"dataschemacontent,omitempty"` // "dataschemacontent": "[schema]", ?
	Extensions        []ServiceTypeExtension `json:"extensions,omitempty"`        // "extensions": [ ?
}

type ServiceTypeExtension struct {
	Name    string `json:"name"`    // "name": "[CE context attribute name]",
	Type    string `json:"type"`    // "type": "[CE type string]",
	SpecURL string `json:"specurl"` // "specurl": "[URL to specification defining the extension]" ?
}
