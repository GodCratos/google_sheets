package main

import (
	"fmt"

	"github.com/GodCratos/google_sheets/services"
)

func main() {
	/*job := func() {
		now := time.Now()
		start := now.Format("2006-01-02 15:04:05")
		fmt.Println("Started job at : ", start)
		fmt.Println("Finished job at : ", time.Now().Format("2006-01-02 15:04:05"))
	}

	if _, err := scheduler.Every().Day().At("9:48").Run(job); err != nil {
		fmt.Println(err)
	}
	if _, err := scheduler.Every().Day().At("9:49").Run(job); err != nil {
		fmt.Println(err)
	}
	runtime.Goexit()*/
	err := services.GoogleSheetsWriteDataInSheet()
	if err != nil {
		fmt.Println(err)
	}

}
