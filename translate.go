package cloud2butane

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
)

func isValidSystemdDirectory(path string) bool {
	const systemdDir = "/etc/systemd/system"
	// to handle cases like /etc/systemd/../systemd/system
	cleanPath := filepath.Clean(path)
	return filepath.Dir(cleanPath) == systemdDir
}

func isValidSystemdUnit(path string) bool {
	// check if the unit is a valid systemd unit
	// a valid systemd unit should have a valid extension
	if !isValidSystemdDirectory(path) {
		return false
	}
	//TODO: Extend this list
	validExtensions := []string{".service", ".timer", ".socket", ".mount", ".target", ".swap"}
	for _, extension := range validExtensions {
		if filepath.Ext(path) == extension {
			return true
		}
	}
	return false
}

func addButaneFile(ButaneStorage *ButaneStorage, file CloudFile) error {
	mode, err := strconv.ParseInt(file.GetPermissions(), 8, 32)
	if err != nil {
		err = fmt.Errorf("invalid permissions %q for file %q: %w", file.Permissions, file.Path, err)
		return err
	}

	butaneFile := ButaneFile{
		Path: file.Path,
		Contents: ButaneContent{
			Inline: file.Content,
		},
		//parse int from string to int
		Mode: Octal(mode),
	}

	ButaneStorage.Files = append(ButaneStorage.Files, butaneFile)
	return nil
}

func addButaneUnit(ButaneSystemd *ButaneSystemd, file CloudFile, runcmd []string) error {
	// extract file name from path by splitting / and getting last element
	fileName := filepath.Base(file.Path)
	// check if the file name is in the runcmd list
	enableUnit := false
	for _, cmd := range runcmd {
		// handle systemctl enable <unit> command
		parts := strings.Fields(cmd)
		// use parts to handle double spaces in command
		if len(parts) == 3 &&
			parts[0] == "systemctl" &&
			parts[1] == "enable" &&
			parts[2] == fileName {
			enableUnit = true
		}
	}
	butaneUnit := ButaneUnit{
		Name:     fileName,
		Enabled:  enableUnit,
		Contents: file.Content,
	}

	ButaneSystemd.Units = append(ButaneSystemd.Units, butaneUnit)
	return nil
}

func TranslateCloudConfig(config CloudConfig) (Butane, error) {
	var butane Butane
	butane.Variant = "fcos"
	butane.Version = "0.0.X"
	ButanePasswd := ButanePasswd{
		Users: []ButaneUser{},
	}
	for _, user := range config.Users {
		butaneUser := ButaneUser{
			Name:   user.Name,
			Groups: user.Groups,
			Shell:  user.Shell,
		}
		ButanePasswd.Users = append(ButanePasswd.Users, butaneUser)
	}
	ButaneStorage := ButaneStorage{
		Files: []ButaneFile{},
	}
	ButaneSystemd := ButaneSystemd{
		Units: []ButaneUnit{},
	}
	for _, file := range config.WriteFiles {
		if isValidSystemdUnit(file.Path) {
			err := addButaneUnit(&ButaneSystemd, file, config.Runcmd)
			if err != nil {
				return Butane{}, err
			}
		} else {
			err := addButaneFile(&ButaneStorage, file)
			if err != nil {
				return Butane{}, err
			}
		}

	}
	butane.Passwd = ButanePasswd
	butane.Storage = ButaneStorage
	butane.Systemd = ButaneSystemd
	return butane, nil
}
