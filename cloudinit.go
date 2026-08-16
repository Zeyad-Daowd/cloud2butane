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

// handle default value for permissions if not provided in the cloud-config file
func (f CloudFile) GetPermissions() string {
	if f.Permissions == "" {
		return "0644"
	}
	return f.Permissions
}
