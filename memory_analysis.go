package main

import "fmt"

func main() {
	const playersCount = 10_000_000
	
	fmt.Println("=== Memory Usage Analysis for 10 Million Players ===")
	fmt.Println()
	
	// Original implementation: [][]int
	originalBytesPerPlayer := 5 * 8 // 5 ints * 8 bytes per int
	originalTotalBytes := playersCount * originalBytesPerPlayer
	originalTotalBits := originalTotalBytes * 8
	
	// Optimized implementation: [5]byte
	optimizedBytesPerPlayer := 5 // 5 bytes per player
	optimizedTotalBytes := playersCount * optimizedBytesPerPlayer
	optimizedTotalBits := optimizedTotalBytes * 8
	
	// Theoretical minimum (if we could use exactly 35 bits per player)
	theoreticalBitsPerPlayer := 35 // 5 numbers * 7 bits each (without space bits)
	theoreticalTotalBits := playersCount * theoreticalBitsPerPlayer
	theoreticalTotalBytes := (theoreticalTotalBits + 7) / 8 // Round up to nearest byte
	
	fmt.Printf("ORIGINAL ([][]int):\n")
	fmt.Printf("  Per player: %d bytes (%d bits)\n", originalBytesPerPlayer, originalBytesPerPlayer*8)
	fmt.Printf("  Total: %d bytes (%.2f MB) = %d bits (%.2f Gb)\n", 
		originalTotalBytes, float64(originalTotalBytes)/(1024*1024), originalTotalBits, float64(originalTotalBits)/(1024*1024*1024))
	fmt.Println()
	
	fmt.Printf("OPTIMIZED ([5]byte with custom bit mapping):\n")
	fmt.Printf("  Per player: %d bytes (%d bits)\n", optimizedBytesPerPlayer, optimizedBytesPerPlayer*8)
	fmt.Printf("  Total: %d bytes (%.2f MB) = %d bits (%.2f Gb)\n", 
		optimizedTotalBytes, float64(optimizedTotalBytes)/(1024*1024), optimizedTotalBits, float64(optimizedTotalBits)/(1024*1024*1024))
	fmt.Println()
	
	fmt.Printf("THEORETICAL MINIMUM (35 bits per player, no byte alignment):\n")
	fmt.Printf("  Per player: %.2f bytes (%d bits)\n", float64(theoreticalBitsPerPlayer)/8, theoreticalBitsPerPlayer)
	fmt.Printf("  Total: %d bytes (%.2f MB) = %d bits (%.2f Gb)\n", 
		theoreticalTotalBytes, float64(theoreticalTotalBytes)/(1024*1024), theoreticalTotalBits, float64(theoreticalTotalBits)/(1024*1024*1024))
	fmt.Println()
	
	memorySaved := originalTotalBytes - optimizedTotalBytes
	compressionRatio := float64(optimizedTotalBytes) / float64(originalTotalBytes)
	theoreticalSavings := originalTotalBytes - theoreticalTotalBytes
	
	fmt.Printf("SAVINGS:\n")
	fmt.Printf("  Memory saved (optimized vs original): %d bytes (%.2f MB)\n", memorySaved, float64(memorySaved)/(1024*1024))
	fmt.Printf("  Compression ratio: %.1f%% (%.1fx smaller)\n", compressionRatio*100, 1/compressionRatio)
	fmt.Printf("  Theoretical maximum savings: %d bytes (%.2f MB)\n", theoreticalSavings, float64(theoreticalSavings)/(1024*1024))
	fmt.Printf("  Our efficiency vs theoretical: %.1f%%\n", float64(theoreticalTotalBytes)/float64(optimizedTotalBytes)*100)
	fmt.Println()
	
	fmt.Printf("USER'S CALCULATION vs REALITY:\n")
	fmt.Printf("  User calculated: 800,000,000 bits\n")
	fmt.Printf("  Actual usage:    %d bits\n", optimizedTotalBits)
	fmt.Printf("  Difference:      %s by %.1fx\n", 
		map[bool]string{true: "OVER-estimated", false: "UNDER-estimated"}[800_000_000 > optimizedTotalBits],
		float64(800_000_000)/float64(optimizedTotalBits))
	fmt.Println()
	
	fmt.Printf("BIT BREAKDOWN PER PLAYER:\n")
	fmt.Printf("  Number 1: [SPACE_BIT][7-bit number] = 8 bits\n")
	fmt.Printf("  Number 2: [SPACE_BIT][7-bit number] = 8 bits\n")
	fmt.Printf("  Number 3: [SPACE_BIT][7-bit number] = 8 bits\n")
	fmt.Printf("  Number 4: [SPACE_BIT][7-bit number] = 8 bits\n")
	fmt.Printf("  Number 5: [0][7-bit number]         = 8 bits\n")
	fmt.Printf("  Total per player:                   = 40 bits (5 bytes)\n")
	fmt.Println()
	
	fmt.Printf("SCALE EXAMPLES:\n")
	scales := []int{1_000, 10_000, 100_000, 1_000_000, 10_000_000, 100_000_000}
	for _, scale := range scales {
		oldMB := float64(scale * originalBytesPerPlayer) / (1024 * 1024)
		newMB := float64(scale * optimizedBytesPerPlayer) / (1024 * 1024)
		fmt.Printf("  %10s players: %.2f MB → %.2f MB (saved %.2f MB)\n", 
			formatNumber(scale), oldMB, newMB, oldMB-newMB)
	}
}

func formatNumber(n int) string {
	if n >= 1_000_000 {
		return fmt.Sprintf("%.1fM", float64(n)/1_000_000)
	} else if n >= 1_000 {
		return fmt.Sprintf("%.1fK", float64(n)/1_000)
	}
	return fmt.Sprintf("%d", n)
}