# \# Hungarian Lottery Problem (GoLang)

# 

# This project solves the Hungarian Lottery winner reporting problem with high efficiency using Go and concurrent processing.

# 

# \## 📌 Problem Description

# 

# In the Hungarian lottery, players select \*\*5 distinct numbers\*\* between \*\*1 and 90\*\*. During a weekly draw, the lottery organization also selects \*\*5 distinct numbers\*\*.

# 

# The system must quickly report how many players matched 2, 3, 4, or all 5 numbers, based on pre-submitted player entries (up to 10 million).

# 

\### 🎯 Sample Output Format

| Number Matching | Winners |
===

# | --------------- | ------- |

# | 5               | 0       |

# | 4               | 12      |

# | 3               | 818     |

| 2               | 22613   |

===

# \## ⚙️ How It Works

# 

# \- Players' entries are preloaded from a file.
# \- Upon receiving the 5 winning numbers, the system calculates match counts for each player.
# \- Results are aggregated and printed in a categorized summary.

# \## 🏗️ Project Structure

# ```
# HungarianLotteryProblem/
# ├── main.go                 # Main application logic
# ├── main_test.go            # Unit tests
# ├── generate_data.go        # Data generator (standard)
# ├── generate_data_fast.go   # Data generator (optimized)
# ├── generate_data.bat       # Windows batch script
# ├── generate_data.ps1       # PowerShell script
# ├── go.mod                  # Go module definition
# ├── go.sum                  # Dependency checksums
# ├── Makefile                # Build and run commands
# ├── .gitignore              # Git ignore rules
# ├── sample_players.txt      # Sample player data
# └── README.md               # This file
# ```

# \## ✨ Features

# \- **Concurrent Processing**: Uses Go goroutines for high-performance parallel processing
# \- **Input Validation**: Validates player numbers (1-90, no duplicates)
# \- **Error Handling**: Robust error handling for file operations and user input
# \- **Memory Efficient**: Processes players in batches to handle large datasets
# \- **Test Coverage**: Comprehensive unit tests for core functionality
# \- **Cross-platform**: Works on Windows, macOS, and Linux

# 

# \## 🚀 Run Instructions

# 

# ### Quick Start
# ```bash
# # Run tests
# go test -v
# 
# # Run with sample data
# go run main.go sample_players.txt 4
# 
# # Or use Makefile commands
# make test
# make run-sample
# make run-sample-auto
# ```

# ### 📊 Generate Test Data
# ```bash
# # Generate 10,000 players (test)
# go run generate_data.go -n 10000 -o test_10k.txt
# 
# # Generate 1,000,000 players
# go run generate_data.go -n 1000000 -o players_1m.txt
# 
# # Generate 10,000,000 players (fast version)
# go run generate_data_fast.go -n 10000000 -o players_10m.txt -workers 8
# 
# # Or use Windows scripts
# .\generate_data.bat
# .\generate_data.ps1
# ```

# 

# ### Command Line Usage
# ```bash
# go run main.go <input_file_path> <number_of_threads>
# ```

# **Parameters:**
# - `input_file_path`: Path to the file containing player entries
# - `number_of_threads`: Number of goroutines to use for concurrent processing

# 

# ### Build and Run
# ```bash
# # Build the executable
# make build
# 
# # Run the executable
# ./lottery sample_players.txt 4
# ```

# 

# 📄 Player File Format

# Each line represents one player, with exactly 5 distinct numbers between 1 and 90, separated by spaces:

# ```
# 4 79 13 80 56
# 71 84 48 85 38
# 41 65 39 82 36
# ...
# ```





