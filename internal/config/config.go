package config

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
