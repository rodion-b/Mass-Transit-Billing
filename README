
# Mass Transit Billing

## Overview

A billing system for mass transit of passengers on a network of train stations, each of which belongs to a pricing zone.

The system parses CSV inputs, such as journey data and zone maps and calculates travel costs per passenger. The resulting output is provided in CSV format in a desirable directory.

### Assumptions

- CSV files are well formatted: no missing data, timestamps are in UTC and correct.
- There is at least one journey in input file.
- Prices of the zones are hardcoded for simplicity as there was no requirments to keep them dynamic.
- At the moment supports correct calculation of monthly batches

## Technical requirements
- **Language:** Golang  
- No External libraries 

### Running the Program
To run the program, use the below command providing the correct file paths:

```bash
go run main.go path_to_zone_map.csv path_to_journey_data.csv path_to_output.csv
```

To run all tests, simply use the below command:

```bash
go test ./...
```
