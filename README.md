# Regman (Registry Manager) 🐳

`regman` is a lightweight, Go-based CLI tool designed to manage Docker V2 private registries. It enables you to list repositories, query tags, and perform deep image deletions (including cleanup of backend storage indexes).

## ✨ Features

- **Multi-layered Configuration**: Supports command-line flags, environment variables (`REGMAN_*`), and a config file (`~/.regman.yaml`).
- **Smart Authentication**:
  - Supports Basic Auth (username/password).
  - **Automatic Docker Credentials Integration**: If no credentials are provided, it automatically reads login info from `~/.docker/config.json`.
- **Compliant Deletion**: Follows the Docker Registry V2 API strictly by retrieving the Manifest Digest before executing the DELETE request.
- **Automated Maintenance**: Includes a generalized Ansible Playbook for server-side Garbage Collection (GC) and empty index cleanup.

## 🚀 Installation

Ensure you have Go 1.21+ installed:

```bash
git clone <your-repo-url>
cd regman
go build -o regman
```

## ⚙️ Configuration

Configure `regman` using the following (ordered by priority):

1. **CLI Flags**: `--registry`, `--user`, `--pass`, `--insecure`
2. **Environment Variables**:
   ```bash
   export REGMAN_REGISTRY="my-registry.com"
   export REGMAN_USER="admin"
   export REGMAN_PASS="password"
   ```
3. **Config File** (`~/.regman.yaml`):
   ```yaml
   registry: "my-registry.com"
   user: "admin"
   pass: "password"
   insecure: true
   ```

## 📖 Usage

### List all repositories
```bash
./regman ls
```

### List all tags for an image
```bash
./regman tags my-app
```

### Delete a specific image/tag
*Note: Deletion must be enabled on the server side.*
```bash
./regman rm my-app:v1.0.1
```

---

## 🛠️ Server-side Maintenance (Critical)

### 1. Enable Deletion
By default, Docker Registry disables deletion. Update your server's `compose.yml`:

```yaml
services:
  registry:
    environment:
      REGISTRY_STORAGE_DELETE_ENABLED: "true"
```

### 2. GC & Index Cleanup
Executing `rm` only removes logical references. To free up disk space and remove empty repository names from the `ls` list, a maintenance script is required.

An automated solution is provided in the `ansible/` directory:

**Deploy the maintenance job:**
```bash
ansible-playbook ansible/setup-registry-gc.yml \
  -e "target_hosts=your-registry-server" \
  -e "registry_project_path=/path/to/registry/compose" \
  -e "registry_storage_path=/path/to/registry/data"
```

**Variables:**
- `target_hosts`: The host group in your inventory (default: `all`).
- `registry_project_path`: Directory containing your `compose.yml` (default: `/root/registry`).
- `registry_storage_path`: Root directory of registry storage on the host (default: `/data/registry`).

## 🏗️ Project Structure

```text
/regman
  /cmd          # Cobra command implementations (ls, tags, rm, root)
  /ansible      # Automation maintenance tools (Playbook)
  main.go       # Program entry point
  design.md     # Design specification
  go.mod        # Dependency management
```

## 📄 License

MIT License
