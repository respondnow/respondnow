package version

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ghodss/yaml"
)

type Metadata struct {
	Name    string `json:"name" yaml:"name"`
	Version string `json:"version" yaml:"version"`
}

type VersionInfo struct {
	Version   string `json:"version" yaml:"version"`
	BuildNo   string `json:"buildNo" yaml:"buildNo"`
	GitCommit string `json:"gitCommit" yaml:"gitCommit"`
	GitBranch string `json:"gitBranch" yaml:"gitBranch"`
	Timestamp string `json:"timestamp" yaml:"timestamp"`
	Patch     string `json:"patch" yaml:"patch"`
}

type Resource struct {
	VersionInfo VersionInfo `json:"versionInfo" yaml:"versionInfo"`
}

type Version struct {
	MetaData         Metadata `json:"metaData" yaml:"metaData"`
	Resource         Resource `json:"resource" yaml:"resource"`
	ResponseMessages []string `json:"responseMessages" yaml:"responseMessages"`
}

func GetVersionInfo() (Version, error) {
	var defaultVersionPath string
	defaultVersionPath, err := filepath.Abs("version/versionInfo.yaml")
	if err != nil {
		return Version{}, fmt.Errorf("error while filepath.Abs %v", err)
	}

	var versionInfoPath string
	if respondNowServerVersionInfoPath, ok := os.LookupEnv("RESPOND_NOW_SERVER_VERSION_PATH"); ok {
		versionInfoPath = respondNowServerVersionInfoPath
	} else {
		versionInfoPath = defaultVersionPath
	}

	file, err := os.ReadFile(versionInfoPath)
	if err != nil {
		return Version{}, fmt.Errorf("error while reading file %v", err)
	}

	var versionYAML Version
	err = yaml.Unmarshal(file, &versionYAML)
	if err != nil {
		return Version{}, fmt.Errorf("error unmarshalling %v", err)
	}

	return versionYAML, nil
}
