
package utils

import (
	"fmt"
	"time"
)



func GetTimeWindow() string {
	now := time.Now().UTC()
	return fmt.Sprintf("%v", now.Year()) + "-" + fmt.Sprintf("%v", now.Month()) + "-" + fmt.Sprintf("%v", now.Day()) + "::" + fmt.Sprintf("%v", now.Hour())
}

func GetDayString() string {
	now := time.Now().UTC()
	return fmt.Sprintf("%v", now.Year()) + "-" + fmt.Sprintf("%v", now.Month()) + "-" + fmt.Sprintf("%v", now.Day())
}

func ConvertMonthStringToInt(m string) int {

	switch m {
    case "January":
        return 1
	case "February":
        return 2
	case "March":
        return 3
	case "April":
        return 4
	case "May":
        return 5
	case "June":
        return 6
	case "July":
        return 7
	case "August":
        return 8
	case "September":
        return 9
	case "October":
        return 10
	case "November":
        return 11
	case "December":
        return 12

    }

	return 0
}
