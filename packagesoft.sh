#!/bin/bash

# [ PACKAGESOFT - BASH BUILDER ]
# [ Author: Muhammad Quwais Saputra ]

RED='\033[0;31m'
GREEN='\033[0;32m'
CYAN='\033[0;36m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
NC='\033[0m'

show_banner() {
    echo -e "${PURPLE}"
    echo "          __            __"
    echo "         /  \\          /  \\"
    echo "        |    \\        /    |"
    echo "        |     \\      /     |"
    echo "         \\     \\    /     /"
    echo "          \\     \\  /     /"
    echo "           \\     \\/     /"
    echo "            \\    /\\    /"
    echo "             \\  /  \\  /"
    echo "              \\/    \\/"
    echo "      [ PACKAGESOFT - BASH BUILDER ]"
    echo "      [ Author: Muhammad Quwais Saputra ]"
    echo -e "${NC}"
}

if [ -z "$1" ]; then
    show_banner
    echo -e "${YELLOW}Usage: ./packagesoft.sh <filename.go>${NC}"
    exit 1
fi

SOURCE_FILE=$1
BIN_NAME=$(basename "$SOURCE_FILE" .go)
OUT_DIR="packagesoft_build"

mkdir -p "$OUT_DIR"

show_banner
echo -e "${BLUE}[*] Compiling Go Source: ${YELLOW}$SOURCE_FILE${NC}"
echo "------------------------------------------------"

# Daftar target
targets=(
    "linux amd64"
    "linux 386"
    "linux arm64"
    "linux arm"
    "windows amd64 .exe"
    "windows 386 .exe"
    "darwin amd64"
    "darwin arm64"
)

for t in "${targets[@]}"; do
    set -- $t
    os=$1
    arch=$2
    ext=$3
    
    suffix="${os}_${arch}"
    out_file="${OUT_DIR}/${BIN_NAME}_${suffix}${ext}"

    echo -ne "${CYAN}[>] Building for ${suffix}... ${NC}"

    # Build dengan CGO_ENABLED=0 untuk statis
    if GOOS=$os GOARCH=$arch CGO_ENABLED=0 go build -ldflags="-s -w" -o "$out_file" "$SOURCE_FILE" > /dev/null 2>&1; then
        echo -e "${GREEN}SUCCESS${NC}"
    else
        echo -e "${RED}FAILED${NC}"
    fi
done

echo "------------------------------------------------"
echo -e "${GREEN}[+] Build finished! Binaries are in '${OUT_DIR}/'${NC}"
