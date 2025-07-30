package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"math/rand"
)

type MatchCount map[int]int

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Generating test file...")
		// Append 5 random numbers per line between 1 and 90 to test_players.txt until 1 million lines are reached
		fileName := os.Args[1]
		lineCount := 1000000
		numsPerLine := 5
		min, max := 1, 90

		// Open file in append mode, create if it doesn't exist
		file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println("Error opening file:", err)
			return
		}
		defer file.Close()

		rand.Seed(time.Now().UnixNano())

		for i := 0; i < lineCount; i++ {
			line := ""
			for j := 0; j < numsPerLine; j++ {
				n := rand.Intn(max-min+1) + min
				line += strconv.Itoa(n)
				if j < numsPerLine-1 {
					line += " "
				}
			}
			_, err := file.WriteString(line + "\n")
			if err != nil {
				fmt.Println("Error writing to file:", err)
				return
			}
		}

		fmt.Println("Appended", lineCount, "lines to", fileName)

		return
	}

	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <input_file> <threads_count>")
		return
	}

	// Read players from file
	players, err := readPlayers(os.Args[1])
	if err != nil {
		fmt.Println("Error reading player file:", err)
		return
	}

	// Get winning numbers
	winning := getWinningNumbers()

	start := time.Now()
	result := MatchCount{2: 0, 3: 0, 4: 0, 5: 0}

	// Count matches
	for _, player := range players {
		match := countMatches(player, winning)
		if match >= 2 && match <= 5 {
			result[match]++
		}
	}

	// Print results
	fmt.Println("Number Matching | Winners")
	fmt.Println("----------------|--------")
	for i := 5; i >= 2; i-- {
		fmt.Printf("%-16d| %d\n", i, result[i])
	}

	elapsed := time.Since(start)
	fmt.Printf("Execution took %s (%d ns)\n", elapsed, elapsed.Nanoseconds())
}

func readPlayers(path string) ([][]int, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var players [][]int
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		nums := strings.Fields(line)
		if len(nums) != 5 {
			continue
		}
		var picks []int
		for _, s := range nums {
			n, err := strconv.Atoi(s)
			if err != nil {
				return nil, fmt.Errorf("invalid number: %v", s)
			}
			picks = append(picks, n)
		}
		players = append(players, picks)
	}
	return players, scanner.Err()
}

func getWinningNumbers() []int {
	fmt.Print("Enter 5 winning numbers (space-separated): ")
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	fields := strings.Fields(line)

	if len(fields) != 5 {
		fmt.Println("Please enter exactly 5 numbers.")
		os.Exit(1)
	}

	var win []int
	for _, s := range fields {
		n, err := strconv.Atoi(s)
		if err != nil {
			fmt.Println("Invalid number:", s)
			os.Exit(1)
		}
		win = append(win, n)
	}
	return win
}

func countMatches(player []int, winning []int) int {
	match := 0
	set := make(map[int]bool)
	for _, n := range winning {
		set[n] = true
	}
	for _, p := range player {
		if set[p] {
			match++
		}
	}
	return match
}
