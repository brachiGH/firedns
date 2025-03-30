// go:build ignore

#include <linux/bpf.h>
#include <bpf/bpf_helpers.h>
#include <bpf_endian.h>
#include <linux/if_ether.h>
#include <linux/ip.h>
#include <linux/udp.h>
#include <linux/tcp.h>
#include <netinet/in.h>
#include "common.h"

struct
{
  __uint(type, BPF_MAP_TYPE_HASH);
  __type(key, __be32);
  __type(value, __u32);
  __uint(max_entries, MAX_ENTRIES_PER_TICK);
} query_count_per_ip SEC(".maps");

struct
{
  __uint(type, BPF_MAP_TYPE_HASH);
  __type(key, __be32);
  __type(value, bool);
  __uint(max_entries, MAX_ENTRIES_LOADED_IN_RAM);
} premium_ips SEC(".maps");

struct
{
  __uint(type, BPF_MAP_TYPE_HASH);
  __type(key, __be32);
  __type(value, bool);
  __uint(max_entries, MAX_ENTRIES_LOADED_IN_RAM);
} usage_limit_exeded_ips SEC(".maps");

// count response queries per ip (ingress)
SEC("xdp")
int query_analyser(struct xdp_md *ctx)
{
  void *data_end = (void *)(unsigned long)ctx->data_end;
  void *data = (void *)(unsigned long)ctx->data;

  // Boundary check: check if packet is larger than a full ethernet + ip header
  if (data + sizeof(struct ethhdr) + sizeof(struct iphdr) > data_end)
  {
    return DEFAULT_ACTION;
  }

  struct ethhdr *eth = data;

  // Ignore packet if ethernet protocol is not IP-based
  if (eth->h_proto != bpf_htons(ETH_P_IP))
  {
    return DEFAULT_ACTION;
  }

  struct iphdr *ip = data + sizeof(*eth);

  if (ip->version != 4)
  {
    return XDP_DROP; // IPv6 packets not yet supported
  }

  if (ip->protocol == IPPROTO_UDP)
  {
    // Boundary check for UDP
    if (data + sizeof(*eth) + sizeof(*ip) + sizeof(struct udphdr) > data_end)
    {
      return DEFAULT_ACTION;
    }

    struct udphdr *udp = data + sizeof(*eth) + sizeof(*ip);

    // Check if packet is a dns message from a user
    if (udp->dest == bpf_htons(FIREDNS_DNS_OVER_UDP_PORT))
    {
      goto record_ip;
    }
  }

  if (ip->protocol == IPPROTO_TCP)
  {
    // Boundary check for TCP
    if (data + sizeof(*eth) + sizeof(*ip) + sizeof(struct tcphdr) > data_end)
    {
      return DEFAULT_ACTION;
    }

    struct tcphdr *tcp = data + sizeof(*eth) + sizeof(*ip);

    // Check if it's your DNS over TCP port
    if (tcp->dest == bpf_htons(FIREDNS_DNS_OVER_TCP_PORT) ||
        tcp->dest == bpf_htons(FIREDNS_DNS_OVER_HTTPS_PORT) ||
        tcp->dest == bpf_htons(FIREDNS_DNS_OVER_TLS_PORT))
    {
      goto record_ip;
    }
  }

  return DEFAULT_ACTION;

record_ip:
{
  __be32 key = ip->addrs.saddr;
  bool *is_usage_exeded = bpf_map_lookup_elem(&usage_limit_exeded_ips, &key);
  if (is_usage_exeded)
  {
    return XDP_DROP;
  }

  __u32 *question_count = bpf_map_lookup_elem(&query_count_per_ip, &key);

  if (question_count)
  {
    // drop users query if the are more MAX_QUERIES_PER_TICK
    // eliminate the risk of dos attack
    __sync_fetch_and_add(question_count, 1);
    if (*question_count > MAX_QUERIES_PER_TICK)
    {
      bool *is_premium_user = bpf_map_lookup_elem(&premium_ips, &key);
      if (is_premium_user)
      {
        return XDP_PASS;
      }
      return XDP_DROP;
    }
  }
  else
  {
    __u32 init_count = 1;
    int err = bpf_map_update_elem(&query_count_per_ip, &key, &init_count, BPF_ANY);
    if (err != 0)
    {
      return XDP_DROP;
    }
  }

  return DEFAULT_ACTION;
}
}

char __license[] SEC("license") = "Dual MIT/GPL";