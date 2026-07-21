package sdui

const SchemaVersion = 1

type Envelope struct {
	SchemaVersion int            `json:"schema_version"`
	Component     string         `json:"component"`
	Props         map[string]any `json:"props"`
}

func New(component string, props map[string]any) Envelope {
	return Envelope{SchemaVersion: SchemaVersion, Component: component, Props: props}
}
