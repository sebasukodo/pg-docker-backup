# pg-docker-backup

A CLI tool for creating, encrypting, decrypting, and restoring PostgreSQL database backups running in Docker containers.

The tool uses `pg_dump` to create backups from a PostgreSQL database inside a Docker container. These backups are then encrypted using AES-256-GCM for secure storage. Encrypted backups can later be decrypted and restored using `pg_restore`.

---

## ⚙️ Features

- Backup PostgreSQL databases from Docker containers
- AES-256-GCM encryption for secure backup storage
- Decrypt previously created backups
- Restore backups back into a PostgreSQL database

---

## 📦 Environment Variables (.env)

Before using the application, you must generate and configure a secure encryption key.

### 1. Generate an Encryption Key

Run the following command in your terminal:

```bash
openssl rand -base64 32
```

This generates a **Base64-encoded string** (≈44 characters) representing a secure 32-byte key for AES-256 encryption.

---

### 2. Create your `.env` file

Copy the example file:

```bash
cp .env.example .env
```

---

### 3. Configure `.env`

Open the `.env` file and set your values:

```env
ENCRYPT_KEY=your_generated_key_here==
BACKUP_FOLDER_PATH=/mnt/pg-backup
CONTAINER_NAME=postgres-container
DB_NAME=mydb
DB_USER=postgres
DB_PASSWORD=secret
RSYNC_DEST_HOST=user@127.0.0.1
RSYNC_DEST_DIR=/home/user/backups/
```

---

### ⚠️ Important Notes

* The encryption key is **required to decrypt your backups**
* Store this key in a **secure location** (e.g. password manager or vault)
* If the key is lost, **your backups cannot be recovered**
* Never commit your `.env` file to version control

* `RSYNC_DEST_HOST` and `RSYNC_DEST_DIR` are optional (enable remote sync if both are set)
* `BACKUP_FOLDER_PATH` defaults to the working directory if not set

---

## 🚀 Available Commands

### `encrypt`

Creates a PostgreSQL backup using `pg_dump` and encrypts it.

### `decrypt`

Decrypts an encrypted backup file.

### `restore`

Restores a decrypted PostgreSQL backup file into a database.

### `help`

Displays help information for commands.

---

## 🔐 Encrypt Command

Creates and encrypts a PostgreSQL backup.

### Flags:

* `-c, --container` Docker container name
* `-n, --db-name` Database name
* `-p, --db-pw` Database password
* `-u, --db-user` Database username
* `-b, --backup-folder-path` Directory where backups will be stored (default: current working directory)
* `-h, --help` Help information

**Notes:**

* All flags are optional if you have configured the corresponding values in your `.env` file.
* You only need to set these flags manually if you are **not using `.env`** or want to override its values.

---

## 🔓 Decrypt Command

Decrypts an encrypted backup file.

### Flags:

* `-f, --file` Path to encrypted file (e.g. `./database-260312-1608.enc`)
* `-o, --output` Output file (default: `decrypted_backup.dump`)
* `-h, --help` Help information

**Notes:**

* The `-o, --output` flag is optional.

---

## ♻️ Restore Command

Restores a decrypted PostgreSQL dump into a running container.

### Flags:

* `-c, --container` Docker container name
* `-d, --db-name` Database name
* `-p, --db-pw` Database password
* `-u, --db-user` Database username
* `-f, --file` Path to decrypted backup file (e.g. `decrypted.dump`)
* `-h, --help` Help information

**Notes:**

* All flags except `-f, --file` are optional if you have configured the corresponding values in your `.env` file.
* You only need to set these flags manually if you are **not using `.env`** or want to override its values.

---

## 📌 Example Workflow

1. Create encrypted backup:

   ```bash
   pg-docker-backup encrypt -c my-postgres -n mydb -u user -p password
   ```

2. Decrypt backup:

   ```bash
   pg-docker-backup decrypt -f database-260312-1608.enc -o backup.dump
   ```

3. Restore backup:

   ```bash
   pg-docker-backup restore -c my-postgres -n mydb -u user -p password -f backup.dump
   ```

---

## ⏱️ Scheduling Backups

You can automate backups using either **Cron** or **systemd timers**.

### Using Cron

To run a backup every 6 hours, add the following line to your user crontab (`crontab -e`):

```bash
0 */6 * * * /path/to/script/run-backup.sh >> /path/to/script/backup.log 2>&1
```

**Notes:**

* Replace `/path/to/script/` with the path to where `run-backup.sh` and the `pg-docker-backup` Go binary are located.
* Output and errors are logged to `backup.log`.
* To remove the cron job, open `crontab -e` again and delete or comment out the line.

---

### Using systemd timers

1. Create a **service** file: `/etc/systemd/system/pg-docker-backup.service`

```ini
[Unit]
Description=Run pg-docker-backup

[Service]
Type=oneshot
WorkingDirectory=/path/to/script
ExecStart=/path/to/script/run-backup.sh
```

2. Create a **timer** file: `/etc/systemd/system/pg-docker-backup.timer`

```ini
[Unit]
Description=Run pg-docker-backup every 6 hours

[Timer]
OnBootSec=5min
OnUnitActiveSec=6h
Persistent=true

[Install]
WantedBy=timers.target
```

3. Reload systemd and start the timer:

```bash
sudo systemctl daemon-reload
sudo systemctl enable --now pg-docker-backup.timer
```

**Notes:**

* Replace `/path/to/script` with the actual path to where `run-backup.sh` and the `pg-docker-backup` go binary are located.
* To stop or remove the timer:

```bash
sudo systemctl stop pg-docker-backup.timer
sudo systemctl disable pg-docker-backup.timer
sudo rm /etc/systemd/system/pg-docker-backup.service
sudo rm /etc/systemd/system/pg-docker-backup.timer
sudo systemctl daemon-reload
```

This setup ensures automated, reliable backups every 6 hours.

---

## 🛡️ Security Notes

* Backups are encrypted using AES-256-GCM
* Store your `ENCRYPT_KEY` securely
* Do not commit `.env` files to version control

---

## 📄 License

This project is licensed under the MIT License.

---
