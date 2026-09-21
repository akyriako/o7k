#!/usr/bin/env bash

set -euo pipefail

if [[ -z "${APT_GPG_KEY_ID:-}" ]]; then
	echo "APT_GPG_KEY_ID is required"
	exit 1
fi

if [[ $# -ne 2 ]]; then
	echo "usage: $0 <dist-directory> <apt-repository-directory>"
	exit 1
fi

DIST_DIR="$(realpath "$1")"
APT_REPO_DIR="$(realpath "$2")"

DISTRIBUTION="stable"
COMPONENT="main"

POOL_DIR="${APT_REPO_DIR}/pool/main/o/o7k"
DIST_ROOT="${APT_REPO_DIR}/dists/${DISTRIBUTION}"

echo "Updating o7k APT repository..."

mkdir -p "${POOL_DIR}"
mkdir -p "${DIST_ROOT}/${COMPONENT}/binary-amd64"
mkdir -p "${DIST_ROOT}/${COMPONENT}/binary-arm64"

echo "Copying Debian packages..."

cp "${DIST_DIR}"/o7k_*_linux_amd64.deb "${POOL_DIR}/"
cp "${DIST_DIR}"/o7k_*_linux_arm64.deb "${POOL_DIR}/"

cd "${APT_REPO_DIR}"

echo "Generating amd64 package index..."

dpkg-scanpackages \
	-a amd64 \
	pool/main/o/o7k \
	/dev/null \
	> "${DIST_ROOT}/${COMPONENT}/binary-amd64/Packages"

gzip -9 -c \
	"${DIST_ROOT}/${COMPONENT}/binary-amd64/Packages" \
	> "${DIST_ROOT}/${COMPONENT}/binary-amd64/Packages.gz"

echo "Generating arm64 package index..."

dpkg-scanpackages \
	-a arm64 \
	pool/main/o/o7k \
	/dev/null \
	> "${DIST_ROOT}/${COMPONENT}/binary-arm64/Packages"

gzip -9 -c \
	"${DIST_ROOT}/${COMPONENT}/binary-arm64/Packages" \
	> "${DIST_ROOT}/${COMPONENT}/binary-arm64/Packages.gz"

echo "Generating Release metadata..."

apt-ftparchive \
	-o APT::FTPArchive::Release::Origin="o7k" \
	-o APT::FTPArchive::Release::Label="o7k" \
	-o APT::FTPArchive::Release::Suite="${DISTRIBUTION}" \
	-o APT::FTPArchive::Release::Codename="${DISTRIBUTION}" \
	-o APT::FTPArchive::Release::Architectures="amd64 arm64" \
	-o APT::FTPArchive::Release::Components="${COMPONENT}" \
	release "dists/${DISTRIBUTION}" \
	> "${DIST_ROOT}/Release"

echo "Signing Release metadata..."

rm -f "${DIST_ROOT}/InRelease" "${DIST_ROOT}/Release.gpg"

gpg \
	--batch \
	--yes \
	--local-user "${APT_GPG_KEY_ID}" \
	--armor \
	--detach-sign \
	--output "${DIST_ROOT}/Release.gpg" \
	"${DIST_ROOT}/Release"

gpg \
	--batch \
	--yes \
	--local-user "${APT_GPG_KEY_ID}" \
	--clearsign \
	--output "${DIST_ROOT}/InRelease" \
	"${DIST_ROOT}/Release"

echo "APT repository metadata generated."