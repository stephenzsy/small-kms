package agentconfigmanager

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	agentmodels "github.com/stephenzsy/small-kms/backend/models/agent"
)

type ConfigDir string

type RootConfigDir ConfigDir

type NamedConfigDir struct {
	RootConfigDir
	Version    string
	configName string
}

type WellKnownConfigFile string

const (
	configFileClientCert WellKnownConfigFile = "client-cert.pem"
)

func (dir RootConfigDir) Config(name agentmodels.AgentConfigName) NamedConfigDir {
	return NamedConfigDir{
		RootConfigDir: dir,
		configName:    string(name),
	}
}

func (dir NamedConfigDir) Versioned(version string) NamedConfigDir {
	dir.Version = version
	return dir
}

func (dir NamedConfigDir) ConfigFile(name WellKnownConfigFile, emptyIfNotExist bool) (fp string) {

	if dir.Version == "" {
		fp = filepath.Join(string(dir.RootConfigDir),
			"active",
			string(dir.configName),
			string(name))
	} else {
		fp = filepath.Join(string(dir.RootConfigDir),
			"versioned",
			fmt.Sprintf("%s.%s", dir.configName, dir.Version),
			string(name))
	}
	if emptyIfNotExist {
		if _, err := os.Stat(fp); err != nil && errors.Is(err, os.ErrNotExist) {
			return ""
		}
	}
	return fp
}

type CertsConfigDir struct {
	ConfigDir
}

func (dir RootConfigDir) Certs() CertsConfigDir {
	return CertsConfigDir{ConfigDir(filepath.Join(string(dir), "certs"))}
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

// type ConfigPath interface {
// 	Path() string
// 	EnsureDirExist() error
// 	WriteFile(data []byte) error
// 	ReadFile() ([]byte, error)
// }

// type ConfigDir interface {
// 	ConfigPath
// 	Active() ConfigDir
// 	Versioned(string) ConfigDir
// 	Dir(paths ...string) ConfigDir
// 	File(paths ...string) ConfigPath
// }

// type configPathImpl struct {
// 	configName  string
// 	path        string
// 	isVersioned bool
// 	isLeaf      bool
// }

// // Active implements ConfigDir.
// func (impl *configPathImpl) Active() ConfigDir {
// 	if impl.isVersioned {
// 		return impl
// 	}
// 	return &configPathImpl{
// 		configName:  impl.configName,
// 		path:        filepath.Join(impl.path, fmt.Sprint(impl.configName, ".active")),
// 		isVersioned: true,
// 	}
// }

// // WriteFile implements ConfigPath.
// func (impl *configPathImpl) WriteFile(data []byte) error {
// 	if !impl.isLeaf {
// 		return errors.New("not a leaf path")
// 	}
// 	return os.WriteFile(impl.path, data, 0640)
// }

// // WriteFile implements ConfigPath.
// func (impl *configPathImpl) ReadFile() ([]byte, error) {
// 	if !impl.isLeaf {
// 		return nil, errors.New("not a leaf path")
// 	}
// 	return os.ReadFile(impl.path)
// }

// // EnsureDirExist implements ConfigDir.
// func (impl *configPathImpl) EnsureDirExist() error {
// 	dirname := impl.path
// 	if impl.isLeaf {
// 		dirname = filepath.Dir(impl.path)
// 	}
// 	if s, err := os.Stat(dirname); err != nil {
// 		if !errors.Is(err, os.ErrNotExist) {
// 			return err
// 		}
// 		if err := os.MkdirAll(dirname, 0750); err != nil {
// 			return err
// 		}
// 	} else if !s.IsDir() {
// 		return errors.New("not a directory")
// 	}
// 	return nil
// }

// // File implements ConfigDir.
// func (impl *configPathImpl) File(paths ...string) ConfigPath {
// 	return &configPathImpl{
// 		configName:  impl.configName,
// 		path:        filepath.Join(impl.path, filepath.Join(paths...)),
// 		isVersioned: impl.isVersioned,
// 		isLeaf:      true,
// 	}
// }

// // Children implements ConfigDir.
// func (impl *configPathImpl) Dir(paths ...string) ConfigDir {
// 	return &configPathImpl{
// 		configName:  impl.configName,
// 		path:        filepath.Join(impl.path, filepath.Join(paths...)),
// 		isVersioned: impl.isVersioned,
// 	}
// }

// // Path implements ConfigDir.
// func (impl *configPathImpl) Path() string {
// 	return impl.path
// }

// func (impl *configPathImpl) Versioned(version string) ConfigDir {
// 	if impl.isVersioned {
// 		return impl
// 	}
// 	return &configPathImpl{
// 		configName:  impl.configName,
// 		path:        filepath.Join(impl.path, "versioned", fmt.Sprint(impl.configName, ".", version)),
// 		isVersioned: true,
// 	}
// }

// var _ ConfigDir = (*configPathImpl)(nil)

// func NewConfigDir(configName, basePath string) ConfigDir {
// 	if absPath, err := filepath.Abs(basePath); err == nil {
// 		return &configPathImpl{
// 			configName: configName,
// 			path:       absPath,
// 		}
// 	}
// 	return &configPathImpl{
// 		configName: configName,
// 		path:       basePath,
// 	}
// }
