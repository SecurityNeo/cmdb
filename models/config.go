package models

type Config struct {
	DB DbConfig `yaml:"db"`
}

type DbConfig struct {
	Host      string `yaml:"host"`
	Port      string `yaml:"port"`
	DbName    string `yaml:"db-name"`
	UnameName string `yaml:"username"`
	Password  string `yaml:"password"`
}
