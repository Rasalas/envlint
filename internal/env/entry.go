package env

// Entry represents a single key-value pair from an env file.
type Entry struct {
	Key      string
	Value    string
	Comment  string // inline comment after the value
	LineNum  int
	Required bool // determined by "# required" annotation or non-empty example value
	IsRef    bool // value contains variable reference ($VAR or ${VAR})
}
