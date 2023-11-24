package agentconfigmanager

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	agentmodels "github.com/stephenzsy/small-kms/backend/models/agent"
)

type ConfigDir string

type RootConfigDir struct {
	ConfigDir
}

type LeafConfigDir struct {
	ConfigDir
}

func (dir RootConfigDir) Active(name agentmodels.AgentConfigName) LeafConfigDir {
	return LeafConfigDir{
		ConfigDir: ConfigDir(filepath.Join(string(dir.ConfigDir), "active", string(name))),
	}
}

func (dir RootConfigDir) Versioned(name agentmodels.AgentConfigName, version string) LeafConfigDir {
	return LeafConfigDir{
		ConfigDir: ConfigDir(filepath.Join(string(dir.ConfigDir), "versioned", fmt.Sprintf("%s.%s", name, version))),
	}
}

type ConfigFile string

type WellKnownConfigFile string

const (
	configFileClientCert WellKnownConfigFile = "client-cert.pem"
	configFileServerCert WellKnownConfigFile = "server-cert.pem"
)

func (f ConfigFile) Exists() (bool, error) {
	if _, err := os.Stat(string(f)); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (f ConfigFile) LinkToAbsolutePath(targetAbsPath string) error {
	linkFileName := string(f)
	if _, err := os.Lstat(linkFileName); err == nil {
		// delete ink
		if err := os.Remove(linkFileName); err != nil {
			return err
		}
	}
	relpath, err := filepath.Rel(filepath.Dir(linkFileName), targetAbsPath)
	if err != nil {
		return err
	}
	return os.Symlink(relpath, linkFileName)
}

func (dir LeafConfigDir) ConfigFile(name WellKnownConfigFile) ConfigFile {
	return ConfigFile(filepath.Join(string(dir.ConfigDir), string(name)))
}

func (dir ConfigDir) EnsureExist() error {
	if _, err := os.Lstat(string(dir)); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if err := os.MkdirAll(string(dir), 0750); err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return nil
}

type CertsConfigDir struct {
	ConfigDir
}

func (dir RootConfigDir) Certs() CertsConfigDir {
	return CertsConfigDir{ConfigDir(filepath.Join(string(dir.ConfigDir), "certs"))}
}

func (dir ConfigDir) OpenFile(filename string, flag int, fileMode os.FileMode, ensureDirExist bool) (*os.File, error) {
	if _, err := os.Stat(string(dir)); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if ensureDirExist {
				if err := os.MkdirAll(string(dir), 0750); err != nil {
					return nil, err
				}
			}
		} else {
			return nil, err
		}
	}
	return os.OpenFile(filepath.Join(string(dir), filename), flag, fileMode)
}

func (dir ConfigDir) File(filename string) string {
	return filepath.Join(string(dir), filename)
}
