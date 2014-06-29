package plugins

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/ehazlett/phoenix"
)

type (
	Plugin interface {
		Handle(*phoenix.Message) (string, error)
		Name() string
		Version() string
		Author() string
	}

	Manager struct {
		plugins        map[string]Plugin
		enabledPlugins []string
	}
)

func New(enabledPlugins []string) *Manager {
	// initialize plugins
	plugins := make(map[string]Plugin)
	plugins["example"] = Example()
	manager := &Manager{
		plugins:        plugins,
		enabledPlugins: enabledPlugins,
	}
	return manager
}

func (manager *Manager) Handle(msg *phoenix.Message) string {
	resp := ""
	// check for enabled plugin
	for _, p := range manager.enabledPlugins {
		// check if plugin matches trigger
		if msg.TriggerWord == p {
			for _, plugin := range manager.plugins {
				// if enabled plugin found, execute
				if plugin.Name() == p {
					logger.WithFields(logrus.Fields{
						"name":    plugin.Name(),
						"version": plugin.Version(),
						"author":  plugin.Author(),
					}).Infof("running plugin")
					r, err := plugin.Handle(msg)
					if err != nil {
						logMsg := fmt.Sprintf("error handling message: %s", err)
						logger.Errorf(logMsg)
						r = logMsg
					}
					resp = r
					break
				}
			}
		}
	}
	return resp
}