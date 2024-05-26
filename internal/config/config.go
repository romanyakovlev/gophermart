package config

import (
	"flag"

	"github.com/caarlos0/env/v6"

	"github.com/romanyakovlev/gophermart/internal/logger"
)

type argConfig struct {
	flagAAddr string
	flagDAddr string
	flagRAddr string
}

type envConfig struct {
	RunAddress           string `env:"RUN_ADDRESS"`
	DatabaseURI          string `env:"DATABASE_URI"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
}

type Config struct {
	RunAddress           string
	DatabaseURI          string
	AccrualSystemAddress string
}

func parseEnvs(s *logger.Logger) envConfig {
	var cfg envConfig
	err := env.Parse(&cfg)
	if err != nil {
		s.Fatal(err)
	}
	return cfg
}

func parseFlags() argConfig {
	var cfg argConfig
	// указываем имя флага, значение по умолчанию и описание
	flag.StringVar(&cfg.flagAAddr, "a", "localhost:8080", "Адрес запуска HTTP-сервера")
	flag.StringVar(&cfg.flagDAddr, "d", "", "Строка с адресом подключения к БД")
	flag.StringVar(&cfg.flagRAddr, "r", "", "Адрес системы расчёта начислений")
	// делаем разбор командной строки
	flag.Parse()
	return cfg
}

func GetConfig(s *logger.Logger) Config {
	argCfg := parseFlags()
	envCfg := parseEnvs(s)

	var RunAddress string
	var DatabaseURI string
	var AccrualSystemAddress string

	if envCfg.RunAddress != "" {
		RunAddress = envCfg.RunAddress
	} else {
		RunAddress = argCfg.flagAAddr
	}
	if envCfg.DatabaseURI != "" {
		DatabaseURI = envCfg.DatabaseURI
	} else {
		DatabaseURI = argCfg.flagDAddr
	}
	if envCfg.AccrualSystemAddress != "" {
		AccrualSystemAddress = envCfg.AccrualSystemAddress
	} else {
		AccrualSystemAddress = argCfg.flagAAddr
	}
	return Config{
		RunAddress:           RunAddress,
		DatabaseURI:          DatabaseURI,
		AccrualSystemAddress: AccrualSystemAddress,
	}
}
