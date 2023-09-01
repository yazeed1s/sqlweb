
BIN_DIR := bin
STATIC_DIR := static
LDFLAGS := -ldflags="-s -w"

clean_ui:
	@echo "Removing frontend artifacts...."
	@rm -rf ${STATIC_DIR}/build/*
	@sleep 1

clean_release:
	@echo "Removing backend artifacts...."
	@rm -rf ${BIN_DIR}/linux-amd-64-bit ${BIN_DIR}/linux-arm-64-bit ${BIN_DIR}/mac-amd-64-bit ${BIN_DIR}/mac-arm-64-bit
	@sleep 1

clean_bin: 
	@echo "Removing backend artifacts...."
	@rm -rf ${BIN_DIR}/sqlweb
	@sleep 1

build_ui:
	@echo "Buidling the frontend...."
	@mkdir -p ${STATIC_DIR}/build
	@cd ui && yarn build
	@mv ui/build/* ${STATIC_DIR}/build
	@rm -rf ui/build
	
build: 
	@echo "Buidling binaries..."
	$ go build ${LDFLAGS} -o $(BIN_DIR)/sqlweb ./cmd/main

build_linux_arm64:
	@echo "Buidling for linux-arm64...."
	$ GOOS=linux GOARCH=arm64 go build ${LDFLAGS} -o ${BIN_DIR}/linux-arm-64-bit/sqlweb ./cmd/main
	@sleep 2

build_linux_amd64:
	@echo "Buidling for linux-amd64...."
	$ GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o ${BIN_DIR}/linux-amd-64-bit/sqlweb ./cmd/main

build_mac_amd64:
	@echo "Buidling for mac-amd64...."
	$ GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o ${BIN_DIR}/mac-amd-64-bit/sqlweb ./cmd/main

build_mac_arm64:
	@echo "Buidling for mac-arm64...."
	$ GOOS=darwin GOARCH=arm64 go build ${LDFLAGS} -o ${BIN_DIR}/mac-arm-64-bit/sqlweb ./cmd/main

build_win_amd64:
	@echo "Buidling for windows-amd64...."
	$ GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o ${BIN_DIR}/windows-amd-64-bit/sqlweb ./cmd/main

build_win_arm64:
	@echo "Buidling for widnows-arm64...."
	$ GOOS=windows GOARCH=arm64 go build ${LDFLAGS} -o ${BIN_DIR}/windows-arm-64-bit/sqlweb ./cmd/main

install-bin:
	@echo "Moving binaries to /usr/local/bin...."
	@mv ${BIN_DIR}/sqlweb /usr/local/bin/

dev: clean build_ui build    
clean: clean_ui clean_bin
build: clean build_ui build
release: clean_release clean_ui build_ui build_linux_arm64 build_linux_amd64 build_mac_amd64 build_mac_arm64 build_win_amd64 build_win_arm64
install: install-bin