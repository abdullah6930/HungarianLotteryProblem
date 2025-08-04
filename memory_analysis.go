package main

import "fmt"

func main() {
	const playersCount = 10_000_000
	
	fmt.Println("=== Memory Usage Analysis: Direct Byte Storage (10M Players) ===")
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
	
	fmt.Printf("OPTIMIZED ([5]byte with direct byte storage):\n")
	fmt.Printf("  Per player: %d bytes (%d bits)\n", optimizedBytesPerPlayer, optimizedBytesPerPlayer*8)
	fmt.Printf("  Total: %d bytes (%.2f MB) = %d bits (%.2f Gb)\n", 
		optimizedTotalBytes, float64(optimizedTotalBytes)/(1024*1024), optimizedTotalBits, float64(optimizedTotalBits)/(1024*1024*1024))
	fmt.Println()
	
	fmt.Printf("THEORETICAL MINIMUM (if we could use fractional bytes):\n")
	fmt.Printf("  Per player: %.2f bytes (%.1f bits needed for 1-90 range)\n", float64(theoreticalBitsPerPlayer)/8, 6.49) // log2(90) ≈ 6.49
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
	
	fmt.Printf("DIRECT BYTE STORAGE EFFICIENCY:\n")
	fmt.Printf("  Numbers 1-90 fit perfectly in 1 byte (max 255)\n")
	fmt.Printf("  No wasted bits or complex encoding needed\n")
	fmt.Printf("  Simple, fast, and memory efficient\n")
	fmt.Println()
	
	fmt.Printf("BYTE STORAGE BREAKDOWN PER PLAYER:\n")
	fmt.Printf("  Number 1: [1 byte] = 8 bits\n")
	fmt.Printf("  Number 2: [1 byte] = 8 bits\n")
	fmt.Printf("  Number 3: [1 byte] = 8 bits\n")
	fmt.Printf("  Number 4: [1 byte] = 8 bits\n")
	fmt.Printf("  Number 5: [1 byte] = 8 bits\n")
	fmt.Printf("  Total per player: = 40 bits (5 bytes)\n")
	fmt.Printf("  Direct storage: Numbers 1-90 → byte values 1-90\n")
	fmt.Println()
	
	fmt.Printf("SCALE EXAMPLES:\n")
	scales := []int{1_000, 10_000, 100_000, 1_000_000, 10_000_000, 100_000_000}
	for _, scale := range scales {
		oldMB := float64(scale * originalBytesPerPlayer) / (1024 * 1024)
		newMB := float64(scale * optimizedBytesPerPlayer) / (1024 * 1024)
		fmt.Printf("  %10s players: %.3f MB → %.3f MB (saved %.3f MB)\n", 
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