# pg-docker-backup

A CLI tool for creating, encrypting, decrypting, and restoring PostgreSQL database backups running in Docker containers.

The tool uses `pg_dump` to create backups from a PostgreSQL database inside a Docker container. These backups are then encrypted using AES-256-GCM for secure storage. Encrypted backups can later be decrypted and restored using `pg_restore`.

---

## ⚙️ Features

- Backup PostgreSQL databases from Docker containers
- AES-256-GCM encryption for secure backup storage
- Decrypt previously created backups
- Restore backups back into a PostgreSQL database
- Optional Docker mode for future container-in-container usage (experimental)

---

## 📦 Environment Variables (.env)

Before using the application, you must generate and configure a secure encryption key.

### 1. Generate an Encryption Key

Run the following command in your terminal:

```bash
openssl rand -base64 32
````

This generates a **Base64-encoded string** (usually ~44 characters) representing a secure 32-byte key for AES-256 encryption.

### 2. Configure `.env`

1. Copy the generated key
2. Open or create your `.env` file
3. Add the key as the value of `ENCRYPT_KEY`:

```env
ENCRYPT_KEY=your_generated_key_here
```

Example:

```env
CONTAINER_NAME=my-postgres
DB_NAME=mydb
DB_USER=postgres
DB_PASSWORD=secret
DOCKER_MODE=false
ENCRYPT_KEY=3K8kF9yourGeneratedKey
```

---

### ⚠️ Important Notes

* The encryption key is **required to decrypt your backups**
* Store this key in a **secure location** (e.g. password manager or vault)
* If the key is lost, **your backups cannot be recovered**
* Never commit your `.env` file to version control

### furthermore:

* `DOCKER_MODE=true` is intended for running the CLI inside a Docker container.
* This mode is **not fully implemented/tested yet**.

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
* `-m, --docker-mode` Run inside Docker container (`true` / `false`)
* `-h, --help` Help information

---

## 🔓 Decrypt Command

Decrypts an encrypted backup file.

### Flags:

* `-f, --file` Path to encrypted file (e.g. `./database-260312-1608.enc`)
* `-o, --output` Output file (default: `decrypted_backup.dump`)
* `-h, --help` Help information

---

## ♻️ Restore Command

Restores a decrypted PostgreSQL dump into a running container.

### Flags:

* `-c, --container` Docker container name
* `-d, --db-name` Database name
* `-p, --db-pw` Database password
* `-u, --db-user` Database username
* `-m, --docker-mode` Run inside Docker container (`true` / `false`)
* `-f, --file` Path to decrypted backup file (e.g. `decrypted.dump`)
* `-h, --help` Help information

---

## 🧪 Docker Mode (Experimental)

The `DOCKER_MODE` / `--docker-mode` flag is intended to allow running this CLI tool inside a Docker container instead of executing commands from the host system.

⚠️ This feature is currently **not fully implemented or tested**.

---

## 📌 Example Workflow

1. Create encrypted backup:

   ```bash
   pg-docker-backup encrypt -c my-postgres -n mydb -u user -p password
   ```

2. Decrypt backup:

   ```bash
   pg-docker-backup decrypt -f backup.enc -o backup.dump
   ```

3. Restore backup:

   ```bash
   pg-docker-backup restore -c my-postgres -n mydb -u user -p password -f backup.dump
   ```

---

## 🛡️ Security Notes

* Backups are encrypted using AES-256-GCM
* Store your `ENCRYPT_KEY` securely
* Do not commit `.env` files to version control

---

## 📄 License

This project is licensed under the MIT License.

---
