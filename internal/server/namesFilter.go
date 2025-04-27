package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/brachiGH/firedns/internal/utils"
	"github.com/brachiGH/firedns/internal/utils/config"
	"github.com/brachiGH/firedns/monitor/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Helper function to convert labels to a reversed domain path string.
// e.g., ["www", "google", "com"] -> "com.google.www"
func labelsToReversedPath(labels []utils.Lable) string {
	parts := make([]string, len(labels))
	for i, label := range labels {
		// Assuming Lable is []byte or similar string-convertible type
		parts[i] = string(label)
	}
	// Reverse the parts
	for i, j := 0, len(parts)-1; i < j; i, j = i+1, j-1 {
		parts[i], parts[j] = parts[j], parts[i]
	}
	return strings.Join(parts, ".")
}

func CheckIfDomainIsBlocked(lables []utils.Lable, IP net.IP) bool {
	db, err := database.GetSettingsDB()
	if err != nil {
		log.Printf("Error accessing settings DB: %v. Allowing request.", err)
		return false // Fail-safe: allow if DB is unavailable
	}

	if len(lables) == 0 {
		return false // Cannot check an empty domain
	}

	ipStr := IP.String() // Use the standard string representation of the IP
	// Replace dots with a character allowed in MongoDB keys if necessary,
	// or ensure your storage mechanism handles dots in keys.
	// For simplicity, assuming dots are handled or replaced elsewhere if needed.
	// If using integer keys like the example '1278264', convert IP accordingly.

	domainPath := labelsToReversedPath(lables)
	if domainPath == "" {
		return false
	}

	// Construct the query path, e.g., "127.0.0.1.com.google.www"
	// This assumes the document structure is { <ip_string>: { <tld>: { <domain>: ... } } }
	// And the IP string itself is the top-level key in the document.
	// If the top-level key is literally "ip", the path would be "ip." + domainPath
	// Adjust based on your exact BSON structure. Using the structure from the example:
	// { "1278264": { "com": { "google": { "www": 0 } } } }
	// We need the IP representation used as the key (e.g., "1278264" or "127.0.0.1")
	// Assuming ipStr is the correct key format used in the DB.
	queryField := fmt.Sprintf("%s.%s", ipStr, domainPath)

	// Filter to check if the specific nested field exists for the given IP key.
	// We look for the document matching the IP key, and within it, the nested domain path.
	filter := bson.M{queryField: bson.M{"$exists": true}}

	// Check if a document matching the filter exists
	var result bson.M
	err = db.BlockedDomainsList.FindOne(context.Background(), filter).Decode(&result)

	if err == mongo.ErrNoDocuments {
		// No document found where this specific IP and domain path exists.
		// Therefore, the domain is not blocked for this IP.
		return false
	} else if err != nil {
		// An error occurred during the database query
		log.Printf("Error querying blocklist for IP '%s', domain path '%s': %v. Allowing request.", ipStr, domainPath, err)
		return false // Fail-safe: allow if query fails
	}

	// A matching document/path was found, the domain is blocked for this IP.
	return true
}

func CreateBlockedDomainDNSMessage(data []byte, hdr *DNSHeader, q *DNSQuestion) *DNSMessage {
	arr := &DNSAnswer{
		name:  []byte{0xC0, 0x0C},
		typ:   1,
		class: 1,
		ttl:   60,
		rdlen: 4,
		rdata: []byte{0x00, 0x00, 0x00, 0x00},
	}

	hdr.ancount = 1
	hdr.arcount = 1

	rd := hdr.flags.RD()

	hdr.flags = NewDNSFlag()
	hdr.flags.SetQR(1)
	hdr.flags.SetRD(rd)
	hdr.flags.SetRA(1)

	dnsMessage := &DNSMessage{
		hdr:          hdr,
		q:            data[q.labelsStartPointer : q.labelsEndPointer+1+QueryInfomationBytesLength],
		AnswerRR:     arr,
		AdditionalRR: config.UDP_Response_Additional_Records,
	}

	return dnsMessage
}
