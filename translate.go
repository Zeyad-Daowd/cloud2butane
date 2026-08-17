package cloud2butane

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
)

func isValidSystemdDirectory(path string) bool {
	const systemdDir = "/etc/systemd/system"
	return filepath.Dir(path) == systemdDir
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

func checkUnitEnabled(fileName string, runcmd []string) bool {
	enableUnit := false
	for _, cmd := range runcmd {
		// handle systemctl enable <unit> command
		parts := strings.Fields(cmd)
		// use parts to handle double spaces in command
		if len(parts) >= 3 &&
			parts[0] == "systemctl" &&
			parts[len(parts)-1] == fileName {
			if parts[1] == "enable" {
				enableUnit = true
			} else if parts[1] == "disable" {
				enableUnit = false
			}
		}
	}
	return enableUnit
}

func addButaneUnit(ButaneSystemd *ButaneSystemd, file CloudFile, runcmd []string) error {
	// extract file name from path by splitting / and getting last element
	fileName := filepath.Base(file.Path)
	// check if the file name is in the runcmd list
	enableUnit := checkUnitEnabled(fileName, runcmd)
	butaneUnit := ButaneUnit{
		Name:     fileName,
		Enabled:  enableUnit,
		Contents: file.Content,
	}

	ButaneSystemd.Units = append(ButaneSystemd.Units, butaneUnit)
	return nil
}

func isValidDropin(file CloudFile) bool {
	if filepath.Ext(file.Path) != ".conf" {
		return false
	}
	directory := filepath.Dir(file.Path)
	if !strings.HasSuffix(directory, ".d") {
		return false
	}
	trimmed := strings.TrimSuffix(directory, ".d")
	return isValidSystemdUnit(trimmed)
}

func getDropinService(path string) string {
	directory := filepath.Dir(path)
	trimmed := strings.TrimSuffix(directory, ".d")
	return filepath.Base(trimmed)
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
	butaneStorage := ButaneStorage{
		Files: []ButaneFile{},
	}
	butaneSystemd := ButaneSystemd{
		Units: []ButaneUnit{},
	}
	dropinFiles := []CloudFile{}
	for _, file := range config.WriteFiles {
		file.Path = filepath.Clean(file.Path)
		if isValidSystemdUnit(file.Path) {
			err := addButaneUnit(&butaneSystemd, file, config.Runcmd)
			if err != nil {
				return Butane{}, err
			}
		} else if isValidDropin(file) {
			dropinFiles = append(dropinFiles, file)
		} else {
			err := addButaneFile(&butaneStorage, file)
			if err != nil {
				return Butane{}, err
			}
		}

	}

	for _, dropinFile := range dropinFiles {
		service := getDropinService(dropinFile.Path)
		foundService := false
		dropin := ButaneDropin{
			Name:     filepath.Base(dropinFile.Path),
			Contents: dropinFile.Content,
		}
		for i := range butaneSystemd.Units {
			butaneUnit := &butaneSystemd.Units[i]
			if service == butaneUnit.Name {
				foundService = true
				butaneUnit.Dropins = append(butaneUnit.Dropins, dropin)
				break
			}
		}
		if !foundService {
			// create stub unit
			enableUnit := checkUnitEnabled(service, config.Runcmd)
			unit := ButaneUnit{
				Name:    service,
				Dropins: []ButaneDropin{},
				Enabled: enableUnit,
			}
			unit.Dropins = append(unit.Dropins, dropin)
			butaneSystemd.Units = append(butaneSystemd.Units, unit)
		}
	}

	butane.Passwd = ButanePasswd
	butane.Storage = butaneStorage
	butane.Systemd = butaneSystemd
	return butane, nil
}
