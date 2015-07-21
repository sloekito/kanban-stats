package main

import (
	"fmt"
	"log"
)


func main(){
	log.Print("---daily_trello: start")
	
	client := Client{
		Key: KEY,
		Token: TOKEN,
	}	
	lists := client.GetTrelloLists("534efa63a3a33edc034ac3d1")
	
	for _, list := range lists {
		fmt.Printf("%s(%s): %d\n", list.Name, list.Id, len(list.Cards))
	}

	PublishListsToInflux("DevProcess", lists)
	
	
	log.Print("---daily_trello: success")	
}