package utils

import "io"

type ChainedWriter struct {
	InnerWriter io.Writer
	count       int
	err         error
}

func (c *ChainedWriter) Write(p []byte) (int, error) {
	if c.err != nil {
		return 0, c.err
	}
	var count int
	count, c.err = c.InnerWriter.Write(p)
	c.count += count
	return count, c.err
}

func (c *ChainedWriter) WriteString(s string) (int, error) {
	if c.err != nil {
		return 0, c.err
	}
	var count int
	count, c.err = io.WriteString(c.InnerWriter, s)
	c.count += count
	return count, c.err
}

func (c *ChainedWriter) Return() (int, error) {
	return c.count, c.err
}
