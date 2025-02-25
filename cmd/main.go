package main

import (
	"fmt"
	"log"
	"mass-transit-billing/internal/app/services"
	"os"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}

func run() error {
	if len(os.Args) != 4 {
		return fmt.Errorf("Wrong args - please use like: go run mass_transit_billing.go <zone_map_file> <journey_map_file> <output_file>")
	}
	zoneMapFile, journeyDataFile, outputFile := os.Args[1], os.Args[2], os.Args[3]
	billingSystem := services.NewBillingService(zoneMapFile, journeyDataFile, outputFile)
	err := billingSystem.ProcessBilling()
	if err != nil {
		return fmt.Errorf("Error processing billing: %v", err)
	}
	return nil
}
