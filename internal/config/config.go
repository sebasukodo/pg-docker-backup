package config

import "os"

type Env struct {
	ContainerName string
	DBName        string
	DBUser        string
	DBPassword    string
	BckFolderPath string
	EncryptKey    string
}

func Load() *Env {
	return &Env{
		ContainerName: os.Getenv("CONTAINER_NAME"),
		DBName:        os.Getenv("DB_NAME"),
		DBUser:        os.Getenv("DB_USER"),
		DBPassword:    os.Getenv("DB_PASSWORD"),
		BckFolderPath: os.Getenv("BACKUP_FOLDER_PATH"),
		EncryptKey:    os.Getenv("ENCRYPT_KEY"),
	}
}
