package managedapp

import "io"

func (c *AgentContainerConfiguration) digest(writer io.Writer) {

	writer.Write([]byte(c.ContainerName))
	writer.Write([]byte(c.ImageRepo))
	writer.Write([]byte(c.ImageTag))
	for _, s := range c.ExposedPortSpecs {
		writer.Write([]byte(s))
	}
	for _, s := range c.HostBinds {
		writer.Write([]byte(s))
	}
	for _, s := range c.Secrets {
		writer.Write([]byte(s.Source))
		writer.Write([]byte(s.TargetName))
	}
	for _, s := range c.Env {
		writer.Write([]byte(s))
	}
	writer.Write([]byte(c.NetworkName))
}
