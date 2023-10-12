package main

type DockerComposeFile struct {
	Version string `yaml:"version"`
}

/*
func generateDockerCompose() {
	content, err := yaml.Marshal(DockerComposeFile{})
	if err != nil {
		panic(err)
	}

	os.WriteFile("compose.yaml", content, 0644)
}
*/
