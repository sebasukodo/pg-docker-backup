package config

import "os"

type Env struct {
	ContainerName string
	DBName        string
	DBUser        string
	DBPassword    string
	DockerMode    string
	EncryptKey    string
}

func Load() *Env {
	return &Env{
		ContainerName: os.Getenv("CONTAINER_NAME"),
		DBName:        os.Getenv("DB_NAME"),
		DBUser:        os.Getenv("DB_USER"),
		DBPassword:    os.Getenv("DB_PASSWORD"),
		DockerMode:    os.Getenv("DOCKER_MODE"),
		EncryptKey:    os.Getenv("ENCRYPT_KEY"),
	}
}
