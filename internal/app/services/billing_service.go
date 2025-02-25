package services

import (
	"encoding/csv"
	"fmt"
	"mass-transit-billing/internal/app/domain"
	"math"
	"os"
	"sort"
	"strconv"
	"time"
)

type BillingService struct {
	zoneMapFile     string
	journeyDataFile string
	outputFile      string
	stationToZone   map[string]int
	journeysPerUser map[string][]*domain.Journey
	userBills       map[string]float64
}

func NewBillingService(zoneMapFile, journeyDataFile, outputFile string) *BillingService {
	return &BillingService{
		zoneMapFile:     zoneMapFile,
		journeyDataFile: journeyDataFile,
		outputFile:      outputFile,
		stationToZone:   make(map[string]int),
		journeysPerUser: make(map[string][]*domain.Journey),
		userBills:       make(map[string]float64),
	}
}

func (s *BillingService) getAdditionalZoneCost(zone int) (float64, error) {
	switch {
	case zone == 1:
		return 0.80, nil
	case zone >= 2 && zone <= 3:
		return 0.50, nil
	case zone >= 4 && zone <= 5:
		return 0.30, nil
	case zone >= 6:
		return 0.10, nil
	default:
		return 0, fmt.Errorf("invalid zone: %d", zone)
	}
}

func (s *BillingService) readZoneMap() error {
	file, err := os.Open(s.zoneMapFile)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	// Skipping header
	_, err = reader.Read()
	if err != nil {
		return err
	}

	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	for _, record := range records {
		zone, _ := strconv.Atoi(record[1])
		s.stationToZone[record[0]] = zone
	}
	return nil
}

func (s *BillingService) readJourneyData() error {
	file, err := os.Open(s.journeyDataFile)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	// Skipping header
	_, err = reader.Read()
	if err != nil {
		return err
	}

	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	for _, record := range records {
		userId := record[0]
		stationName := record[1]
		direction := record[2]
		dateAndTimeRaw := record[3]

		// converting station to zone
		zone := s.stationToZone[stationName]

		// converting timestamp to time.Time
		timestamp, err := time.Parse("2006-01-02T15:04:05", dateAndTimeRaw)
		if err != nil {
			return err
		}

		//creating new journey record
		journey := domain.NewJourney(zone, direction, timestamp)

		// creating a map of journeys per user
		s.journeysPerUser[userId] = append(s.journeysPerUser[userId], journey)
	}

	return nil
}

func (s *BillingService) calculateBillingAmountsPerUser() error {
	for userID, journeys := range s.journeysPerUser {
		var dailyJourneys []int
		var dailyTotal, monthlyTotal float64
		curDay := journeys[0].GetTimestamp().Day()

		for _, journey := range journeys {
			journeyDay := journey.GetTimestamp().Day()
			if journeyDay != curDay {
				monthlyTotal += math.Min(15.00, dailyTotal)
				dailyTotal = 0.00
				curDay = journeyDay
			}

			if journey.GetDirection() == "IN" {
				dailyJourneys = append(dailyJourneys, journey.GetZone())
			} else if journey.GetDirection() == "OUT" {
				if len(dailyJourneys) == 0 {
					dailyTotal += 5.00
				} else {
					entryZone := dailyJourneys[len(dailyJourneys)-1]
					dailyJourneys = dailyJourneys[:len(dailyJourneys)-1]
					entryCost, err := s.getAdditionalZoneCost(entryZone)
					if err != nil {
						return fmt.Errorf("error getting additional zone cost: %v", err)
					}
					exitCost, err := s.getAdditionalZoneCost(journey.GetZone())
					if err != nil {
						return fmt.Errorf("error getting additional zone cost: %v", err)
					}
					journeyCost := 2.0 + entryCost + exitCost
					dailyTotal += journeyCost
				}
			}
		}

		if len(dailyJourneys) > 0 {
			dailyTotal += float64(len(dailyJourneys)) * 5.00
		}
		monthlyTotal += math.Min(15.00, dailyTotal)
		s.userBills[userID] = math.Min(100.00, monthlyTotal)
	}
	return nil
}

func (s *BillingService) writeBillingOutput() error {
	file, err := os.Create(s.outputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	userIDs := make([]string, 0, len(s.userBills))
	for userID := range s.userBills {
		userIDs = append(userIDs, userID)
	}

	//Sorting by userids
	sort.Strings(userIDs)

	for _, userID := range userIDs {
		amount := fmt.Sprintf("%.2f", s.userBills[userID])
		writer.Write([]string{userID, amount})
	}
	return nil
}

func (s *BillingService) ProcessBilling() error {
	if err := s.readZoneMap(); err != nil {
		return err
	}
	if err := s.readJourneyData(); err != nil {
		return err
	}
	//calculating the amounts per user
	s.calculateBillingAmountsPerUser()

	//writting output
	err := s.writeBillingOutput()
	if err != nil {
		return err
	}
	return nil
}
