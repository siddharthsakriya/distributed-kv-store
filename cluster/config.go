package cluster

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Nodes []NodeConfig `yaml:"nodes"`
}

type NodeConfig struct {
	ID         string `yaml:"id"`
	Addr       string `yaml:"addr"`
	ClientAddr string `yaml:"client_addr"`
}

func Load(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("read config %s: %w", path, err)
	}
	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return Config{}, fmt.Errorf("parse config %s: %w", path, err)
	}
	if len(cfg.Nodes) == 0 {
		return Config{}, fmt.Errorf("config %s has no nodes", path)
	}
	return cfg, nil
}

func (c Config) View(myID string) (myAddr string, myClientAddr string, peers map[string]string, err error) {
	peers = make(map[string]string)
	found := false
	for _, n := range c.Nodes {
		if n.ID == myID {
			myAddr = n.Addr
			myClientAddr = n.ClientAddr
			found = true
		} else {
			peers[n.ID] = n.Addr
		}
	}
	if !found {
		return "", "", nil, fmt.Errorf("id %q not found in config", myID)
	}
	return myAddr, myClientAddr, peers, nil
}
