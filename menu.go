package main

import (
	"fmt"
	"strings"
	"os"
	"bufio"
	"github.com/bondar-aleksandr/cisco_route_parser/parser/entities"
)

type RT interface {
	FindRoutes(string, bool) (int, <-chan *entities.Route, error)
	FindRoutesByNH(string) (int, <-chan *entities.Route)
	FindUniqNexthops(bool) (int, <-chan *entities.NextHop)
	String() string
}

const MENU_TEXT = `
======================================
Possible values for selection:
1 - do route lookup based on entered IP
2 - do route lookup based on entered Next-hop
3 - print list of all unique Next-hops
8 - print raw routingTable object (for debug)
9 - exit the program
======================================`

// main menu func
func Menu(rt RT) {
mainLoop:
	for {
		fmt.Println(MENU_TEXT)
		choise := requestUserInput("Enter your choise:")
		switch {
		case choise == "1":
			ip := requestUserInput("Enter IP:")
			n, routes, err := rt.FindRoutes(ip, true)
			if err != nil {
				logger.Errorf("Cannot parse IP because of: %q", err)
			}
			fmt.Printf("Found %d routes:\n", n)
			for r := range routes {
				fmt.Println(r)
			}
		case choise == "2":
			nh := requestUserInput("Enter Next-hop value, either IP or interface format accepted:")
			n, routes := rt.FindRoutesByNH(nh)
			fmt.Printf("Found %d routes:\n", n)
			for r := range routes {
				fmt.Println(r)
			}
		case choise == "3":
			n, nhs := rt.FindUniqNexthops(false)
			fmt.Printf("Found %d unique nexthops:\n", n)
			for nh := range nhs {
				fmt.Println(nh)
			}
		case choise == "8":
			fmt.Println(rt)
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