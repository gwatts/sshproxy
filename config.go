// Copyright (c) 2016, Gareth Watts
// All rights reserved.

package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

type config struct {
	filename    string
	ProxyListen string            `yaml:"proxy-listen"`
	SSHHost     string            `yaml:"ssh-host"`
	SSHUser     string            `yaml:"ssh-user"`
	HostSig     string            `yaml:"host-sig,omitempty"`
	DNSMap      map[string]string `yaml:"dns-map"`
	SSHKey      string            `yaml:"ssh-key"`
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func loadConfig() (cfg config) {
	if *cfgFilename == "" {
		// try the current directory first
		dn := filepath.Dir(os.Args[0])
		for _, fn := range []string{defaultConfigFilename, filepath.Join(dn, defaultConfigFilename)} {
			if fileExists(fn) {
				return readConfig(fn)
			}
		}
		log.Fatal("No configuration file found")
	}
	return readConfig(*cfgFilename)
}

func readConfig(filename string) (cfg config) {
	cfg.filename = filename
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Failed to read %s: %v", filename, err)
	}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Fatalf("Failed to parse %s: %v", filename, err)
	}
	if cfg.ProxyListen == "" {
		cfg.ProxyListen = defaultProxyListen
	}
	log.Println("Loaded configuration from", filename)
	return cfg
}

func (cfg config) save() {
	data, err := yaml.Marshal(&cfg)
	if err != nil {
		log.Fatal("Failed to marshal yaml data", err)
	}
	if err := ioutil.WriteFile(cfg.filename, data, os.ModePerm); err != nil {
		log.Fatal("Failed to update configuration file:", err)
	}
}

func (cfg config) mapHost(name string) string {
	host := strings.ToLower(name)
	for k, v := range cfg.DNSMap {
		matched, err := filepath.Match(k, host)
		if err != nil {
			log.Fatalf("map host error for entry=%s host=%s", k, host)
		}
		if matched {
			return v
		}
	}
	return ""
}

func (cfg *config) hostSig() *hostSig {
	if cfg.HostSig == "" {
		return nil
	}
	s, err := hostSigFromString(cfg.HostSig)
	if err != nil {
		log.Println("Invalid host signature found in configuration file:", err)
	}
	return s
}

func (cfg *config) updateHostSig(sig *hostSig) {
	log.Println("Updated host-sig in configuration to", sig.String())
	cfg.HostSig = sig.String()
	cfg.save()
}
