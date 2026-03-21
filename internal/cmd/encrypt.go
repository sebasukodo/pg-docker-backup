package cmd

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/spf13/cobra"
)

var encryptCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "Backup and encrypt a PostgreSQL database from Docker",
	Long: `Creates a PostgreSQL backup from a Docker container using pg_dump,
encrypts it using AES-256-GCM, and removes the unencrypted dump file.`,

	RunE: func(cmd *cobra.Command, args []string) error {

		t := time.Now()
		timeText := t.Format("060102-1504")

		encFileName := fmt.Sprintf("%s-%s.enc", dbName, timeText)

		if key == "" {
			return fmt.Errorf("Could not read 'ENCRYPT_KEY' variable.")
		}

		fmt.Println("Starting encrypted backup process...")

		if err := checkValid(); err != nil {
			return err
		}

		// get pg_dump data
		fmt.Printf("Dumping database '%s' from container '%s'...\n", dbName, containerName)

		encryptKey, err := base64.StdEncoding.DecodeString(key)
		if err != nil {
			log.Fatal(err)
		}

		var command string
		var commandArgs []string

		if dockerMode == "true" {
			command = "pg_dump"
			commandArgs = append(commandArgs, "-h", containerName, "-U", dbUser, "-d", dbName, "-Fc")
		} else {
			command = "docker"
			commandArgs = append(commandArgs, "exec", "-e", "PGPASSWORD="+dbPW, containerName, "pg_dump", "-U", dbUser, "-d", dbName, "-Fc")
		}

		fmt.Println("Running command...")

		dockerCmd := exec.Command(command, commandArgs...)

		stdout, err := dockerCmd.Output()
		if err != nil {
			return fmt.Errorf("pg_dump failed: %w", err)
		}

		// encrypt pg_dump with .env key
		fmt.Println("Encrypting backup (AES-256-GCM)...")

		block, err := aes.NewCipher(encryptKey)
		if err != nil {
			panic(err.Error())
		}

		aesgcm, err := cipher.NewGCM(block)
		if err != nil {
			panic(err.Error())
		}

		nonce := make([]byte, aesgcm.NonceSize())
		if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
			panic(err.Error())
		}

		ciphertext := aesgcm.Seal(nonce, nonce, stdout, nil)

		// write encrypted data to file
		fmt.Println("Writing encrypted data into file...")

		err = os.WriteFile(encFileName, ciphertext, 0644)
		if err != nil {
			os.Remove(encFileName)
			panic(err)
		}

		fmt.Printf("Encrypted backup created: %s\n", encFileName)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(encryptCmd)

	encryptCmd.Flags().StringVarP(&containerName, "container", "c", containerName, "Docker Container Name")
	encryptCmd.Flags().StringVarP(&dbName, "db-name", "n", dbName, "Database Name")
	encryptCmd.Flags().StringVarP(&dbUser, "db-user", "u", dbUser, "Database Username")
	encryptCmd.Flags().StringVarP(&dbPW, "db-pw", "p", dbPW, "Database Password")
	encryptCmd.Flags().StringVarP(&dockerMode, "docker-mode", "m", dockerMode, "Are you running this application inside of a Docker container? - true or false")
}
