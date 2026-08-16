package cloud2butane

import (
	"fmt"
	"strconv"
)

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
	for _, file := range config.WriteFiles {
		mode, err := strconv.ParseInt(file.GetPermissions(), 8, 32)
		if err != nil {
			err = fmt.Errorf("invalid permissions %q for file %q: %w", file.Permissions, file.Path, err)
			return Butane{}, err
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
	}
	butane.Passwd = ButanePasswd
	butane.Storage = ButaneStorage
	return butane, nil
}
