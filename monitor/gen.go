package monitor

//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -tags linux -cc clang -cflags "-I/usr/include/bpf" nic_monitor XDP/nic_monitor.c
