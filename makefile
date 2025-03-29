sqlc:
	cd ./src/db && sqlc generate

BUILD_DIR = senatus
build:
	docker build --no-cache -t fedora-build .
	docker run --platform linux/amd64 -v $(shell pwd):/app fedora-build

build-docker: clean
	@echo ">>>>>>>>>>>> --------- Creating directory structure... --------- <<<<<<<<<<<<"
	@mkdir -p $(BUILD_DIR)/static

	@echo ">>>>>>>>>>>> --------- Building binary for Linux... --------- <<<<<<<<<<<<"
	@CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/main ./src/main.go

	@echo ">>>>>>>>>>>> --------- Copying static files... --------- <<<<<<<<<<<<"
	@cp -r ./static/* $(BUILD_DIR)/static/ 2>/dev/null || true

	@echo ">>>>>>>>>>>> --------- Creating zip archive... --------- <<<<<<<<<<<<"
	@zip -r $(BUILD_DIR).zip $(BUILD_DIR)

	# @rm -rf $(BUILD_DIR)
	@echo ">>>>>>>>>>>> --------- Build complete! Output: $(BUILD_DIR).zip --------- <<<<<<<<<<<<"

clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f $(BUILD_DIR).zip