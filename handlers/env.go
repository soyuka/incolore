package handlers

import (
	t "github.com/soyuka/incolore/transports"
	c "github.com/soyuka/incolore/config"
)


type Env struct {
	Transport t.Transport
	Config c.Config
}
