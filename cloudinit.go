package cloud2butane

type CloudConfig struct {
	Users      []User   `yaml:"users"`
	WriteFiles []File   `yaml:"write_files"`
	Runcmd     []string `yaml:"runcmd"`
}

type User struct {
	Name   string   `yaml:"name"`
	Groups []string `yaml:"groups"`
	Shell  string   `yaml:"shell"`
}

type File struct {
	Path        string `yaml:"path"`
	Content     string `yaml:"content"`
	Append      bool   `yaml:"append"`
	Permissions string `yaml:"permissions"`
}
