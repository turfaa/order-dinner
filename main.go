package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/turfaa/order-dinner/dinner"
	"github.com/turfaa/order-dinner/service"
)

func main() {
	_ = godotenv.Load()

	token, ok := os.LookupEnv("TOKEN")
	if !ok {

		log.Fatalln("No TOKEN")
	}

	sfids, ok := os.LookupEnv("FOOD_IDS")
	if !ok {
		log.Fatalln(("No FOOD_IDS"))
	}

	afids := strings.Split(sfids, ",")

	fids := make([]int, 0, len(afids))
	for _, sfid := range afids {
		if fid, err := strconv.Atoi(sfid); err != nil {
			log.Fatalf("Error converting food id: %s", err.Error())
		} else {
			fids = append(fids, fid)
		}
	}

	ctx := context.Background()
	c, err := dinner.NewDinnerClient(ctx, "https://dinner.seagroup.com/api", token, fids)
	if err != nil {
		log.Fatalf("Error creating dinner client: %s", err.Error())
	}

	s, err := service.NewDinnerService(ctx, c, 1000, time.Now().UTC().Add(time.Hour*11+time.Minute*31))
	if err != nil {
		log.Fatalf("Error creating dinner service: %s", err.Error())
	}

	log.Print("Starting dinner service")
	err = s.Serve()
	log.Printf("Dinner service closed: %s", err)
}
