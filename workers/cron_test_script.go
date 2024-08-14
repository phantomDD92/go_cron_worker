package workers

import (
	"github.com/go-co-op/gocron"
	// "go_proxy_worker/models"
	// "go_proxy_worker/logger"
	// "go_proxy_worker/utils"
	// "go_proxy_worker/dbRedisQueries"
	// "github.com/fatih/structs"
	// "go_proxy_worker/db"
	// // "gorm.io/gorm"
	// "strconv"
	// "strings"
	"log"
	"time"
	// "fmt"
	// "os"
)

func CronTestScript() {
	s := gocron.NewScheduler(time.UTC)
	s.Every(60).Minutes().Do(RunTestScript)
	s.StartBlocking()
}

func RunTestScript() {

	// if utils.OnlyRunTestAccounts() {
	// 	log.Println("helloe")
	// }

	log.Println("helloe")

}
