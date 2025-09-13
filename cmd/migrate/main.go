package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/sanmu2018/word-hero/internal/conf"
	"github.com/sanmu2018/word-hero/internal/dao"
)

func main() {
	// Parse command line flags
	excelFile := flag.String("excel", "", "Path to Excel file (default: configs/words/IELTS.xlsx)")
	force := flag.Bool("force", false, "Force import even if data already exists")
	clean := flag.Bool("clean", false, "Clean existing data before import")
	verbose := flag.Bool("verbose", false, "Verbose logging")
	help := flag.Bool("help", false, "Show help message")

	flag.Parse()

	if *help {
		showHelp()
		return
	}

	// Set log level based on verbose flag
	if *verbose {
		// Note: In a real implementation, you'd configure the logger level
		fmt.Println("Verbose mode enabled")
	}

	fmt.Println("=== Word Hero Data Migration Tool ===")
	fmt.Println()

	// Load configuration
	config, err := conf.LoadConfig()
	if err != nil {
		fmt.Printf("❌ Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize database
	fmt.Println("🔌 Connecting to database...")
	if err := dao.InitDatabase(&config.Database); err != nil {
		fmt.Printf("❌ Failed to initialize database: %v\n", err)
		os.Exit(1)
	}

	// Run database migrations
	fmt.Println("🔄 Running database migrations...")
	if err := dao.RunMigrations(); err != nil {
		fmt.Printf("❌ Failed to run database migrations: %v\n", err)
		os.Exit(1)
	}

	// Determine Excel file path
	if *excelFile == "" {
		*excelFile = config.App.ExcelFile
	}

	// Check if Excel file exists
	if _, err := os.Stat(*excelFile); err != nil {
		fmt.Printf("❌ Excel file not found: %s\n", *excelFile)
		os.Exit(1)
	}

	// Initialize WordDAO
	wordDAO := dao.NewWordDAO()

	// Check if data already exists
	isEmpty, err := wordDAO.IsEmpty()
	if err != nil {
		fmt.Printf("❌ Failed to check database status: %v\n", err)
		os.Exit(1)
	}

	if !isEmpty && !*force {
		fmt.Printf("⚠️  Words table already contains data. Use --force to overwrite or --clean to clean first.\n")
		fmt.Printf("   Current word count: ")
		if count, err := wordDAO.GetWordCount(); err == nil {
			fmt.Printf("%d\n", count)
		} else {
			fmt.Printf("Unknown\n")
		}
		os.Exit(1)
	}

	// Clean existing data if requested
	if *clean && !isEmpty {
		fmt.Println("🧹 Cleaning existing data...")
		if err := wordDAO.DeleteAllWords(); err != nil {
			fmt.Printf("❌ Failed to clean existing data: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✅ Existing data cleaned")
	}

	// Start migration
	startTime := time.Now()
	fmt.Printf("📖 Reading Excel file: %s\n", *excelFile)

	// Read Excel file
	excelReader := dao.NewExcelReader(*excelFile)
	if err := excelReader.ValidateFile(); err != nil {
		fmt.Printf("❌ Invalid Excel file: %v\n", err)
		os.Exit(1)
	}

	wordList, err := excelReader.ReadWords()
	if err != nil {
		fmt.Printf("❌ Failed to read Excel file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("📊 Found %d words in Excel file\n", len(wordList.Words))

	// Import to database
	fmt.Println("💾 Importing words to database...")
	if err := wordDAO.BulkImport(wordList.Words); err != nil {
		fmt.Printf("❌ Failed to import words: %v\n", err)
		os.Exit(1)
	}

	// Verify import
	importedCount, err := wordDAO.GetWordCount()
	if err != nil {
		fmt.Printf("❌ Failed to verify import: %v\n", err)
		os.Exit(1)
	}

	duration := time.Since(startTime)
	fmt.Println()
	fmt.Println("🎉 Migration completed successfully!")
	fmt.Printf("⏱️  Duration: %v\n", duration)
	fmt.Printf("📝 Words imported: %d\n", importedCount)
	fmt.Printf("📄 Source file: %s\n", *excelFile)
	fmt.Printf("🗄️  Database: %s\n", config.Database.DBName)
}

func showHelp() {
	fmt.Printf("Word Hero Data Migration Tool\n\n")
	fmt.Printf("Usage: %s [options]\n\n", os.Args[0])
	fmt.Printf("Options:\n")
	fmt.Printf("  --excel string   Path to Excel file (default: configs/words/IELTS.xlsx)\n")
	fmt.Printf("  --force          Force import even if data already exists\n")
	fmt.Printf("  --clean          Clean existing data before import\n")
	fmt.Printf("  --verbose        Enable verbose logging\n")
	fmt.Printf("  --help           Show this help message\n\n")
	fmt.Printf("Examples:\n")
	fmt.Printf("  %s                                    # Use default Excel file\n", os.Args[0])
	fmt.Printf("  %s --excel /path/to/words.xlsx       # Custom Excel file\n", os.Args[0])
	fmt.Printf("  %s --force --clean                    # Clean and force import\n", os.Args[0])
	fmt.Printf("  %s --help                             # Show help\n", os.Args[0])
}