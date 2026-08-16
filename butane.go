package cloud2butane

type Butane struct {
	Variant string        `yaml:"variant"`
	Version string        `yaml:"version"`
	Passwd  ButanePasswd  `yaml:"passwd"`
	Storage ButaneStorage `yaml:"storage"`
}

type ButanePasswd struct {
	Users []ButaneUser `yaml:"users"`
}

type ButaneStorage struct {
	Files []ButaneFile `yaml:"files"`
}
type ButaneUser struct {
	Name   string   `yaml:"name"`
	Groups []string `yaml:"groups"`
	Shell  string   `yaml:"shell"`
}

type ButaneFile struct {
	Path     string        `yaml:"path"`
	Contents ButaneContent `yaml:"contents"`
	// mode should be displayed as an octal number in the butane file, but it is stored as an int in the struct
	Mode int `yaml:"mode"`
}
type ButaneContent struct {
	Inline string `yaml:"inline"`
}
