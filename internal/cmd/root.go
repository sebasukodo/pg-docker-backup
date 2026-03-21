package cmd

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/sebasukodo/pg-docker-backup/internal/config"
	"github.com/spf13/cobra"
)

var (
	containerName string
	dbName        string
	dbUser        string
	dbPW          string
	key           string
	dockerMode    string
)

var rootCmd = &cobra.Command{
	Use:   "pg-docker-backup",
	Short: "Backup and encrypt PostgresQL databases in Docker containers.",
	Long:  `pg-docker-backup is a simple CLI tool to create PostgresQL backups from Docker containers, with built-in encryption and decryption support.`,
}

func Execute() {

	godotenv.Load()
	env := config.Load()
	containerName = env.ContainerName
	dbName = env.DBName
	dbUser = env.DBUser
	dbPW = env.DBPassword
	key = env.EncryptKey
	dockerMode = env.DockerMode

	err := rootCmd.Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Oops. An error occured while executing pg-docker-backup:\n '%s'\n\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func checkValid() error {
	// Check if ENVIRONMENT Variables are set or not null
	if containerName == "" {
		return fmt.Errorf("--container flag is required if you haven't set Environment Variables")
	}
	if dbName == "" {
		return fmt.Errorf("--db-name flag is required if you haven't set Environment Variables")
	}
	if dbUser == "" {
		return fmt.Errorf("--db-user flag is required if you haven't set Environment Variables")
	}
	if dbPW == "" {
		return fmt.Errorf("--db-pw flag is required if you haven't set Environment Variables")
	}
	if dockerMode == "" {
		return fmt.Errorf("--docker-mode flag is required if you haven't set Environment Variables")
	}
	return nil
}
