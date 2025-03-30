#include <stdint.h>
#define bool __u8

#define DEBUG
#define DEFAULT_ACTION XDP_PASS
#define FIREDNS_DNS_OVER_UDP_PORT 2053
#define FIREDNS_DNS_OVER_TCP_PORT 2053
#define FIREDNS_DNS_OVER_HTTPS_PORT 8443
#define FIREDNS_DNS_OVER_TLS_PORT 8853

#define MAX_ENTRIES_LOADED_IN_RAM 40000
#define MAX_ENTRIES_PER_TICK 10000
#define MAX_QUERIES_PER_TICK 256 	// Prevent DOS attack

