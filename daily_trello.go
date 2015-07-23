package main

import (
	"fmt"
	"log"
	"flag"
)


func main(){
	log.Print("---daily_trello: start")
	
	trelloKey := flag.String("trellokey", "", "Trello application key")
	trelloToken :=flag.String("trellotoken", "", "Trello access token")
	trelloBoardID := flag.String("boardid", "", "Trello board ID")
	flag.Parse()
	
	client := Client{
		Key: *trelloKey,
		Token: *trelloToken,
	}
	fmt.Println(client)	
	lists := client.GetTrelloLists(*trelloBoardID)
	
	for _, list := range lists {
		fmt.Printf("%s(%s): %d\n", list.Name, list.Id, len(list.Cards))
	}

	PublishListsToInflux(*trelloBoardID, lists)
	
	
	log.Print("---daily_trello: success")	
}