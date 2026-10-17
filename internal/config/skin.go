package config

type ColorPair struct {
	Foreground string `yaml:"foreground"`
	Background string `yaml:"background"`
}

type Skin struct {
	ResourceTag   ColorPair `yaml:"resource_tag"`
	NavigationTag ColorPair `yaml:"navigation_tag"`
	Status        ColorPair `yaml:"status"`
	Loading       ColorPair `yaml:"loading"`
	Info          ColorPair `yaml:"info"`
	StateOn       ColorPair `yaml:"state_on"`
	StateOff      ColorPair `yaml:"state_off"`
	ModalButton   ColorPair `yaml:"modal_button"`
	TableSelected ColorPair `yaml:"table_selected"`

	Command            string `yaml:"command"`
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

	TableTitle   string `yaml:"table_title"`
	TableBorder  string `yaml:"table_border"`
	DetailBorder string `yaml:"detail_border"`

	ErrorModal  ColorPair `yaml:"error_modal"`
	ErrorBorder string    `yaml:"error_border"`
	ErrorTitle  string    `yaml:"error_title"`
}
