package config

type ColorPair struct {
	Foreground string `yaml:"foreground"`
	Background string `yaml:"background"`
}

type Skin struct {
	Main    ColorPair `yaml:"main"`
	Table   ColorPair `yaml:"table"`
	Context ColorPair `yaml:"context"`
	Command ColorPair `yaml:"command"`

	ResourceTag   ColorPair `yaml:"resource_tag"`
	NavigationTag ColorPair `yaml:"navigation_tag"`
	Status        ColorPair `yaml:"status"`
	Loading       ColorPair `yaml:"loading"`
	Info          ColorPair `yaml:"info"`
	StateOn       ColorPair `yaml:"state_on"`
	StateOff      ColorPair `yaml:"state_off"`
	ModalButton   ColorPair `yaml:"modal_button"`
	TableSelected ColorPair `yaml:"table_selected"`

	CommandGhost       string `yaml:"command_ghost"`
	CommandPlaceholder string `yaml:"command_placeholder"`

	HeaderLabel       string `yaml:"header_label"`
	HeaderValue       string `yaml:"header_value"`
	HeaderCommandKey  string `yaml:"header_command_key"`
	HeaderCommandText string `yaml:"header_command_text"`
	Logo              string `yaml:"logo"`

	JSONKey    string `yaml:"json_key"`
	JSONString string `yaml:"json_string"`
	JSONNumber string `yaml:"json_number"`
	JSONBool   string `yaml:"json_bool"`
	JSONNull   string `yaml:"json_null"`

	TableTitle  string `yaml:"table_title"`
	TableHeader string `yaml:"table_header"`
	ErrorTitle  string `yaml:"error_title"`
}
