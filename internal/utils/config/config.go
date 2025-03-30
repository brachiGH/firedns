package config

import "time"

const UDP_ns_addr = "45.90.28.128:53" //name server
const UDP_local_ns_addr = "127.0.0.1:2053"
const TCP_local_ns_addr = "127.0.0.1:2053"

const NO_PROFILE_ID = ""

const ClearCache__TickDuration = time.Minute
const UpdateQuestions__TickDuration = 30 * time.Second
const UpdateUsageLimitIps__TickDuration = time.Minute * 10
const UpdatePremiumIps__TickDuration = time.Minute * 10
const NICMonitor__TickDuration = time.Minute

var UDP_Response_Additional_Records = []byte{0x0, 0x0, 0x29, 0x4, 0xd0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x18, 0x0, 0xf, 0x0, 0x14, 0x0, 0x11, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x65, 0x64, 0x20, 0x62, 0x79, 0x20, 0x46, 0x69, 0x72, 0x65, 0x44, 0x4e, 0x53}
