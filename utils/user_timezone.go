package utils

import (
	"encoding/json"
	"fmt"
	"go_proxy_worker/models"
	"log"
	"strconv"
	"strings"
	"time"
)

// GetCurrentDateForOffset returns the current time for a given UTC offset
func GetCurrentDateTimeForOffset(accountItem models.Account) (time.Time, error) {

	// Get Account Timezone Offset
	timeZoneString := accountItem.Timezone
	timeZoneMap := make(map[string]string)

	err := json.Unmarshal([]byte(timeZoneString), &timeZoneMap)
	if err != nil {
		log.Println("ERROR - error converting timezone string to map")
		return time.Time{}, fmt.Errorf("ERROR - error converting timezone string to map")
	}

	// adjustedJobStartTime := jobStartTime
	offsetString, present := timeZoneMap["timeValue"]

	if present {

		// Ensure the offset string has correct format and split it
		if len(offsetString) != 6 || (offsetString[0] != '+' && offsetString[0] != '-') {
			return time.Time{}, fmt.Errorf("invalid offset format")
		}

		// Extract and convert hours and minutes from the offset string
		hours, err1 := strconv.Atoi(offsetString[:3]) // Extracts "+hh" or "-hh"
		if err1 != nil {
			return time.Time{}, fmt.Errorf("invalid hour format: %v", err1)
		}

		// Skip the colon and convert the minutes
		minutes, err2 := strconv.Atoi(offsetString[4:]) // Skips the colon, extracts "mm"
		if err2 != nil {
			return time.Time{}, fmt.Errorf("invalid minute format: %v", err2)
		}

		fmt.Println("****** OFFSET CALC")
		fmt.Println(hours)
		fmt.Println(minutes)

		// Calculate the total offset in seconds
		totalOffset := hours*3600 + minutes*60
		if hours < 0 {
			// Adjust negative offset calculation because hours impact the sign of minutes too
			totalOffset = hours*3600 - minutes*60
		}

		// Create a location with the given offset
		location := time.FixedZone("", totalOffset)

		// Get the current time in the created location
		currentDateTime := time.Now().In(location)
		fmt.Println(currentDateTime)

		return currentDateTime, nil
	} else {

		currentDateTime := time.Now()
		return currentDateTime, nil
	}

}

func GetUsersTimezoneDayStart(jobStartTime time.Time, accountItem models.Account) time.Time {

	/*
		Get the dayStartDate for a user based on their timezone
	*/

	// Get Local Time
	localDateTime := jobStartTime

	// Get Account Timezone Offset
	timeZoneString := accountItem.Timezone
	timeZoneMap := make(map[string]string)

	err := json.Unmarshal([]byte(timeZoneString), &timeZoneMap)
	if err != nil {
		log.Println("ERROR - error converting timezone string to map")
	}

	adjustedJobStartTime := jobStartTime
	timeOffsetString, present := timeZoneMap["timeValue"]

	// Day Start Time
	dayStartTime := time.Date(adjustedJobStartTime.Year(), adjustedJobStartTime.Month(), adjustedJobStartTime.Day(), 0, 0, 0, 0, time.UTC)
	adjustedDayStartTime := dayStartTime

	if present {

		// Add Or Subtract Time
		firstCharacter := timeOffsetString[0:1]
		timeArray := strings.Split(timeOffsetString[1:], ":")
		hours, _ := strconv.Atoi(timeArray[0])
		minutes, _ := strconv.Atoi(timeArray[1])

		if firstCharacter == "-" {
			adjustedDayStartTime = adjustedDayStartTime.Add(time.Hour * time.Duration(hours))
			adjustedDayStartTime = adjustedDayStartTime.Add(time.Minute * time.Duration(minutes))

			// Localtime
			localDateTime = localDateTime.Add(-time.Hour * time.Duration(hours))
			localDateTime = localDateTime.Add(-time.Minute * time.Duration(minutes))

		} else if firstCharacter == "+" {
			adjustedDayStartTime = adjustedDayStartTime.Add(-time.Hour * time.Duration(hours))
			adjustedDayStartTime = adjustedDayStartTime.Add(-time.Minute * time.Duration(minutes))

			// Localtime
			localDateTime = localDateTime.Add(time.Hour * time.Duration(hours))
			localDateTime = localDateTime.Add(time.Minute * time.Duration(minutes))
		}

	}

	// Bug Fix --> If it has just crossed over the 24 hour mark
	dayDifference := adjustedDayStartTime.YearDay() - localDateTime.YearDay()
	if dayDifference < -1 {
		adjustedDayStartTime = adjustedDayStartTime.Add(time.Hour * time.Duration(24))
	}

	if dayDifference > 1 {
		adjustedDayStartTime = adjustedDayStartTime.Add(-time.Hour * time.Duration(24))
	}

	dayStartDateString := adjustedDayStartTime.Format("2006-01-02 15:04")
	dayStartDate, _ := time.Parse("2006-01-02 15:04", dayStartDateString)

	return dayStartDate.UTC()

}
