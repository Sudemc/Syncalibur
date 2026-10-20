#!/usr/bin/env bash
# ==============================================================================
# Syncalibur One-Liner Public Installer (install.sh)
# ------------------------------------------------------------------------------
# Install or update Syncalibur globally in seconds:
#   curl -fsSL https://raw.githubusercontent.com/YusufKaya00/Syncalibur/main/install.sh | bash
#
# Or pass a target repository path directly:
#   curl -fsSL https://raw.githubusercontent.com/YusufKaya00/Syncalibur/main/install.sh | bash -s -- /path/to/repo
# ==============================================================================

set -euo pipefail

export PATH="${HOME}/.local/bin:/opt/homebrew/bin:/usr/local/bin:${PATH}"

INSTALL_DIR="${SYNCALIBUR_HOME:-${HOME}/Syncalibur}"
BIN_DIR="${HOME}/.local/bin"
REPO_URL="https://github.com/YusufKaya00/Syncalibur.git"
TARGET_WATCH="${1:-$(pwd)}"

C_RESET=$'\033[0m'
C_BOLD=$'\033[1m'
C_GOLD=$'\033[38;2;245;194;66m'
C_CYAN=$'\033[38;2;56;189;248m'
C_GREEN=$'\033[38;2;52;211;153m'
C_SLATE=$'\033[38;2;148;163;184m'

printf '\n%s%s◈ S Y N C A L I B U R   I N S T A L L E R%s\n' "${C_GOLD}" "${C_BOLD}" "${C_RESET}"
printf '%sAutonomous Git Auto-Sync Daemon — Zero-Friction Global Setup%s\n\n' "${C_SLATE}" "${C_RESET}"

command -v git >/dev/null 2>&1 || { echo "✕ Git is required but not installed." >&2; exit 1; }
command -v bash >/dev/null 2>&1 || { echo "✕ Bash is required." >&2; exit 1; }

if [[ -d "${INSTALL_DIR}/.git" ]]; then
  printf '  %s⟳ Updating existing Syncalibur installation in %s...%s\n' "${C_CYAN}" "${INSTALL_DIR}" "${C_RESET}"
  git -C "${INSTALL_DIR}" pull --rebase --autostash >/dev/null 2>&1 || true
else
  printf '  %s✦ Cloning Syncalibur into %s...%s\n' "${C_CYAN}" "${INSTALL_DIR}" "${C_RESET}"
  rm -rf "${INSTALL_DIR}"
  git clone --depth 1 "${REPO_URL}" "${INSTALL_DIR}" >/dev/null 2>&1
fi

chmod +x "${INSTALL_DIR}/syncalibur.sh"

# Compile Go engine if Go is installed on the machine
if command -v go >/dev/null 2>&1 && [[ -f "${INSTALL_DIR}/cmd/syncalibur-engine/main.go" ]]; then
  mkdir -p "${INSTALL_DIR}/bin"
  (cd "${INSTALL_DIR}" && go build -ldflags="-s -w" -o "${INSTALL_DIR}/bin/syncalibur-engine" ./cmd/syncalibur-engine 2>/dev/null) \
    && printf '  %s✓ Compiled native Go engine (bin/syncalibur-engine)%s\n' "${C_GREEN}" "${C_RESET}" || true
fi

# Create global CLI symlink (`syncalibur`)
mkdir -p "${BIN_DIR}"
ln -sf "${INSTALL_DIR}/syncalibur.sh" "${BIN_DIR}/syncalibur"
printf '  %s✓ Linked global CLI command: %s/syncalibur%s\n' "${C_GREEN}" "${BIN_DIR}" "${C_RESET}"

# Ensure ~/.local/bin is in user's shell RC
for rc in "${HOME}/.bashrc" "${HOME}/.zshrc" "${HOME}/.profile"; do
  if [[ -f "${rc}" ]] && ! grep -q '\.local/bin' "${rc}" 2>/dev/null; then
    printf '\n# Syncalibur CLI\nexport PATH="$HOME/.local/bin:$PATH"\n' >> "${rc}"
  fi
done

"${INSTALL_DIR}/syncalibur.sh" install "${TARGET_WATCH}"
