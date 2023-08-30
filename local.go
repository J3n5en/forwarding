package main

import (
	"os"
	"path"
	"strings"

	"github.com/mailway-app/config"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

func getDomainConfigFile(domain string) string {
	originConfigFilePath := path.Join(config.ROOT_LOCATION, "domain.d", domain+".yaml")
	_, err := os.Stat(originConfigFilePath)
	if err != nil {
		return originConfigFilePath
	}

	entries, err := os.ReadDir(path.Join(config.ROOT_LOCATION, "domain.d"))
	if err != nil {
		return ""
	}
	for _, e := range entries {
		if strings.HasSuffix(domain+".yaml", e.Name()) {
			return path.Join(config.ROOT_LOCATION, "domain.d", e.Name())
		}
	}
	return ""
}

func getLocalDomainConfig(instance *config.Config, domain string) (*Domain, error) {
	status := DOMAIN_UNCOMPLETE
	if getDomainConfigFile(domain) != "" {
		status = DOMAIN_ACTIVE
	} else {
		log.Warnf("No configuration for domain %s not found", domain)
	}
	return &Domain{
		Name:   domain,
		Status: status,
	}, nil
}

func getLocalDomainRules(instance *config.Config, domain string) (DomainRules, error) {
	config := DomainRules{}

	file := getDomainConfigFile(domain)
	content, err := os.ReadFile(file)
	if err != nil {
		return config, errors.Wrap(err, "could not read domain config")
	}

	err = yaml.Unmarshal(content, &config)
	if err != nil {
		return config, errors.Wrap(err, "failed to parse")
	}

	// legalize uuid
	for i := range config.Rules {
		config.Rules[i].Id = RuleId(uuid.Nil.String())
	}

	return config, nil
}
