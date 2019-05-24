APP_VERSION := "1.0.0"
COMMIT_ID := $(shell git log --format=%H -n 1)
PROJECT := $(shell pwd)
BUILD_DIR := $(PROJECT)/build
LIB_DIR := $(BUILD_DIR)/lib
SERVER_DIR := ./server

# packages to cover
PKGS := ./sources/reddit \
	./sources/stackexchange

LIBS := reddit \
	stackexchange

.PHONY: all
all: clean assets server libs run
	@echo "Build done"

.PHONY: run
run:
	$(BUILD_DIR)/server

.PHONY: server
server:
	@echo "Building server..."; \
	go build -o $(BUILD_DIR)/server -ldflags "-X main.version=$(APP_VERSION) -X main.commit=$(COMMIT_ID)" ./server

.PHONY: libs
libs:
	@$(foreach lib, $(LIBS), echo "Building '$(lib)' library"; \
		go build -buildmode=plugin -o $(LIB_DIR)/$(lib).so ./sources/$(lib);)

.PHONY: clean
clean:
	@if [ -d "$(BUILD_DIR)" ]; then \
		rm -rf "$(BUILD_DIR)"; \
		echo "Removed $(BUILD_DIR)"; \
	fi;

.PHONY: assets
assets:
	@echo "Copying assets..."; \
	if [ ! -d "$(BUILD_DIR)" ]; then \
		mkdir "$(BUILD_DIR)"; \
	fi; \
	cp $(SERVER_DIR)/config.json $(BUILD_DIR)/config.json; \
	cp -rf $(SERVER_DIR)/public $(BUILD_DIR)/public;

.PHONY: watch
watch:
	@if ! command -v "inotifywait" > /dev/null; then \
		echo "Please install 'inotifywait' tool:"; \
		echo "sudo apt install inotify-tools libnotify-bin"; \
		exit 1; \
	fi; \
	echo "Tracking changes in '$(PROJECT)/server'..."; \
	inotifywait -mrq -e create -e modify $(PROJECT)/server | \
	while read file; do (echo "Changed: $$file"; make all&) done

.PHONY: cover
cover:
	@$(foreach pkg,$(PKGS), go test -coverprofile=/tmp/cover.out $(pkg);\
	rm /tmp/cover.out;)
