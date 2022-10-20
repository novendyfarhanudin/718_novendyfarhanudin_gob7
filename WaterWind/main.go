package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"example.id/waterwind/models"
	"github.com/gin-gonic/gin"
)

func main() {
	go jsonUpdate()
	router := gin.Default()
	router.Static("/assets", "./assets")
	router.LoadHTMLGlob("templates/*.html")
	router.GET("/", func(c *gin.Context) {
		jsonFile, err := os.Open("waterwind.json")
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			log.Println(err.Error())
			return
		}
		bytes, err := ioutil.ReadAll(jsonFile)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			log.Println(err.Error())
			return
		}
		var status models.Status
		json.Unmarshal(bytes, &status)
		log.Printf("in router: %+v\n", status)
		waterStatus, windStatus, waterClass, windClass := "", "", "", ""
		switch {
		case status.Status.Water <= 5:
			waterStatus = "Safe"
			waterClass = "list-group-item-success"
		case status.Status.Water <= 8:
			waterStatus = "Alert"
			waterClass = "list-group-item-warning"
		default:
			waterStatus = "Danger"
			waterClass = "list-group-item-danger"
		}
		switch {
		case status.Status.Wind <= 6:
			windStatus = "Safe"
			windClass = "list-group-item-success"
		case status.Status.Wind <= 15:
			windStatus = "Alert"
			windClass = "list-group-item-warning"
		default:
			windStatus = "Danger"
			windClass = "list-group-item-danger"
		}
		c.HTML(http.StatusOK, "index.html", gin.H{
			"water":       status.Status.Water,
			"waterStatus": waterStatus,
			"waterClass":  waterClass,
			"wind":        status.Status.Wind,
			"windStatus":  windStatus,
			"windClass":   windClass,
		})
	})
	router.Run("localhost:8080")
}

func jsonUpdate() {
	for {
		waterStatus, windStatus := rand.Intn(3), rand.Intn(3)
		var water, wind int
		switch waterStatus {
		case 0:
			water = rand.Intn(5) + 1
		case 1:
			water = rand.Intn(3) + 6
		default:
			water = rand.Intn(92) + 9
		}
		switch windStatus {
		case 0:
			wind = rand.Intn(6) + 1
		case 1:
			wind = rand.Intn(9) + 7
		default:
			wind = rand.Intn(85) + 16
		}
		jsonString := fmt.Sprintf(`{
	"status": {
		"water": %d,
		"wind": %d
	}
}`, water, wind)
		jsonFile, err := os.Create("waterwind.json")
		if err != nil {
			log.Println("Error in jsonUpdate:", err.Error())
			continue
		}
		_, err = jsonFile.Write([]byte(jsonString))
		if err != nil {
			log.Println("Error in jsonUpdate:", err.Error())
			continue
		}
		log.Printf("in jsonUpdater: {Status:{Water:%d Wind:%d}}\n", water, wind)
		time.Sleep(15 * time.Second)
	}
}