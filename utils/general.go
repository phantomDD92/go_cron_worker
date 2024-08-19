package utils

import (
	"context"
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	//"log"
)

func GetRedisCtx() context.Context {
	// ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	// defer cancel()

	// select {
	// case <-ctx.Done():
	// 	if ctx.Err() != nil {
	// 		log.Println("Redis Context Errror:", ctx.Err())
	// 	}
	// }
	return context.TODO()
}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// Max returns the larger of x or y.
func MaxFloat32(x, y float32) float32 {
	if x < y {
		return y
	}
	return x
}

// Min returns the smaller of x or y.
func MinFloat32(x, y float32) float32 {
	if x > y {
		return y
	}
	return x
}

// Avg
func AvgInt(x, y int) int {
	if x == 0 && y == 0 {
		return 0
	}
	return (x + y) / 2
}

func AvgFloat32(x, y float32) float32 {
	if x == 0 && y == 0 {
		return 0
	}
	return (x + y) / 2
}

func AvgFloat64(x, y float64) float64 {
	if x == 0 && y == 0 {
		return 0
	}
	return (x + y) / 2
}

// Time
func MaxTime(x, y time.Time) time.Time {
	if x.After(y) {
		return x
	}
	return y
}

func ParseIntFromPattern(text string, pattern string, splitter ...string) (int, error) {
	// Compile the regular expression pattern
	re := regexp.MustCompile(pattern)
	// Find all matches of the pattern in the input string
	matches := re.FindAllStringSubmatch(text, -1)
	// Extract the first match (which should be the number)
	var result string
	if len(matches) > 0 {
		result = matches[0][1]
	} else {
		return 0, errors.New("pattern not found")
	}
	var split string
	if len(splitter) > 0 {
		split = splitter[0]
	} else {
		split = ","
	}
	num, err := strconv.Atoi(strings.ReplaceAll(result, split, ""))
	return num, err
}

func ParseTextFromPattern(text string, pattern string) (string, error) {
	// Compile the regular expression pattern
	re := regexp.MustCompile(pattern)
	// Find all matches of the pattern in the input string
	matches := re.FindAllStringSubmatch(text, -1)
	// Extract the first match (which should be the number)
	var result string
	if len(matches) > 0 {
		result = matches[0][1]
	} else {
		return "", errors.New("pattern not found")
	}
	return result, nil
}

func ExtractTextFromTag(textTag *goquery.Selection) string {
	text := textTag.Text()
	re := regexp.MustCompile(`\s+`)
	return re.ReplaceAllString(strings.TrimSpace(text), " ")
}
