## Summary

<!-- Briefly describe what this Pull Request adds, fixes, or improves in Syncalibur. -->

## Component Area

- [ ] **Go Core Engine** (`cmd/syncalibur-engine/main.go`)
- [ ] **Daemon & Service Supervisor** (`syncalibur.sh`)
- [ ] **Public Installer** (`install.sh`)
- [ ] **Documentation / Config Templates** (`README.md`, `syncalibur.env.example`, `repos.list.example`)

## Verification & Safety Checklist

- [ ] `bash -n syncalibur.sh && bash -n install.sh` passes without syntax errors
- [ ] `go build ./cmd/syncalibur-engine/main.go` compiles cleanly with zero external dependencies
- [ ] `./syncalibur.sh doctor` passes all checks
- [ ] Personal `syncalibur.env` and `repos.list` files are **not** included in this PR
