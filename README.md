# Hungarian Lottery Problem (GoLang)

This project solves the Hungarian Lottery winner reporting problem with high efficiency using Go and **segmented parallel file processing**.

## 📌 Problem Description

In the Hungarian lottery, players select **5 distinct numbers** between **1 and 90**. During a weekly draw, the lottery organization also selects **5 distinct numbers**.

The system must quickly report how many players matched 2, 3, 4, or all 5 numbers, based on pre-submitted player entries (up to 10 million).

### 🎯 Sample Output Format

| Number Matching | Winners |
| --------------- | ------- |
| 5               | 0       |
| 4               | 52      |
| 3               | 1511    |
| 2               | 25949   |

## ⚙️ Segmented Parallel File Processing Algorithm

### 🚀 Revolutionary Approach
This implementation uses a **segmented parallel file reading** strategy that achieves true parallel I/O by eliminating the single-threaded file reader bottleneck.

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
- **True Parallel I/O**: Each thread has independent file access
- **Zero Data Loss**: Careful line boundary management ensures no missed entries  
- **Maximum Throughput**: Scales directly with storage parallelism
- **Memory Efficient**: Each thread processes only its segment

## 📈 Performance Benchmarks

### 🖥️ Test Environment
- **CPU**: AMD Ryzen 5 5600G (6 cores, 12 threads)
- **Dataset**: 1,000,000 player entries
- **Threads**: 12 (matching hardware threads)

### ⚡ Performance Results
```
Reading players from file...
Threads: 12
Reading players Execution took 96.3782ms (96378200 ns)

Enter 5 winning numbers (space-separated): 1 3 5 6 8
Counting matches with 12 threads...

Number Matching | Winners
----------------|--------
5               | 0
4               | 52
3               | 1511
2               | 25949

Counting matches Execution took 19.4563499s (19456349900 ns)
```

### 🎯 Key Metrics
- **File Reading**: 96.38ms for 1M entries ≈ **10.4M entries/second**
- **Match Processing**: 19.46s for 1M comparisons ≈ **51.4K comparisons/second**
- **Total Throughput**: **~20.1 seconds** for complete lottery analysis
- **Memory Usage**: Efficient segmented processing with minimal memory footprint

## 🏗️ Project Structure

```
HungarianLotteryProblem/
├── main.go                 # Main application logic with segmented parallel processing
├── test_players.txt        # Generated test data (1M players)
└── README.md               # This file
```

## ✨ Features

- **Segmented Parallel File Processing**: Revolutionary approach eliminating I/O bottlenecks
- **True Parallel I/O**: Each thread has independent file access for maximum throughput
- **Smart Line Boundary Handling**: Zero data loss with careful segment boundary management
- **Parallel Match Counting**: Multi-threaded result computation for optimal performance
- **Input Validation**: Validates player numbers (1-90, exactly 5 numbers per player)
- **Error Handling**: Robust error handling for file operations and user input
- **Memory Efficient**: Segmented processing with optimized memory allocation
- **Hardware Optimization**: Automatic CPU core detection and threading recommendations
- **Cross-platform**: Works on Windows, macOS, and Linux

## 🚀 Run Instructions

### Quick Start
```bash
# Generate test data (1M players)
go run main.go test_players.txt

# Run with optimal threading (match your CPU cores)
go run main.go test_players.txt 12
```

### Command Line Usage
```bash
go run main.go <input_file_path> <number_of_threads>
```

**Parameters:**
- `input_file_path`: Path to the file containing player entries (or new filename to generate test data)
- `number_of_threads`: Number of parallel threads for segmented file processing

**Threading Recommendations:**
- Use your CPU's thread count for optimal performance (e.g., 12 threads for 6-core/12-thread CPU)
- The system will provide automatic recommendations based on detected CPU cores

### Performance Optimization Tips
```bash
# For AMD Ryzen 5 5600G (6 cores, 12 threads) - optimal setting:
go run main.go test_players.txt 12

# For Intel i7-8700K (6 cores, 12 threads):
go run main.go test_players.txt 12

# For AMD Ryzen 9 5900X (12 cores, 24 threads):
go run main.go test_players.txt 24
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
If the specified file doesn't exist, the program automatically generates 1,000,000 test players:
```bash
# This will create test_players.txt with 1M random entries
go run main.go test_players.txt 12
```

## 🏆 Algorithm Innovation

This implementation represents a significant advancement in parallel file processing:

- **Eliminates Single-Point Bottlenecks**: Traditional producer-consumer patterns create I/O bottlenecks
- **True Parallel I/O**: Each thread directly accesses the file system independently  
- **Optimal Resource Utilization**: Scales with both CPU cores and storage parallelism
- **Zero Data Loss**: Sophisticated line boundary handling ensures perfect data integrity

The segmented approach achieves **10.4M entries/second** reading performance, demonstrating the effectiveness of eliminating sequential file access patterns in favor of true parallel I/O operations.
