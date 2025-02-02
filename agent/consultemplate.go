package agent

import (
	"log"

	"github.com/pkg/errors"

	consultemplateconfig "github.com/hashicorp/consul-template/config"
	"github.com/hashicorp/consul-template/manager"
)

func NewTemplateRunner(source string, destination string) (*manager.Runner, error) {
	config := consultemplateconfig.DefaultConfig()
	*config.Templates = append(
		*config.Templates,
		&consultemplateconfig.TemplateConfig{
			Source:      &source,
			Destination: &destination,
		},
	)

	runner, err := manager.NewRunner(config, false)
	if err != nil {
		return nil, errors.Wrap(err, "could not start consul-template")
	}

	return runner, nil
}

func (a *Agent) startConsulTemplate() error {
	go a.templateRunner.Start()
	defer a.stopConsulTemplate()

	for {
		select {
		case <-a.templateRunner.TemplateRenderedCh():
			log.Printf("Template rendered. Reloading agent configuration...")
			err := a.consul.Agent().Reload()
			if err != nil {
				log.Printf("Error reloading agent meta-data information: %s", err)
			} else {
				log.Print("Agent meta-data correctly reloaded")
			}
		case <-a.ctx.Done():
			return nil
		}
	}
}

func (a *Agent) stopConsulTemplate() {
	log.Println("Stopping consul-template")
	a.templateRunner.StopImmediately()
	log.Println("Stopped consul-template")
}
