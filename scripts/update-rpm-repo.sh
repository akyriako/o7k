#!/usr/bin/env bash

set -euo pipefail

if [[ -z "${RPM_GPG_KEY_ID:-}" ]]; then
	echo "RPM_GPG_KEY_ID is required"
	exit 1
fi

if [[ $# -ne 2 ]]; then
	echo "usage: $0 <dist-directory> <rpm-repository-directory>"
	exit 1
fi

DIST_DIR="$(realpath "$1")"
RPM_REPO_DIR="$(realpath "$2")"

X86_64_DIR="${RPM_REPO_DIR}/packages/x86_64"
AARCH64_DIR="${RPM_REPO_DIR}/packages/aarch64"

echo "Updating o7k RPM repository..."

mkdir -p "${X86_64_DIR}"
mkdir -p "${AARCH64_DIR}"

echo "Copying RPM packages..."

cp "${DIST_DIR}"/o7k_*_linux_amd64.rpm "${X86_64_DIR}/"
cp "${DIST_DIR}"/o7k_*_linux_arm64.rpm "${AARCH64_DIR}/"

echo "Generating repository metadata..."

createrepo_c "${RPM_REPO_DIR}"

echo "Signing repository metadata..."

rm -f "${RPM_REPO_DIR}/repodata/repomd.xml.asc"

gpg \
	--batch \
	--yes \
	--local-user "${RPM_GPG_KEY_ID}" \
	--armor \
	--detach-sign \
	--output "${RPM_REPO_DIR}/repodata/repomd.xml.asc" \
	"${RPM_REPO_DIR}/repodata/repomd.xml"

echo "RPM repository metadata generated."