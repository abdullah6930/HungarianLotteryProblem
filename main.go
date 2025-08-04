package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
	"math/rand"
)

type MatchCount map[int]int

// Player represents a player's numbers as 5 bytes
type Player [5]byte

// createPlayer creates a player from 5 numbers (1-90 each fits in a byte)
func createPlayer(numbers []int) Player {
	var player Player
	for i, num := range numbers {
		player[i] = byte(num)
	}
	return player
}

func main() {


	if len(os.Args) < 3 {
		fmt.Println("Generating test file...")
		// Append 5 random numbers per line between 1 and 90 to test_players.txt until 1 million lines are reached
		fileName := os.Args[1]
		lineCount := 10000000
		numsPerLine := 5
		min, max := 1, 90

		// Open file in append mode, create if it doesn't exist
		file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)

		if err != nil {
			fmt.Println("Error opening file:", err)
			return
		}
		defer file.Close()

		rand.Seed(time.Now().UnixNano())

		// get line count from file
		var i int
		i, err = getLineCount(file)
		if err != nil {
			fmt.Println("Error getting line count:", err)
			return
		}

		for i < lineCount {
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
			i++
		}

		fmt.Println("Appended", lineCount, "lines to", fileName)

		return
	}

	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <input_file> <threads_count>")
		return
	}

	// Read players from file
	threadsInt, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("Error converting threads count to integer:", err)
		return
	}
	
	// Convert to byte (max 256 threads)
	if threadsInt > 255 {
		fmt.Println("Warning: Thread count capped at 255")
		threadsInt = 255
	}
	threads := byte(threadsInt)

	cores := runtime.NumCPU()
	if cores > 255 {
		cores = 255
	}
	cpuCores := byte(cores)

	fmt.Println("CPU cores:", cpuCores)
	if threads > cpuCores {
		fmt.Printf("Warning: Using %d threads with only %d CPU cores may cause overhead.\n", threadsInt, cpuCores)
	}
	if threads < cpuCores {
		fmt.Printf("Tip: Consider using %d threads to match your CPU cores for optimal performance.\n", cpuCores)
	}
	
	fmt.Println("Reading players from file...")
	fmt.Println("Threads:", threads)
	
	start := time.Now()
	players, err := readPlayers(os.Args[1], threads)
	fmt.Println("Players read:", len(players))
	elapsed := time.Since(start)
	fmt.Printf("Reading players Execution took %s (%d ns)\n", elapsed, elapsed.Nanoseconds())
	
	if err != nil {
		fmt.Println("Error reading player file:", err)
		return
	}
	
	// Get winning numbers
	winning := getWinningNumbers()
	
	start = time.Now()
	// Parallel match counting for maximum speed
	fmt.Printf("Counting matches with %d threads...\n", threads)
	result := countMatchesParallel(players, winning, threads)

	// Print results
	fmt.Println("Number Matching | Winners")
	fmt.Println("----------------|--------")
	for i := 5; i >= 2; i-- {
		fmt.Printf("%-16d| %d\n", i, result[i])
	}

	elapsed = time.Since(start)
	fmt.Printf("Counting matches Execution took %s (%d ns)\n", elapsed, elapsed.Nanoseconds())
}

func readPlayers(path string, threads byte) ([]Player, error) {
	// First, count total lines in the file
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	
	fmt.Println("Counting total lines...")
	totalLines, err := countLines(file)
	if err != nil {
		file.Close()
		return nil, err
	}
	file.Close()
	
	fmt.Printf("Total lines in file: %d\n", totalLines)
	
	// Calculate lines per thread
	linesPerThread := totalLines / int(threads)
	if linesPerThread == 0 {
		linesPerThread = totalLines
		threads = 1
	}
	
	fmt.Printf("Lines per thread: %d\n", linesPerThread)

	// Channel to collect results from each thread
	resultsChan := make(chan []Player, threads)
	var wg sync.WaitGroup

	// Launch reader threads - each reads its assigned line range
	for threadID := byte(0); threadID < threads; threadID++ {
		wg.Add(1)
		go func(id byte) {
			defer wg.Done()
			
			// Calculate start and end line numbers for this thread
			startLine := int(id) * linesPerThread
			endLine := startLine + linesPerThread
			if id == threads - 1 {
				endLine = totalLines // Last thread reads until end of file
			}
			
			fmt.Printf("Thread %d: lines %d to %d\n", id, startLine, endLine-1)

			players, err := readFileLines(path, startLine, endLine)
			fmt.Println("Thread", id, "players:", len(players))
			if err != nil {
				fmt.Printf("Thread %d error: %v\n", id, err)
				resultsChan <- nil
				return
			}
			
			resultsChan <- players
		}(threadID)
	}

	// Wait for all threads and close results channel
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	// Combine results from all threads
	var allPlayers []Player
	totalPlayers := 0
	
	for threadResults := range resultsChan {
		if threadResults != nil {
			allPlayers = append(allPlayers, threadResults...)
			totalPlayers += len(threadResults)
		}
	}

	return allPlayers, nil
}



// Parallel match counting for maximum performance
func countMatchesParallel(players []Player, winning []int, threads byte) MatchCount {
	if len(players) == 0 {
		return MatchCount{2: 0, 3: 0, 4: 0, 5: 0}
	}

	totalPlayers := len(players)
	chunkSize := totalPlayers / int(threads)
	if chunkSize == 0 {
		chunkSize = 1
	}

	resultsChan := make(chan MatchCount, threads)
	var wg sync.WaitGroup

	for i := 0; i < int(threads); i++ {
		start := i * chunkSize
		end := start + chunkSize
		if i == int(threads)-1 {
			end = totalPlayers
		}
		if start >= totalPlayers {
			break
		}

		wg.Add(1)
		go func(playerChunk []Player) {
			defer wg.Done()
			localResult := MatchCount{2: 0, 3: 0, 4: 0, 5: 0}
			var set map[int]bool
			var match int
			var p, n int

			for _, player := range playerChunk {
				// Convert player bytes to int slice for matching
				numbers := make([]int, 5)
				for i := 0; i < 5; i++ {
					numbers[i] = int(player[i])
				}
				match = countMatches(numbers, winning, set, match, n, p)
				if match >= 2 && match <= 5 {
					localResult[match]++
				}
			}

			resultsChan <- localResult
		}(players[start:end])
	}

	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	finalResult := MatchCount{2: 0, 3: 0, 4: 0, 5: 0}
	for localResult := range resultsChan {
		for k, v := range localResult {
			finalResult[k] += v
		}
	}

	return finalResult
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

func countMatches(player []int, winning []int, set map[int]bool, match int, n int, p int) int {
	match = 0
	set = make(map[int]bool)
	n = 0
	p = 0
	for _, n = range winning {
		set[n] = true
	}
	for _, p = range player {
		if set[p] {
			match++
		}
	}
	return match
}

func getLineCount(file *os.File) (int, error) {
	// Reset file pointer to beginning
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return 0, err
	}

	scanner := bufio.NewScanner(file)
	lineCount := 0

	for scanner.Scan() {
		lineCount++
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}

	return lineCount, nil
}

// countLines efficiently counts total lines in a file
func countLines(file *os.File) (int, error) {
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return 0, err
	}

	scanner := bufio.NewScanner(file)
	lineCount := 0

	for scanner.Scan() {
		lineCount++
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}

	return lineCount, nil
}

// readFileLines reads specific line ranges from a file (startLine inclusive, endLine exclusive)
func readFileLines(path string, startLine, endLine int) ([]Player, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var players []Player
	currentLine := 0

	// Skip lines before our start
	for currentLine < startLine && scanner.Scan() {
		currentLine++
	}

	// Process our assigned lines
	nums := make([]string, 0, 5)
	picks := make([]int, 0, 5)
	
	for currentLine < endLine && scanner.Scan() {
		line := scanner.Text()
		currentLine++

		nums = nums[:0]
		nums = strings.Fields(line)
		if len(nums) != 5 {
			continue // Skip malformed lines
		}
		
		picks = picks[:0]
		var parseErr error
		for _, s := range nums {
			n, parseErr := strconv.Atoi(s)
			if parseErr != nil {
				break
			}
			picks = append(picks, n)
		}

		if parseErr != nil {
			continue // Skip lines with parse errors
		}
		
		if len(picks) == 5 {
			player := createPlayer(picks)
			players = append(players, player)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return players, nil
}
