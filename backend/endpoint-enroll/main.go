package main

import (
	"flag"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/stephenzsy/small-kms/backend/endpoint-enroll/tasks"
)

func main() {
	// Find .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	if len(os.Args) < 2 {
		log.Fatalf("Usage: enroll <task>")
		log.Fatalf("	install start			-- install starter credentials")
		log.Fatalf("	install view			-- view install certificate enrollment receipt")
		log.Fatalf("	install complete		-- complete starter certificate enrollment")
		log.Fatalf("	install complete user	-- complete starter certificate enrollment, install to user store")

		os.Exit(1)
	}
	receiptFileName := "enroll-receipt.json"

	flag.Parse()

	switch os.Args[1] {
	case "install":
		switch os.Args[2] {
		case "start":

			file, err := os.OpenFile(receiptFileName, os.O_CREATE|os.O_WRONLY, 0600)
			if err != nil {
				panic(err)
			}
			defer file.Close()
			if err := tasks.InstallStart(file, receiptFileName); err != nil {
				panic(err)
			}
		case "view":
			file, err := os.OpenFile(receiptFileName, os.O_RDONLY, 0600)
			if err != nil {
				panic(err)
			}
			if err := tasks.InstallView(file); err != nil {
				panic(err)
			}
		case "complete":
			file, err := os.OpenFile(receiptFileName, os.O_RDONLY, 0600)
			if err != nil {
				panic(err)
			}
			if err := tasks.InstallComplete(file, len(os.Args) >= 4 && (os.Args[3] == "user")); err != nil {
				panic(err)
			}
		}
	}
}
