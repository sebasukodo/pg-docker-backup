/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var (
	encryptFilename   string
	toDecryptFilename string
	encryptionKey     string
)

const stdOutputFile = "./decrypted_backup.dump"

var decryptCmd = &cobra.Command{
	Use:   "decrypt",
	Short: "Decrypt an encrypted backup file",
	Long: `Decrypts a previously encrypted backup file created by pg-docker-backup,
	restoring it to its original SQL dump format for recovery or inspection.`,

	RunE: func(cmd *cobra.Command, args []string) error {

		if encryptFilename == "" {
			return fmt.Errorf("--file flag is required, use pg-docker-backup decrypt --help for more information.")
		}

		if toDecryptFilename == "" {
			toDecryptFilename = stdOutputFile
		}

		if key == "" && encryptionKey == "" {
			return fmt.Errorf("Please set 'ENCRYPT_KEY' in .env or use --encryption-key flag")
		}

		if encryptionKey != "" {
			key = encryptionKey
		}

		fmt.Println("Reading decryption key")

		encryptKey, err := base64.StdEncoding.DecodeString(key)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Reading encrypted data...")

		data, err := os.ReadFile(encryptFilename)
		if err != nil {
			log.Fatal(err)
		}

		block, err := aes.NewCipher(encryptKey)
		if err != nil {
			log.Fatal(err.Error())
		}

		gcm, err := cipher.NewGCM(block)
		if err != nil {
			log.Fatal(err.Error())
		}

		nonceSize := gcm.NonceSize()

		if len(data) < nonceSize {
			log.Fatal("filesize too short to decrypt")
		}

		nonce := data[:nonceSize]
		ciphertext := data[nonceSize:]

		fmt.Println("Decrypt data...")

		plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Writing decrypted data into file...")

		err = os.WriteFile(toDecryptFilename, plaintext, 0644)
		if err != nil {
			os.Remove(toDecryptFilename)
			log.Fatal(err)
		}

		fmt.Printf("Decrypted backup saved in: %s\n", toDecryptFilename)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(decryptCmd)

	decryptCmd.Flags().StringVarP(&encryptionKey, "encryption-key", "e", "", "encryption key (use 'openssl rand -base64 32')")
	decryptCmd.Flags().StringVarP(&encryptFilename, "file", "f", "", "Path to encrypted file (e.g. ./database-260312-1608.enc)")
	decryptCmd.Flags().StringVarP(&toDecryptFilename, "output", "o", stdOutputFile, fmt.Sprintf("Output file for decrypted data (default: %v)", stdOutputFile))

	decryptCmd.MarkFlagRequired("file")

	decryptCmd.MarkFlagFilename("file", ".enc")
	decryptCmd.MarkFlagFilename("output", ".dump")
}
