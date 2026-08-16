package cloud2butane

type CloudConfig struct {
	Users      []CloudUser `yaml:"users"`
	WriteFiles []CloudFile `yaml:"write_files"`
	Runcmd     []string    `yaml:"runcmd"`
}

type CloudUser struct {
	Name   string   `yaml:"name"`
	Groups []string `yaml:"groups"`
	Shell  string   `yaml:"shell"`
}

type CloudFile struct {
	Path        string `yaml:"path"`
	Content     string `yaml:"content"`
	Append      bool   `yaml:"append"`
	Permissions string `yaml:"permissions"`
}
