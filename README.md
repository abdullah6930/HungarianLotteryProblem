# Hungarian Lottery Problem (GoLang)

This project solves the Hungarian Lottery winner reporting problem with **revolutionary memory efficiency** using Go, **segmented parallel file processing**, and **optimized byte storage**.

## 📌 Problem Description

In the Hungarian lottery, players select **5 distinct numbers** between **1 and 90**. Players submit their entries in advance, which must be loaded and processed by the system **before** the weekly draw occurs.

During the weekly draw, the lottery organization selects **5 distinct numbers**. The system must then quickly report how many players matched 2, 3, 4, or all 5 numbers from the pre-loaded player entries (up to 10 million).

### 🔄 Processing Sequence
1. **Pre-processing**: Load all player entries into memory with optimized storage
2. **Draw Event**: Lottery organization announces the 5 winning numbers  
3. **Analysis**: Rapidly compare all players against winning numbers
4. **Report**: Generate winner statistics in milliseconds

### 🎯 Sample Output Format

| Number Matching | Winners |
| --------------- | ------- |
| 5               | 0       |
| 4               | 52      |
| 3               | 1511    |
| 2               | 25949   |

## ⚙️ Revolutionary Optimizations

### 🚀 1. Optimized Byte Storage
This implementation uses **direct byte storage** for ultra-efficient memory usage:

- **Numbers 1-90**: Each stored directly in 1 byte (numbers ≤ 255)
- **No overhead**: Simple, direct storage without bit manipulation
- **Per player**: 5 bytes total
- **10M players**: 47.68 MB total memory usage

#### Storage Example
```
Original numbers: 1 4 22 56 89
Stored as bytes:  [1] [4] [22] [56] [89]
Memory layout:    5 bytes per player, direct access
```

### 🚀 2. Segmented Parallel File Processing Algorithm
Additionally, a **segmented parallel file reading** strategy achieves true parallel I/O by eliminating the single-threaded file reader bottleneck.

### 📊 Algorithm Steps
1. **File Segmentation**: Divide file into `n` equal segments based on byte positions
2. **Thread Assignment**: Each thread gets a segment: `segmentSize = fileSize / threads`  
3. **Parallel File Access**: Each thread opens its own file handle for true parallel I/O
4. **Smart Line Boundary Handling**:
   - **Thread 0**: Reads from file start, processes all complete lines
   - **Thread n**: Reads from `(segmentSize * n) + 1`, skips first partial line
   - **All Threads**: Complete final partial lines beyond segment boundaries
5. **Result Aggregation**: Main thread combines results from all worker threads

### 🎯 Key Advantages
- **Ultra-Compact Memory**: Direct byte storage achieves 5 bytes per player
- **True Parallel I/O**: Each thread has independent file access
- **Zero Data Loss**: Careful line boundary management ensures no missed entries  
- **Maximum Throughput**: Scales directly with storage parallelism
- **Simple & Fast**: No complex encoding/decoding overhead
- **Memory Efficient**: Each thread processes only its segment with optimized data

## 📈 Performance Benchmarks

### 🖥️ Test Environment
- **CPU**: AMD Ryzen 5 5600G (6 cores, 12 threads)
- **Dataset**: 10,000,000 player entries
- **Threads**: 12 (matching hardware threads)

### ⚡ Performance Results
```
CPU cores: 12
Reading players from file...
Threads: 12
Reading players Execution took 186.0733ms (186073300 ns)

Enter 5 winning numbers (space-separated): 1 2 3 4 5
Counting matches with 12 threads...

Number Matching | Winners
----------------|--------
5               | 4
4               | 306
3               | 9132
2               | 157194

Counting matches Execution took 89.5442ms (89544200 ns)
```

### 🎯 Key Metrics
- **File Reading**: 186ms for 10M entries ≈ **53.8M entries/second**
- **Memory Usage**: **47.68 MB** for 10M players
- **Match Processing**: 89.5ms for 10M comparisons ≈ **111.7M comparisons/second**
- **Total Throughput**: **~275ms** for complete 10M player lottery analysis
- **Memory Efficiency**: Ultra-compact 5 bytes per player with direct byte storage

## 🏗️ Project Structure

```
HungarianLotteryProblem/
├── main.go                 # Main application with memory optimization & parallel processing
├── memory_analysis.go      # Comprehensive memory usage analysis tool
├── test_players.txt        # Generated test data (10M players)
└── README.md               # This file
```

## Run Instructions

### Quick Start
```bash
# Step 1: Generate test data (10M players) - run this first
go run main.go test_players.txt

# Step 2: Run lottery analysis with optimal threading (match your CPU cores)
go run main.go test_players.txt 12

# Step 3: Analyze memory optimization results
go run memory_analysis.go
```

### Command Line Usage

**Generation Mode** (creates test data):
```bash
go run main.go <filename>
```

**Analysis Mode** (runs lottery analysis):
```bash
go run main.go <input_file_path> <number_of_threads>
```

**Parameters:**
- `input_file_path`: Path to the file containing player entries
- `number_of_threads`: Number of parallel threads for segmented file processing (omit to generate test data)

**Threading Recommendations:**
- Use your CPU's thread count for optimal performance (e.g., 12 threads for 6-core/12-thread CPU)
- The system will provide automatic recommendations based on detected CPU cores

### Memory Analysis Tool
```bash
# Run comprehensive memory analysis
go run memory_analysis.go
```

**Features:**
- **Memory usage analysis**: Detailed memory calculations
- **Byte-level breakdown**: Detailed storage analysis  
- **Scale examples**: Memory usage from 1K to 100M players
- **Efficiency metrics**: Performance characteristics

### Performance Optimization Tips
```bash
# For AMD Ryzen 5 5600G (6 cores, 12 threads):
go run main.go test_players.txt          # Generate data first
go run main.go test_players.txt 12       # Run analysis

# For Intel i7-8700K (6 cores, 12 threads):
go run main.go test_players.txt          # Generate data first  
go run main.go test_players.txt 12       # Run analysis

# For AMD Ryzen 9 5900X (12 cores, 24 threads):
go run main.go test_players.txt          # Generate data first
go run main.go test_players.txt 24       # Run analysis
```

### Memory Usage at Scale
```bash
# Memory usage with direct byte storage:
     1,000 players:   0.005 MB
    10,000 players:   0.048 MB
   100,000 players:   0.477 MB
 1,000,000 players:   4.768 MB
10,000,000 players:  47.68 MB
```

## 📄 Player File Format

Each line represents one player, with exactly 5 distinct numbers between 1 and 90, separated by spaces:

```
4 79 13 80 56
71 84 48 85 38
41 65 39 82 36
...
```

### Auto-Generation Feature
To generate test data, run the program with only the filename (no thread count):
```bash
# Step 1: Generate test_players.txt with 10M random entries
go run main.go test_players.txt

# Step 2: Run analysis after file is generated
go run main.go test_players.txt 12
```

**Note**: The program detects when no thread count is provided and automatically generates 10,000,000 test players.

## 🧬 Technical Implementation Details

### Direct Byte Storage Algorithm
- **Simple storage**: Store values 1-90 directly in bytes
- **No overhead**: No bit manipulation or complex encoding
- **Byte alignment**: Each number fits in exactly 1 byte
- **Memory efficiency**: Ultra-compact 5 bytes per player

### Player Structure
```go
type Player [5]byte

// Example storage for [1, 4, 22, 56, 89]:
// [1] [4] [22] [56] [89]
// Direct byte storage - no encoding needed
```

### Big O Complexity Analysis

#### Time Complexity
- **File Reading**: O(n) - Linear with number of players
- **Player Storage**: O(n) - O(1) per player × n players
- **Match Processing**: O(n) - Each player compared against 5 winning numbers (constant)
- **Result Aggregation**: O(t) - Where t is number of threads (typically << n)
- **Overall System**: O(n) - Linear scalability with player count

#### Space Complexity
- **Raw Player Data**: O(n) - 5 bytes per player with direct storage
- **Working Memory**: O(n/t) per thread - Segmented processing
- **Result Storage**: O(1) - Fixed-size match counters
- **Total Memory**: O(n) - Optimal linear space usage

#### Parallel Processing Complexity
- **Thread Coordination**: O(t) - Where t = number of threads
- **File Segmentation**: O(1) - Constant time segment calculation  
- **Load Balancing**: O(n/t) - Even distribution across threads
- **Result Merging**: O(t) - Combine results from all threads

### Performance Characteristics
- **Storage**: O(1) per player - direct byte assignment
- **Access**: O(1) per player - direct byte-to-int conversion
- **Cache Efficiency**: 8x better locality due to compact representation
- **Memory Bandwidth**: Reduced by 87.5% compared to standard integer arrays
- **GC Pressure**: Minimal due to compact allocations
- **Simplicity**: No encoding/decoding overhead
