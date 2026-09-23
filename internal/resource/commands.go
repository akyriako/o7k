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

type DetailsMsg struct {
	ID      string
	Content any
	Err     error
}

//type ResErrMsg struct {
//	Err error
//	Op  string
//}
