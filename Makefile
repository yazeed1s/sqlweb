BIN_DIR := bin
STATIC_DIR := static
LDFLAGS := -ldflags="-s -w"

clean_ui:
	@echo "Removing frontend artifacts...."
ifeq ($(OS),Windows_NT) # Windows
	rm -rf ${STATIC_DIR}\\build\\
else
	rm -rf ${STATIC_DIR}/build/
endif
	@sleep 1

clean_release:
	@echo "Removing backend artifacts...."
	@rm -rf ${BIN_DIR}/linux-amd-64-bit ${BIN_DIR}/linux-arm-64-bit ${BIN_DIR}/mac-amd-64-bit ${BIN_DIR}/mac-arm-64-bit
	@sleep 1

clean_bin: 
	@echo "Removing backend artifacts...."
ifeq ($(OS),Windows_NT) # Windows
	@rm -rf ${BIN_DIR}\\sqlweb
else
	@rm -rf ${BIN_DIR}/sqlweb
endif
	@sleep 1

build_ui:
	@echo "Buidling the frontend...."
ifeq ($(OS),Windows_NT) # Windows
	@mkdir ${STATIC_DIR}\\build
	@cd ui && yarn build
	@mv ui$\\build$\\* ${STATIC_DIR}$\\build
	@rm -rf ui\\build
else
	@mkdir -p ${STATIC_DIR}/build
	@cd ui && yarn build
	@mv ui/build/* ${STATIC_DIR}/build
	@rm -rf ui/build
endif
	@echo "Successfully built the frontend"
	
build: 
	@echo "Buidling binaries..."
ifeq ($(OS),Windows_NT) # Windows
	$ go build ${LDFLAGS} -o $(BIN_DIR)\\sqlweb .\\cmd\\main
else
	$ go build ${LDFLAGS} -o $(BIN_DIR)/sqlweb ./cmd/main
endif
	@echo "Successfully built sqlweb"

build_linux_arm64:
	@echo "Buidling for linux-arm64...."
	$ GOOS=linux GOARCH=arm64 go build ${LDFLAGS} -o ${BIN_DIR}/linux-arm-64-bit/sqlweb ./cmd/main
	@sleep 2

build_linux_amd64:
	@echo "Buidling for linux-amd64...."
	$ GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o ${BIN_DIR}/linux-amd-64-bit/sqlweb ./cmd/main
	@echo "Successfully built for linux-amd64"

build_mac_amd64:
	@echo "Buidling for mac-amd64...."
	$ GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o ${BIN_DIR}/mac-amd-64-bit/sqlweb ./cmd/main
	@echo "Successfully built for mac-amd64"

build_mac_arm64:
	@echo "Buidling for mac-arm64...."
	$ GOOS=darwin GOARCH=arm64 go build ${LDFLAGS} -o ${BIN_DIR}/mac-arm-64-bit/sqlweb ./cmd/main
	@echo "Successfully built for mac-arm64"

build_win_amd64:
	@echo "Buidling for windows-amd64...."
	$ set GOOS=windows GOARCH=amd64
	$ go build ${LDFLAGS} -o ${BIN_DIR}\\windows-amd-64-bit\\sqlweb.exe .\\cmd\\main
	@echo "Successfully built for windows-amd64"

install-bin:
ifeq ($(OS),Windows_NT) # Windows
	@echo "Moving binaries to $(USERPROFILE)\\sqlweb ...."
	@rm -rf $(USERPROFILE)\\sqlweb
	@mkdir $(USERPROFILE)\\sqlweb
	@mv ${BIN_DIR}\\windows-amd-64-bit\\sqlweb.exe $(USERPROFILE)\\sqlweb
else
	@echo "Moving binaries to /usr/local/bin...."
	@mv ${BIN_DIR}/sqlweb /usr/local/bin/
endif
	@echo "Successfully installed sqlweb"

uninstall-bin:
ifeq ($(OS),Windows_NT) # Windows
	@echo "Removing binaries from $(USERPROFILE)\\sqlweb ...."
	@rm -rf $(USERPROFILE)\\sqlweb
else
	@echo "Removing binaries from /usr/local/bin...."
	@rm -rf /usr/local/bin/sqlweb
endif
	@echo "Successfully uninstalled sqlweb"

dev: clean build_ui build
clean: clean_ui clean_bin
build: clean build_ui build
release: clean_release clean_ui build_ui build_linux_arm64 build_linux_amd64 build_mac_amd64 build_mac_arm64 build_win_amd64 build_win_arm64
install: install-bin
uninstall: uninstall-bin