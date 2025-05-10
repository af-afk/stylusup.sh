#!/usr/bin/env -S bash -e

# This script will do the following:

# 1. Ask the user how much reporting they're comfortable with, and prompt them to simply
# press enter if they're okay with this script reporting the language spoken on their
# computer, their architecture, and their operating system. At the end of this script, it
# will report whether the installation was a success. This will form the basis of a
# Popularity Contest, which is accessible here: https://popcon.stylusup.sh. The rest of the
# repository hosted at https://github.com/af-afk/stylusup.sh contains the webapp which
# simply renders a graph using Stellate served from Postgres on a micro ec2 instance,
# intermediated by Cloudflare. IP addresses collected are the HMAC'd form of the
# architecture, language, and operating system details.

# 2. Check if Rust is installed. If it isn't, then it will use Rustup to make an
# installation. It'll check the success of that installation.

# 3. Install the missing dependencies involving Cargo.

# 4. Suggest to the user where to go from here.

os="$(uname -s)"
arch="$(uname -m)"
lang="${LANG%%_*}"

NC=$'\033[0m'; BLUE=$'\033[1;34m'; RED=$'\033[1;31m'; GREEN=$'\033[1;32m'

log() { >&2 printf "${BLUE}ℹ️  %s${NC}\n"  "$*"; }
die() { >&2 printf "${RED}❌ %s${NC}\n"  "$*"; exit 1; }

if ! which curl >/dev/null; then
	die "curl is needed for installation. You can install it with your package manager."
fi

if ! which cc >/dev/null; then
	die "cc is needed for installation. You can install it with your package manager."
fi

if [ "$os" = "Linux" ]; then
	. /etc/os-release
	distro="-$NAME"
fi

report_popcon() {
	curl \
		-sd "{\"query\":\"mutation {\n  register(arch: \\\"$arch\\\", lang: \\\"$lang\\\", os:\\\"$os$distro\\\")\n}\"}" \
		-H 'Content-Type: application/json' \
		https://markov-geist-research.stellate.sh/ 2>&1 >/dev/null
}

check_has_rust() {
	which cargo 2>&1 >/dev/null
}

if [ -z "$STYLUS_POPCON_OFF" ]; then
	log "This installer will record the language, the os, and the architecture to
https://popcon.stylusup.sh. To disable this functionality, control-c now, and set
STYLUS_POPCON_OFF to anything."
	log "Press enter to continue…"
	read -r < /dev/tty > /dev/tty
	report_popcon &
fi

if ! check_has_rust; then
	log "Installing Rust with Rustup..."
	# Use the Rust installer with the stable toolchain and no interactivity.
	curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- --default-toolchain stable -y
	# We need to get this into our PATH.
	. ~/.cargo/env
fi

if ! rustup target list | grep wasm32-unknown-unknown >/dev/null; then
	log "Installing wasm32-unknown-unknown..."
	rustup target add wasm32-unknown-unknown
fi

if ! cargo stylus 2>&1 >/dev/null; then
	log "Installing cargo-stylus..."
	cargo install cargo-stylus
fi

>&2 cat <<EOF
${GREEN}🎉  Congratulations!!! You're ready to develop with Stylus!${NC}

💡  Use "cargo stylus new" to get started with your first project!

🔧  . \$HOME/.cargo/env

🚀  cargo stylus new hello-world${NC}
EOF
