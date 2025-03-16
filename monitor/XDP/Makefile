##*****
# this make file is for testing
##*****

CLANG = clang
#DEBUG = y  enables printk in the BPF program
DEBUG ?= n

#Path to include files of respectively libbpf-dev
LIBBPF_INCLUDE ?= /usr/include/bpf

CFLAGS ?= -I$(LIBBPF_INCLUDE) -I $(IPROUTE_INCLUDE) -g -static


BPF_CFLAGS ?= -I $(LIBBPF_INCLUDE) \
		-I $(IPROUTE_INCLUDE)

ifeq ($(DEBUG),y)
	BPF_EXTRA_FLAGS += -D DEBUG
endif

all: nic_monitor

.PHONY: clean $(CLANG)

clean:
	rm -f *.o

nic_monitor: %: %.c Makefile
	$(CLANG) \
	    -target bpf \
	    $(BPF_CFLAGS) \
	    -Wall $(BPF_EXTRA_FLAGS) \
            -O2 -c -o $@.o $<
