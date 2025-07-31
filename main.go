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

		// get line count from file
		lineCount, err = getLineCount(file)
		if err != nil {
			fmt.Println("Error getting line count:", err)
			return
		}

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

func readPlayers(path string, threads byte) ([][]int, error) {
	// Get file size for segment calculation
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	fileInfo, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, err
	}
	fileSize := fileInfo.Size()
	

	// Calculate segment size per thread
	segmentSize := fileSize / int64(threads)
	if segmentSize == 0 {
		segmentSize = fileSize
		threads = 1
	}

	// Channel to collect results from each thread
	resultsChan := make(chan [][]int, threads)
	var wg sync.WaitGroup

	optimalBufferSize := int(segmentSize) + (int(segmentSize)/2)
	optimalBufferSize = nextPowerOf2(int64(optimalBufferSize))

	// Launch reader threads - each reads its own file segment
	for threadID := byte(0); threadID < threads; threadID++ {
		wg.Add(1)
		go func(id byte) {
			defer wg.Done()
			
			// Calculate start and end positions for this thread
			startPos := int64(id) * segmentSize
			endPos := startPos + segmentSize
			if id == threads - 1 {
				endPos = fileSize // Last thread reads until end of file
			}

			players, err := readFileSegment(file, int(id), startPos, endPos, optimalBufferSize)
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
		file.Close()
	}()

	// Combine results from all threads
	var allPlayers [][]int
	totalPlayers := 0
	
	for threadResults := range resultsChan {
		if threadResults != nil {
			allPlayers = append(allPlayers, threadResults...)
			totalPlayers += len(threadResults)
		}
	}

	return allPlayers, nil
}

// readFileSegment reads a specific segment of the file from startPos to endPos
func readFileSegment(file *os.File, threadID int, startPos, endPos int64, optimalBufferSize int) ([][]int, error) {
	// Seek to the start position
	_, err := file.Seek(startPos, 0)
	if err != nil {
		return nil, fmt.Errorf("thread %d: failed to seek to position %d: %v", threadID, startPos, err)
	}

	// Create a limited reader that stops at endPos
	limitedReader := &io.LimitedReader{R: file, N: endPos - startPos}
	scanner := bufio.NewScanner(limitedReader)

	buf := make([]byte, optimalBufferSize)
	scanner.Buffer(buf, optimalBufferSize)

	var players [][]int
	lineCount := 0

	// Thread 0 starts from beginning, others skip first partial line
	if threadID > 0 {
		// Skip the first (potentially partial) line
		if scanner.Scan() {
			// This line is skipped as it might be partial
			lineCount++
		}
	}

	nums := make([]string, 5)
	var picks []int
	var n int
	var line string

	// Read and process complete lines within our segment
	for scanner.Scan() {
		line = scanner.Text()
		lineCount++

		nums = nums[:0]
		nums = strings.Fields(line)
		if len(nums) != 5 {
			nums = nil
		}
		
		picks = nil
		for _, s := range nums {
			n, err = strconv.Atoi(s)
			if err != nil {
				return nil, fmt.Errorf("invalid number: %v", s)
			}
			picks = append(picks, n)
		}

		if err != nil {
			// Log error but continue processing
			fmt.Printf("Thread %d: error parsing line %d: %v\n", threadID, lineCount, err)
			continue
		}
		if picks != nil {
			players = append(players, picks)
		}
	}


	// Check if we need to read one more complete line (if we ended mid-line)
	if limitedReader.N == 0 {
		// We've reached our limit, but might be in the middle of a line
		// Read until the end of the current line
		originalScanner := bufio.NewScanner(file)
		if originalScanner.Scan() {
			line = originalScanner.Text()
			lineCount++

			nums = strings.Fields(line)
			if len(nums) != 5 {
				nums = nil
			}
			
			picks = nil
			for _, s := range nums {
				n, err = strconv.Atoi(s)
				if err != nil {
					return nil, fmt.Errorf("invalid number: %v", s)
				}
				picks = append(picks, n)
			}

			if err != nil {
				fmt.Printf("Thread %d: error parsing final line %d: %v\n", threadID, lineCount, err)
			} else if picks != nil {
				players = append(players, picks)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("thread %d: scanner error: %v", threadID, err)
	}

	return players, nil
}

func nextPowerOf2(n int64) int {
	if n <= 0 {
		return 1
	}
	if n == 1 {
		return 1
	}
	
	// Find the highest bit set
	power := 1
	for power < int(n) {
		power *= 2
		if power > 1048576 { // Cap at 1MB
			return 1048576
		}
	}
	return power
}

// Parallel match counting for maximum performance
func countMatchesParallel(players [][]int, winning []int, threads byte) MatchCount {
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
		go func(playerChunk [][]int) {
			defer wg.Done()
			localResult := MatchCount{2: 0, 3: 0, 4: 0, 5: 0}
			var set map[int]bool
			var match int
			var p, n int

			for _, player := range playerChunk {
				match = countMatches(player, winning, set, match, n, p)
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
