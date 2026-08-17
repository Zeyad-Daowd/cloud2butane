package cloud2butane

import "testing"

func TestIsValidSystemdUnit(t *testing.T) {
	tests := []struct {
		name string
		path string
		want bool
	}{
		{
			name: "service unit",
			path: "/etc/systemd/system/example.service",
			want: true,
		},
		{
			name: "timer unit",
			path: "/etc/systemd/system/example.timer",
			want: true,
		},
		{
			name: "socket unit",
			path: "/etc/systemd/system/example.socket",
			want: true,
		},
		{
			name: "mount unit",
			path: "/etc/systemd/system/example.mount",
			want: true,
		},
		{
			name: "target unit",
			path: "/etc/systemd/system/example.target",
			want: true,
		},
		{
			name: "swap unit",
			path: "/etc/systemd/system/example.swap",
			want: true,
		},
		{
			name: "not a systemd unit",
			path: "/etc/systemd/system/example.txt",
			want: false,
		},
		{
			name: "service outside systemd directory",
			path: "/etc/example.service",
			want: false,
		},
		{
			name: "wrong directory",
			path: "/etc/systemd/example.service",
			want: false,
		},
		{
			name: "no extension",
			path: "/etc/systemd/system/example",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isValidSystemdUnit(tt.path)

			if got != tt.want {
				t.Errorf("isValidSystemdUnit(%q) = %v, want %v",
					tt.path, got, tt.want)
			}
		})
	}
}

func TestAddButaneUnit(t *testing.T) {
	tests := []struct {
		name   string
		runcmd []string
		want   bool
	}{
		{
			name:   "enable unit",
			runcmd: []string{"systemctl enable example.service"},
			want:   true,
		},
		{
			name:   "enable unit with extra spaces",
			runcmd: []string{"systemctl  enable    example.service"},
			want:   true,
		},
		{
			name:   "start unit does not enable it",
			runcmd: []string{"systemctl start example.service"},
			want:   false,
		},
		{
			name:   "different unit does not enable it",
			runcmd: []string{"systemctl enable other.service"},
			want:   false,
		},
		{
			name:   "empty commands",
			runcmd: []string{},
			want:   false,
		},
		{
			name:   "disable unit",
			runcmd: []string{"systemctl disable example.service"},
			want:   false,
		},
		{
			name: "enable then disable",
			runcmd: []string{
				"systemctl enable example.service",
				"systemctl disable example.service",
			},
			want: false,
		},
		{
			name: "disable then enable",
			runcmd: []string{
				"systemctl disable example.service",
				"systemctl enable example.service",
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			systemd := ButaneSystemd{
				Units: []ButaneUnit{},
			}

			file := CloudFile{
				Path:    "/etc/systemd/system/example.service",
				Content: "test content",
			}

			err := addButaneUnit(&systemd, file, tt.runcmd)
			if err != nil {
				t.Fatalf("addButaneUnit() returned error: %v", err)
			}

			if len(systemd.Units) != 1 {
				t.Fatalf("expected 1 unit, got %d", len(systemd.Units))
			}

			if systemd.Units[0].Enabled != tt.want {
				t.Errorf(
					"Enabled = %v, want %v",
					systemd.Units[0].Enabled,
					tt.want,
				)
			}
		})
	}
}
