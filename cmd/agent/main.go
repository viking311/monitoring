package main

import (
	"time"

	"github.com/viking311/monitoring/internal/entity"
)

func main() {
	collector := entity.NewCollector("http://localhost:8080/", 2*time.Second, 10*time.Second)

	collector.Do()

}
