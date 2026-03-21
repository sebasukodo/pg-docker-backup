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
)

const stdOutputFile = "decrypted_backup.dump"

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

		if key == "" {
			return fmt.Errorf("Could not read 'ENCRYPT_KEY' variable.")
		}

		fmt.Println("Reading decryption key")

		encryptKey, err := base64.StdEncoding.DecodeString(key)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Reading encrypted data...")

		data, err := os.ReadFile(encryptFilename)
		if err != nil {
			panic(err)
		}

		block, err := aes.NewCipher(encryptKey)
		if err != nil {
			panic(err.Error())
		}

		gcm, err := cipher.NewGCM(block)
		if err != nil {
			panic(err.Error())
		}

		nonceSize := gcm.NonceSize()

		if len(data) < nonceSize {
			panic("filesize too short to decrypt")
		}

		nonce := data[:nonceSize]
		ciphertext := data[nonceSize:]

		fmt.Println("Decrypt data...")

		plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
		if err != nil {
			panic(err)
		}

		fmt.Println("Writing decrypted data into file...")

		err = os.WriteFile(toDecryptFilename, plaintext, 0644)
		if err != nil {
			os.Remove(toDecryptFilename)
			panic(err)
		}

		fmt.Printf("Decrypted backup saved in: %s\n", toDecryptFilename)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(decryptCmd)

	decryptCmd.Flags().StringVarP(&encryptFilename, "file", "f", "", "Path to encrypted file (e.g. ./database-260312-1608.enc)")
	decryptCmd.Flags().StringVarP(&toDecryptFilename, "output", "o", stdOutputFile, fmt.Sprintf("Output file for decrypted data (default: %v)", stdOutputFile))

	decryptCmd.MarkFlagRequired("file")

	decryptCmd.MarkFlagFilename("file", ".enc")
	decryptCmd.MarkFlagFilename("output", ".dump")
}
