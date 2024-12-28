package config

import "github.com/spf13/viper"

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Postgres PostgresConfig `yaml:"postgres"`
}

type ServerConfig struct {
	Port string `yaml:"port"`
}

type PostgresConfig struct {
	MasterDB Postgres `yaml:"master_db"`
	SlaveDB  Postgres `yaml:"slave_db"`
}

type Postgres struct {
	Host             string `yaml:"host"`
	Port             string `yaml:"port"`
	User             string `yaml:"user"`
	Password         string `yaml:"password"`
	DBName           string `yaml:"db_name"`
	AutoCreateTables bool   `yaml:"auto_create_tables"`
}

func LoadConfig(filename string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(filename)
	v.SetConfigType("yml")
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}
	cfg := &Config{}
	if err := v.Unmarshal(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
