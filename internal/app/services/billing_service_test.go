package services

import (
	"bufio"
	"encoding/csv"
	"mass-transit-billing/internal/app/domain"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)


func TestBillingServiceE2E(t *testing.T) {
	testOutputPath := filepath.Join("testdata", "test_output.csv")

	os.Remove(testOutputPath)

	service := NewBillingService(
		filepath.Join("testdata", "zone_map.csv"),
		filepath.Join("testdata", "journey_data.csv"),
		testOutputPath,
	)

	t.Run("Test Process Billing", func(t *testing.T) {
		err := service.ProcessBilling()
		if err != nil {
			t.Fatalf("Failed to process billing: %v", err)
		}

		// Validate the output
		expectedLines := []string{
			"user1,3.30",
			"user2,5.00",
			"user3,18.30",
		}

		file, err := os.Open(testOutputPath)
		if err != nil {
			t.Fatalf("Failed to open test output file: %v", err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		var outputLines []string
		for scanner.Scan() {
			outputLines = append(outputLines, strings.TrimSpace(scanner.Text()))
		}

		if err := scanner.Err(); err != nil {
			t.Fatalf("Error reading output file: %v", err)
		}

		if len(outputLines) != len(expectedLines) {
			t.Fatalf("Expected %d lines, but got %d", len(expectedLines), len(outputLines))
		}

		for i, expected := range expectedLines {
			if outputLines[i] != expected {
				t.Errorf("Mismatch on line %d: expected %q, got %q", i+1, expected, outputLines[i])
			}
		}
	})
}

func TestGetAdditionalZoneCost(t *testing.T) {
	service := NewBillingService("", "", "")
	tests := []struct {
		zone     int
		expected float64
		wantErr  bool
	}{
		{1, 0.80, false},
		{2, 0.50, false},
		{4, 0.30, false},
		{6, 0.10, false},
		{10, 0.10, false},
		{0, 0, true},
	}

	for _, tt := range tests {
		cost, err := service.getAdditionalZoneCost(tt.zone)
		if (err != nil) != tt.wantErr {
			t.Errorf("zone %d: expected error %v, got %v", tt.zone, tt.wantErr, err)
		}
		if !tt.wantErr && cost != tt.expected {
			t.Errorf("zone %d: expected cost %.2f, got %.2f", tt.zone, tt.expected, cost)
		}
	}
}

func TestCalculateBillingAmountsPerUser(t *testing.T) {
	service := NewBillingService("", "", "")

	// Mock journey data
	service.journeysPerUser = map[string][]*domain.Journey{
		"user1": {
			domain.NewJourney(1, "IN", time.Now()),
			domain.NewJourney(1, "OUT", time.Now().Add(1*time.Hour)),
		},
	}

	service.calculateBillingAmountsPerUser()

	if service.userBills["user1"] == 0 {
		t.Errorf("Expected non-zero bill for user1, got 0")
	}
}

func TestWriteBillingOutput(t *testing.T) {
	service := NewBillingService("", "", "test_output.csv")
	service.userBills = map[string]float64{
		"user1": 10.00,
		"user2": 15.00,
	}

	err := service.writeBillingOutput()
	if err != nil {
		t.Fatalf("Failed to write billing output: %v", err)
	}
	defer os.Remove("test_output.csv")

	file, err := os.Open("test_output.csv")
	if err != nil {
		t.Fatalf("Failed to open output file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	if len(records) != 2 {
		t.Errorf("Expected 2 records, got %d", len(records))
	}

	for _, record := range records {
		if _, err := strconv.ParseFloat(record[1], 64); err != nil {
			t.Errorf("Expected valid float value for amount, got %s", record[1])
		}
	}
}
