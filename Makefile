GO = go
XDP_FOLDER = monitor

all: main

.PHONY: clean $(GO)

clean:
	rm -f $(XDP_FOLDER)/*.o $(XDP_FOLDER)/*_bpfeb.go $(XDP_FOLDER)/*_bpfel.go main

$(XDP_FOLDER)/nic_monitor_bpfeb.go:$(XDP_FOLDER)/gen.go
	cd $(XDP_FOLDER)/ && $(GO) generate

main: cmd/main.go $(XDP_FOLDER)/nic_monitor_bpfeb.go
	$(GO) build cmd/main.go

cleanmake:
	$(MAKE) clean
	$(MAKE)