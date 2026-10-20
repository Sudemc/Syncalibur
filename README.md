<div align="center">

<img src="assets/syncalibur-logo.png" alt="Syncalibur Emblem" width="220" />

# Syncalibur

**Hybrid Go + Bash Autonomous Git Multi-Remote Sync Engine, Multi-Repo Fleet Supervisor & Live TUI Cockpit**

[![Platform](https://img.shields.io/badge/Platform-Linux%20%7C%20macOS%20%7C%20Windows-0ea5e9?style=for-the-badge)](#os-service-integration)
[![Engine](https://img.shields.io/badge/Core-Go%201.23%20%2B%20Bash-00ADD8?style=for-the-badge&logo=go&logoColor=white)](#architecture)
[![Git Remotes](https://img.shields.io/badge/Multi--Remote-GitHub%20%2B%20GitLab%20%2B%20Self--Hosted-10b981?style=for-the-badge)](#multi-remote--account-discovery)
[![Self-Updating](https://img.shields.io/badge/Self--Updating-Auto%20Git%20Pull%20%2B%20Go%20Rebuild-a855f7?style=for-the-badge)](#autonomous-self-update-watcher)

*Forge-grade background synchronization for your repositories — powered by a compiled Go core (`syncalibur-engine`) for sub-millisecond 2×2 Subpixel TrueColor emblem rendering, Semantic Conventional Commit generation, concurrent Multi-Repo Fleet orchestration, Multi-Channel Webhooks, and an Interactive Live TUI Cockpit.*

</div>

---

## One-Liner Quick Install

Anyone can install Syncalibur globally in seconds:

```bash
curl -fsSL https://raw.githubusercontent.com/YusufKaya00/Syncalibur/main/install.sh | bash
```

Or clone manually and watch any local repository:

```bash
git clone https://github.com/YusufKaya00/Syncalibur.git ~/Syncalibur
cd ~/Syncalibur
./syncalibur.sh install /path/to/your/repo
```

Once installed, the global `syncalibur` command is automatically linked into `~/.local/bin/syncalibur` so you can run `syncalibur ui`, `syncalibur status`, `syncalibur fleet`, or `syncalibur doctor` from any directory.

---

## What's New in Syncalibur v2 (Go Engine + Fleet + TUI Edition)

### Domain 1 — Git Intelligence, Bidirectional Sync, Fleet & Event-Driven Watching
1. **Compiled Go Core (`cmd/syncalibur-engine/main.go`):** Zero external Go dependencies (100% Go standard library). Automatically compiled into `bin/syncalibur-engine` on first run or after any self-update.
2. **Semantic Conventional Commit Engine:** Instead of generic timestamps, the Go AST/diff analyzer inspects staged changes, detects common directory scopes, classifies commit types (`feat`, `fix`, `docs`, `refactor`, `test`, `build`, `ci`, `chore`), and appends line addition/deletion metrics:
   - `feat(scripts): update syncalibur.sh, main.go (+142/-18)`
   - `docs: update README.md (+85/-12)`
   - Optional local **Ollama AI Commit Hook** (`SYNCALIBUR_AI_COMMIT="true"`) with a strict 2-second non-blocking timeout and automatic heuristic fallback.
3. **Safe Bidirectional Sync (`--rebase --autostash`):** Periodically fetches from your primary remote (`origin`), detects when remote commits are ahead of local `HEAD`, and safely runs `git pull --rebase --autostash` when the working tree is clean. If a merge conflict is detected, Syncalibur **immediately aborts the rebase** (`git rebase --abort`), leaves your working tree untouched, and pops open a diagnostic terminal window.
4. **Native Filesystem Event Watching (`inotifywait` / `fswatch`):** Zero-CPU event-driven watching when `inotifywait` (Linux) or `fswatch` (macOS) is installed, with automatic fallback to lightweight polling.
5. **Multi-Repository Fleet Supervisor (`repos.list`):** Watch 5, 10, or 50 repositories simultaneously! Run `syncalibur fleet scan ~/projects` to auto-discover all Git repositories using concurrent Go goroutines and supervise them in parallel.

### Domain 3 — Interactive Live TUI Cockpit, Multi-Channel Webhooks & Diagnostics
1. **Interactive Live TUI Cockpit (`syncalibur ui`):** Full-screen, flicker-free terminal control center built in Go using raw terminal mode and alternate screen buffering. Features the 2×2 subpixel TrueColor emblem, live daemon telemetry, 7-day commit velocity sparkline (`▂█▅··▅▁`), remote matrix, live log stream, and single-keystroke hotkeys:
   - `[p]` Pause / Resume Auto-Sync
   - `[s]` Force Immediate Commit & Push
   - `[b]` Force Bidirectional Pull (`--rebase`)
   - `[u]` Check & Apply Syncalibur Self-Update
   - `[d]` Run 10-Point System Doctor
   - `[q]` Quit TUI cleanly
2. **Multi-Channel Webhooks (`Discord`, `Slack`, `Telegram`, `Generic JSON`):** Configure `SYNCALIBUR_WEBHOOK_URL` (or Telegram Bot credentials) in `syncalibur.env` to receive instant rich notifications on push failures, rebase conflicts, self-updates, or every sync.
3. **10-Point `syncalibur doctor` Diagnostic Suite:** Verifies Git identity, Go engine binary health, filesystem watcher availability, SSH connectivity (`github.com` / `gitlab.com`), `gh` / `glab` CLI authentication, OS background service registration, and remote push permissions.

---

## Terminal Preview (`syncalibur status` & `syncalibur ui`)

When you run `syncalibur status`, `syncalibur logs`, or `syncalibur ui`, the Go engine renders [`assets/syncalibur-logo.png`](assets/syncalibur-logo.png) directly inside your terminal using **24-bit TrueColor ANSI escape sequences** and **2×2 Unicode Quadrant Subpixel Blocks** (`▘ ▝ ▀ ▖ ▌ ▞ ▛ ▗ ▚ ▐ ▜ ▄ ▙ ▟ █`) alongside live repository telemetry:

```text
          ▄▄▄▄▄▄▄▄▄▄                    Syncalibur • Autonomous Git Sync Engine
       ▄██████████████▄                 ────────────────────────────────────────────────────────
     ▄██████████████████▄               Engine      ● ACTIVE (PID 115794) • Go Engine Accelerated
    ██████████████████████              Service     systemd enabled (active)
   ████████████████████████             Repository  /home/user/my-project
  ██████████████████████████            Branch      main @ a1b2c3d
  ██████████████████████████            Strategy    all • debounce 30s • check 3s
  ██████████████████████████            Velocity    ▂█▅··▅▁ (7d: 220 commits • peak 85/day)
   ████████████████████████             Self-Update enabled • branch main (every 300s)
    ██████████████████████              Accounts    GitHub: @username • GitLab: @username
     ▀██████████████████▀               Remotes     2 configured (push target: all)
       ▀██████████████▀                             ├─ origin [GitHub] → username/my-project (acct: username)
          ▀▀▀▀▀▀▀▀▀▀                                └─ gitlab [GitLab] → group/my-project (acct: username)
```

---

## Architecture

```mermaid
flowchart TD
    A["Filesystem Change or Timer Tick"] --> B{"Event Watcher"}
    B -->|inotifywait / fswatch| C["Debounce Window (30s)"]
    B -->|Polling Fallback (3s)| C
    C --> D["Bidirectional Sync Check"]
    D -->|Remote Ahead| E["git pull --rebase --autostash"]
    E -->|Conflict Detected| F["git rebase --abort + Popup Alert Terminal + Webhook"]
    D -->|In Sync / Rebased| G["Safety Guard (<50 MiB & Secret Filter)"]
    G --> H["Go Engine: Semantic Commit Generator"]
    H -->|Heuristic or Ollama AI| I["git commit -m 'feat(scope): ... (+X/-Y)'"]
    I --> J{"Multi-Remote Push"}
    J -->|origin| K["GitHub"]
    J -->|gitlab| L["GitLab / Self-Hosted"]
    K & L -->|Success / Failure| M["Go Engine: Discord / Slack / Telegram Webhook"]
```

---

## CLI Command Reference

| Command | Description |
| :--- | :--- |
| `syncalibur ui` | Launch the **Interactive Full-Screen Live TUI Cockpit** (`[p]` Pause, `[s]` Sync, `[b]` Pull, `[u]` Update, `[d]` Doctor, `[q]` Quit) |
| `syncalibur install [PATH]` | Auto-detect repo, create missing remote repo if needed, link `~/.local/bin/syncalibur`, and start the OS service |
| `syncalibur status` | Render the 2×2 Subpixel TrueColor emblem, 7-day commit sparkline, remote matrix, and recent events |
| `syncalibur logs` | Render the emblem header and stream live color-coded log events (`tail -f`) |
| `syncalibur sync` | Trigger an immediate one-shot bidirectional pull, semantic commit, and multi-remote push |
| `syncalibur pull` | Trigger an immediate safe bidirectional pull (`git pull --rebase --autostash`) |
| `syncalibur pause` / `resume` | Temporarily pause or resume background auto-commits without uninstalling the service |
| `syncalibur doctor` | Run the 10-point diagnostic health check (Git, Go Engine, inotify, SSH, `gh`/`glab`, remotes) |
| `syncalibur fleet status` | Display a concurrent status table of all repositories in `repos.list` |
| `syncalibur fleet scan [DIR]` | Recursively scan a directory for Git repositories using Go goroutines and append them to `repos.list` |
| `syncalibur fleet start` | Spawn background Syncalibur daemons for every active repository in `repos.list` |
| `syncalibur fleet stop` | Stop all background fleet daemons |
| `syncalibur self-update` | Immediately check Syncalibur's own repository for updates, rebuild the Go engine, and reload |
| `syncalibur alert-test` | Test the graphical pop-up diagnostic terminal window on your desktop environment |
| `syncalibur uninstall` | Stop the background daemon and remove OS startup registration |

---

## Configuration (`syncalibur.env.example` → `syncalibur.env`)

Your personal `syncalibur.env` and `repos.list` files are **gitignored** so your local paths and webhook URLs are never pushed to GitHub. To customize settings, copy the template:

```bash
cp syncalibur.env.example syncalibur.env
```

```dotenv
# Target Git repository path to watch and auto-sync
SYNCALIBUR_REPO_DIR="/path/to/your/git/repository"

# Push target strategy:
#   "all"    -> Push to ALL configured git remotes (e.g., GitHub + GitLab simultaneously)
#   "origin" -> Push only to origin
SYNCALIBUR_REMOTE="all"

# Commit message style:
#   "semantic" -> Go Engine generates Conventional Commits: feat(scope): update X (+12/-3)
#   "simple"   -> Classic timestamped auto: file1, file2 (YYYY-MM-DD HH:MM:SS)
SYNCALIBUR_COMMIT_STYLE="semantic"

# Optional Local AI Commit Message Generation via Ollama (true / false)
SYNCALIBUR_AI_COMMIT="false"
SYNCALIBUR_OLLAMA_MODEL="qwen2.5-coder:1.5b"

# Safe Bidirectional Sync (auto-pull --rebase --autostash before pushing)
SYNCALIBUR_BIDIRECTIONAL="true"
SYNCALIBUR_PULL_INTERVAL="60"

# Native Filesystem Event Watching (inotifywait / fswatch with polling fallback)
SYNCALIBUR_USE_FS_EVENTS="true"

# Multi-Channel Webhooks (Discord, Slack, Telegram, or Generic JSON)
SYNCALIBUR_WEBHOOK_URL=""
SYNCALIBUR_WEBHOOK_ON_SUCCESS="false"
SYNCALIBUR_TELEGRAM_BOT_TOKEN=""
SYNCALIBUR_TELEGRAM_CHAT_ID=""
```

---

## Multi-Repository Fleet Mode (`repos.list`)

Want to watch multiple projects at once? Use **Fleet Mode**:

```bash
# Automatically scan a workspace directory and add all Git repositories to repos.list
syncalibur fleet scan ~/workspace

# View real-time parallel status across all repositories in your fleet
syncalibur fleet status

# Start background sync daemons for every repository in repos.list
syncalibur fleet start
```

---

## Safety & Self-Defense Mechanisms

1. **50 MiB File Size Guard:** Automatically skips files larger than `50 MiB` (`SYNCALIBUR_MAX_FILE_SIZE_BYTES`) so GitHub pushes never fail due to large binary blobs.
2. **Sensitive File Filter:** Automatically blocks `.env`, `.env.*`, `*.pem`, `*.key`, `id_rsa*`, `id_ed25519*`, `.terraform/`, `*.tfstate*`, `*.kubeconfig`, and `*.secret*` while allowing `.env.example` and `.env.template`.
3. **Conflict-Safe Rebase Abort:** If a remote change conflicts with local edits during bidirectional sync, Syncalibur immediately aborts the rebase, pauses auto-sync, sends a webhook alert, and opens a graphical diagnostic terminal window.
4. **Autonomous Self-Update Watcher:** Syncalibur monitors its own codebase (`~/Syncalibur`), pulls upstream improvements, recompiles `bin/syncalibur-engine` in Go, and hot-reloads the running daemon automatically.

---

## Repository Structure

```text
Syncalibur/
├── assets/
│   └── syncalibur-logo.png        # Official Syncalibur Emblem (rendered in 2x2 TrueColor ANSI)
├── cmd/
│   └── syncalibur-engine/
│       └── main.go                # Go Core Engine (PNG Renderer, Semantic Commits, Fleet, Webhooks, Live TUI)
├── go.mod                         # Go module definition (zero external dependencies)
├── install.sh                     # 1-Liner public curl|bash installer
├── syncalibur.sh                  # Cross-platform Service Manager & Daemon Supervisor
├── syncalibur.env.example         # Public configuration template (syncalibur.env is gitignored)
├── repos.list.example             # Multi-repository fleet configuration template
└── README.md                      # Documentation
```

---

## License

MIT License — Free for personal and commercial use.
