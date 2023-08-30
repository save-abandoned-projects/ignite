package client

import (
	"github.com/save-abandoned-projects/ignite/pkg/client"
	"github.com/save-abandoned-projects/ignite/pkg/providers"
	log "github.com/sirupsen/logrus"
)

func SetClient() (err error) {
	log.Trace("Initializing the Client provider...")
	providers.Client = client.NewClient(providers.Storage)
	return
}
