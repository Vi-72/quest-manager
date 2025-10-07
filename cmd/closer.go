package cmd

import "log"

type Closer interface {
	Close() error
}

func (c *Container) RegisterCloser(closer Closer) {
	c.closers = append(c.closers, closer)
}

func (c *Container) CloseAll() {
	for _, closer := range c.closers {
		if err := closer.Close(); err != nil {
			log.Printf("error closing resource: %v", err)
		}
	}
}
