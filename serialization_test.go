package kubernetes

import (
	"bytes"
	"compress/gzip"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"testing"

	"go-micro.dev/v4/registry"
)

func TestSerializationSizes(t *testing.T) {
	// Create a dummy service struct
	service := &registry.Service{
		Name:    "test-service",
		Version: "1.0.0",
		Nodes: []*registry.Node{
			{
				Id:      "node-1",
				Address: "192.168.1.100:8080",
				Metadata: map[string]string{
					"env":     "production",
					"region":  "us-west-2",
					"version": "1.0.0",
				},
			},
			{
				Id:      "node-2", 
				Address: "192.168.1.101:8080",
				Metadata: map[string]string{
					"env":     "production",
					"region":  "us-west-2",
					"version": "1.0.0",
				},
			},
		},
	}

	// Test JSON serialization
	jsonData, err := json.Marshal(service)
	if err != nil {
		t.Fatalf("JSON marshal failed: %v", err)
	}

	// Test gob serialization
	var gobBuf bytes.Buffer
	gobEnc := gob.NewEncoder(&gobBuf)
	if err := gobEnc.Encode(service); err != nil {
		t.Fatalf("Gob encode failed: %v", err)
	}
	gobData := gobBuf.Bytes()

	// Test gzip+JSON compression
	var gzipBuf bytes.Buffer
	gzipWriter := gzip.NewWriter(&gzipBuf)
	if _, err := gzipWriter.Write(jsonData); err != nil {
		t.Fatalf("Gzip write failed: %v", err)
	}
	gzipWriter.Close()
	gzipData := gzipBuf.Bytes()

	// Test custom compact serialization
	compactData, err := compactEncode(service)
	if err != nil {
		t.Fatalf("Compact encode failed: %v", err)
	}

	// Test that compact decode works correctly
	decodedService, err := compactDecode(compactData)
	if err != nil {
		t.Fatalf("Compact decode failed: %v", err)
	}

	// Verify decoded service matches original
	if decodedService.Name != service.Name {
		t.Errorf("Name mismatch: got %s, want %s", decodedService.Name, service.Name)
	}
	if decodedService.Version != service.Version {
		t.Errorf("Version mismatch: got %s, want %s", decodedService.Version, service.Version)
	}
	if len(decodedService.Nodes) != len(service.Nodes) {
		t.Errorf("Node count mismatch: got %d, want %d", len(decodedService.Nodes), len(service.Nodes))
	}

	// Print size comparison
	fmt.Printf("\n=== Serialization Size Comparison ===\n")
	fmt.Printf("JSON:      %d bytes\n", len(jsonData))
	fmt.Printf("Gob:       %d bytes\n", len(gobData))
	fmt.Printf("Gzip+JSON: %d bytes\n", len(gzipData))
	fmt.Printf("Compact:   %d bytes\n", len(compactData))
	fmt.Printf("\nGzip+JSON vs JSON: %.1f%% savings\n", float64(len(jsonData)-len(gzipData))/float64(len(jsonData))*100)
	fmt.Printf("Compact vs JSON:   %.1f%% savings\n", float64(len(jsonData)-len(compactData))/float64(len(jsonData))*100)

	// Show actual serialized data (truncated for readability)
	fmt.Printf("\n=== Serialized Data Examples ===\n")
	fmt.Printf("JSON (first 100 chars): %s...\n", string(jsonData)[:min(100, len(jsonData))])
	fmt.Printf("Gob (first 50 bytes):   %q...\n", string(gobData[:min(50, len(gobData))]))
	fmt.Printf("Compact (first 50 bytes): %q...\n", string(compactData[:min(50, len(compactData))]))

	// Ensure compact format is significantly smaller
	if len(compactData) >= len(jsonData) {
		t.Errorf("Compact format should be smaller than JSON: compact=%d bytes, json=%d bytes", len(compactData), len(jsonData))
	}
	if len(compactData) >= len(gobData) {
		t.Errorf("Compact format should be smaller than gob: compact=%d bytes, gob=%d bytes", len(compactData), len(gobData))
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
