package cmd

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

var restoreFile string

// restoreCmd represents the restore command
var restoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Restore a decrypted PostgreSQL backup into a container",
	Long: `Restores a decrypted pg_dump file into a running PostgreSQL container.

	The command copies the backup file into the container and runs pg_restore
	against the specified database.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		if restoreFile == "" {
			return fmt.Errorf("--file flag is required, use pg-docker-backup restore --help for more information")
		}

		if err := checkValid(); err != nil {
			return err
		}

		backupFileInsideContainer := "/tmp/decrypted_backup.dump"

		cpBackupToContainer := exec.Command(
			"docker",
			"cp", restoreFile,
			fmt.Sprintf("%v:%v", containerName, backupFileInsideContainer),
		)

		fmt.Printf("Copying file %v inside container %v", restoreFile, containerName)

		if err := cpBackupToContainer.Run(); err != nil {
			return fmt.Errorf("could not copy backup to container: %v", err)
		}

		dockerCmd := exec.Command(
			"docker",
			"exec", "-e",
			"PGPASSWORD="+dbPW,
			containerName,
			"pg_restore",
			"-d", dbName,
			"-U", dbUser,
			"--clean",
			backupFileInsideContainer,
		)

		fmt.Println("Running command...")

		if err := dockerCmd.Run(); err != nil {
			return fmt.Errorf("Restore failed: %v", err)
		}

		fmt.Println("Restore completed successfully")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(restoreCmd)

	restoreCmd.Flags().StringVarP(&restoreFile, "file", "f", "", "Filepath to decrypted backup, e.g. decrypted.dump")

	restoreCmd.Flags().StringVarP(&containerName, "container", "c", containerName, "Docker Container Name")
	restoreCmd.Flags().StringVarP(&dbName, "db-name", "d", dbName, "Database Name")
	restoreCmd.Flags().StringVarP(&dbUser, "db-user", "u", dbUser, "Database Username")
	restoreCmd.Flags().StringVarP(&dbPW, "db-pw", "p", dbPW, "Database Password")

	restoreCmd.MarkFlagRequired("file")

	restoreCmd.MarkFlagFilename("file", ".dump")
}
