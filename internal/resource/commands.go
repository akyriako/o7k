package resource

type NavigateMsg struct {
	Resource string
	ID       string
}

type NavigateFilteredMsg struct {
	Resource string
	Field    string
	Value    string
}

type NavigateFilteredMultiMsg struct {
	Resource string
	Field    string
	Values   []string
}

type DetailsMsg struct {
	ID      string
	Content any
	Err     error
}

type NavigateScopedMsg struct {
	Resource string
	Scope    map[string]string
}

type ErrorMsg struct {
	Err error
}

type CommandCompletedMsg struct{}
