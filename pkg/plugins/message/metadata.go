package message

// Metadata ...
type Metadata map[string]string

// Get ...
func (m Metadata) Get(key string) string {
	if v, ok := m[key]; ok {
		return v
	}

	return ""
}

// Set ...
func (m Metadata) Set(key, value string) {
	m[key] = value
}
