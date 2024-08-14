package utils

import (
    "time"
    "context"
    //"log"
)


func GetRedisCtx() context.Context{
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