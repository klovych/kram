
BINARY_NAME=kram


INSTALL_DIR=/usr/local/bin


build:
	@echo "Building the project..."
	go build -o $(BINARY_NAME) kram.go


install: build
	@echo "Installing $(BINARY_NAME) to $(INSTALL_DIR)..."
	sudo cp $(BINARY_NAME) $(INSTALL_DIR)


uninstall:
	@echo "Uninstalling $(BINARY_NAME)..."
	sudo rm -f $(INSTALL_DIR)/$(BINARY_NAME)


clean:
	@echo "Cleaning up..."
	rm -f $(BINARY_NAME)


.PHONY: build install uninstall clean
