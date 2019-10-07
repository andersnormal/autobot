package message

const (
	// Src ...
	Src = "src"
	// Dest ...
	Dest = "dest"
)

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

// Src ...
func (m Metadata) Src(value string) {
	m[Src] = value
}

// Dest ...
func (m Metadata) Dest(value string) {
	m[Dest] = value
}
