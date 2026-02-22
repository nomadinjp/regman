# Task Specification: Docker Private Registry Manager (Go CLI)

## 1. Project Overview
Create a CLI tool named `regman` (Registry Manager) in Go. The tool interacts with the Docker Registry V2 API to list repositories, list tags, and delete images.

## 2. Technical Stack
- **Language**: Go 1.21+
- **Key Library**: `github.com/google/go-containerregistry` (Industry standard for registry interaction)
- **CLI Framework**: `github.com/spf13/cobra`
- **Authentication**: Must support Basic Auth (Username/Password).

## 3. Core Functionalities & Commands

### Global Flags
- `--registry`: The URL of the private registry (e.g., `https://my-registry.com`).
- `--user`: Username for authentication.
- `--pass`: Password/Token for authentication.
- `--insecure`: Allow HTTP or skip TLS verification (Boolean).

### Command: `ls` (List Repositories)
- **Action**: Fetch and display all repository names from the `/_catalog` endpoint.
- **Output**: A simple list of image names.

### Command: `tags <image_name>`
- **Action**: Fetch all tags for a specific image using the `/<name>/tags/list` endpoint.
- **Output**: A list of tags associated with the image.

### Command: `rm <image_name>[:tag]`
- **Action**: Delete a specific tag or the entire image.
- **Implementation Logic (Critical)**:
    1.  The Registry API does **not** allow deletion by tag name directly.
    2.  First, fetch the `Manifest Digest` (Header: `Docker-Content-Digest`) via a `HEAD` request for the specific tag.
    3.  Execute a `DELETE` request using the retrieved `Digest`.
- **Note**: Ensure the tool handles the `Accept` header correctly to get the correct Digest.

## 4. Implementation Details (Guidance for AI)

### A. Authentication Wrapper
Use `remote.WithAuth` and `authn.Basic` from the `go-containerregistry` library. Provide a helper function to initialize `remote.Option`.

### B. Code Structure
```text
/regman
  /cmd
    root.go   (Global flags & help)
    ls.go     (List repos)
    tags.go   (List tags)
    rm.go     (Delete logic)
  main.go
  go.mod
```

### C. Example Code Snippet Logic (For AI Reference)
```go
// To get Digest for deletion
ref, _ := name.ParseReference("my-registry.com/my-image:tag")
img, _ := remote.Head(ref, remote.WithAuth(auth))
digest := img.Digest.String()

// To Delete
deleteRef, _ := name.NewDigest("my-registry.com/my-image@" + digest)
remote.Delete(deleteRef, remote.WithAuth(auth))
```

## 5. Deliverables Required
1.  A complete `go.mod` file.
2.  `main.go` entry point.
3.  Modular Cobra command files (`root.go`, `ls.go`, `tags.go`, `rm.go`).
4.  Robust error handling for network issues and 401 Unauthorized errors.
