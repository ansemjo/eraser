# Copyright (c) 2018 Anton Semjonov
# Licensed under the MIT License

# ---------- compile ----------

NAME := eraser

# build static binary with embedded version
VERSION     := $(shell sh version.sh describe)
BUILD_ENV   := CGO_ENABLED=0
BUILD_FLAGS := -ldflags='-s -w -X main.version=$(VERSION)'

# build binary for host system
.PHONY: build
build : $(NAME)
$(NAME) : $(shell find * -type f -name '*.go') go.mod
	env $(BUILD_ENV) go build $(BUILD_FLAGS) -o $@

# cross-compile binaries with gox
.PHONY: release
release :
	env $(BUILD_ENV) gox $(BUILD_FLAGS) -output='$@/$(NAME)-{{.OS}}-{{.Arch}}'

# run tests verbosely
.PHONY: test
test :
	go test ./... -timeout 30s -cover -v

# ---------- install ----------

DESTDIR :=
PREFIX 	:= /usr

# install binary
.PHONY: install
install : $(DESTDIR)$(PREFIX)/bin/$(NAME)

$(DESTDIR)$(PREFIX)/bin/$(NAME) : $(NAME)
	install -m 755 -D $< $@

# ---------- packaging ----------

# package metadata
PKGNAME     := $(NAME)
PKGVERSION  := $(shell echo $(VERSION) | sed -e 's/-\([0-9]\+\)/.r\1/' -e 's/-/./')
PKGAUTHOR   := 'ansemjo <anton@semjonov.de>'
PKGLICENSE  := MIT
PKGURL      := https://github.com/ansemjo/$(PKGNAME)
PKGFORMATS  := rpm deb apk
PKGARCH     := $(shell uname -m)

# how to execute fpm
DOCKER_BIN := podman
FPM := $(DOCKER_BIN) run --rm --net none -v $$PWD:/src -w /src ansemjo/fpm:alpine

# build a package
.PHONY: package-%
package-% :
	make --no-print-directory install DESTDIR=package
	mkdir -p release
	$(FPM) -s dir -t $* -f --chdir package \
		--name $(PKGNAME) \
		--version $(PKGVERSION) \
		--maintainer $(PKGAUTHOR) \
		--license $(PKGLICENSE) \
		--url $(PKGURL) \
		--architecture $(PKGARCH) \
		--package release/$(PKGNAME)-$(PKGVERSION)-$(PKGARCH).$*

# build all package formats with fpm
.PHONY: packages
packages : $(addprefix package-,$(PKGFORMATS))

# ---------- misc ----------

# clean untracked files and directories
.PHONY: clean
clean :
	git clean -fdx
