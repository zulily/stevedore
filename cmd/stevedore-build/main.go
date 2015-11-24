package main

import (
	"log"

	"github.com/zulily/stevedore"
	"github.com/zulily/stevedore/cmd"
)

func main() {
	images, err := stevedore.FindImagesInCwd(cmd.Filter)
	if err != nil {
		log.Fatal("error finding images:", err)
	}
	for _, img := range images {
		log.Println("Building", img)
		if err = img.Build(); err != nil {
			log.Println("error building", img)
			log.Fatal(err)
		}
	}

	log.Println("Done")
}
