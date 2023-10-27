package main

import (
	"fmt"
	"strings"
	"os"
	"bufio"
)

const MENU_TEXT = `
======================================
Possible values for selection:
1 - do route lookup based on entered IP
2 - do route lookup based on entered Next-hop
3 - Print list of all unique Next-hops
9 - Exit the program
======================================`

// main menu func
func Menu() {
mainLoop:
	for {
		fmt.Println(MENU_TEXT)
		choise := requestUserInput("Enter your choise:")
		switch {
		case choise == "1":
			ip := requestUserInput("Enter IP:")
			routes, err := allRoutes.FindRoutes(ip, true)
			if err != nil {
				ErrorLogger.Printf("Cannot parse IP because of: %q", err)
			}
			fmt.Println("Matched routes:")
			for r := range routes {
				fmt.Println(r)
			}
		case choise == "2":
			nh := requestUserInput("Enter Next-hop value, either IP or interface format accepted:")
			routes := allRoutes.FindRoutesByNH(nh)
			fmt.Println("Matched routes:")
			for r := range routes {
				fmt.Println(r)
			}
		case choise == "3":
			nhs := allRoutes.FindUniqNexthops(false)
			fmt.Println("Unique nexthops:")
			for nh := range nhs {
				fmt.Println(nh)
			}
		case choise == "9":
			break mainLoop
		}
	}
}

// helper func to ask for user input
func requestUserInput(prompt string) string{
	fmt.Println(prompt)
	for {
		reader := bufio.NewReader(os.Stdin)
		line , err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Incorrect value entered, try again:")
			continue
		}
		input := strings.TrimSpace(line)
		return input
	}
}