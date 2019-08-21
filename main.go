package main

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/turfaa/order-dinner/dinner"
	"github.com/turfaa/order-dinner/service"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	_ = godotenv.Load()

	token, ok := os.LookupEnv("TOKEN")
	if !ok {

		log.Fatalln("No TOKEN")
	}
	smid, ok := os.LookupEnv("MENU_ID")
	if !ok {
		log.Fatalln("No MENU_ID")
	}

	sfids, ok := os.LookupEnv("FOOD_IDS")
	if !ok {
		log.Fatalln(("No FOOD_IDS"))
	}

	mid, err := strconv.Atoi(smid)
	if err != nil {
		log.Fatalf("Error converting menu id: %s", err.Error())
	}

	afids := strings.Split(sfids, ",")

	fids := make([]int, len(afids))
	for _, sfid := range afids {
		if fid, err := strconv.Atoi(sfid); err != nil {
			log.Fatalf("Error converting food id: %s", err.Error())
		} else {
			fids = append(fids, fid)
		}
	}

	ctx := context.Background()
	c, err := dinner.NewDinnerClient(ctx, "https://dinner.seagroup.com/api", token, mid, fids)
	if err != nil {
		log.Fatalf("Error creating dinner client: %s", err.Error())
	}

	s, err := service.NewDinnerService(ctx, c, 1000, time.Date(2019, 8, 22, 0, 0, 0, 0, time.UTC))
	if err != nil {
		log.Fatalf("Error creating dinner service: %s", err.Error())
	}

	log.Print("Starting dinner service")
	err = s.Serve()
	log.Printf("Dinner service closed: %s", err)
}
