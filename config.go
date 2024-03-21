package main

import (
	"encoding/json"
	"errors"
	"log"
	"os"
)

type config struct {
	cfgPath   string `json:"-"`
	Directory string `json:"directory"`
}

func InitialiseConfig(cfgPath string) config {
	cfg := config{
		cfgPath:   cfgPath,
		Directory: "",
	}
	err := cfg.ensureCFG()
	if err != nil {
		log.Panic(err)
	}
	err = cfg.loadCFG()
	if err != nil {
		log.Panic(err)
	}
	return cfg
}

func (cfg *config) ensureCFG() error {
	_, err := os.ReadFile(cfg.cfgPath)
	if errors.Is(err, os.ErrNotExist) {
		err = cfg.writeCFG()
	}
	return err
}

func (cfg *config) loadCFG() error {
	data, err := os.ReadFile(cfg.cfgPath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &cfg)
	return err
}

func (cfg *config) writeCFG() error {
	data, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	err = os.WriteFile(cfg.cfgPath, data, 0777)
	return err
}
