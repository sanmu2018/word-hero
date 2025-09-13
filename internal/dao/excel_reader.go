package dao

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/tealeg/xlsx/v3"
	"github.com/sanmu2018/word-hero/internal/dto"
	"github.com/sanmu2018/word-hero/internal/table"
	"github.com/sanmu2018/word-hero/internal/utils"
	"github.com/sanmu2018/word-hero/log"
)

// ExcelReader handles reading vocabulary from Excel files
type ExcelReader struct {
	filePath string
}

// NewExcelReader creates a new Excel reader instance
func NewExcelReader(filePath string) *ExcelReader {
	return &ExcelReader{
		filePath: filePath,
	}
}

// GetFilePath returns the Excel file path
func (er *ExcelReader) GetFilePath() string {
	return er.filePath
}

// ReadWords reads words from Excel file
func (er *ExcelReader) ReadWords() (*dto.VocabularyList, error) {
	log.Info().Str("file", er.filePath).Msg("Reading Excel file")

	// Open the Excel file
	xlFile, err := xlsx.OpenFile(er.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open Excel file: %v", err)
	}

	var words []table.Word

	// Iterate through all sheets
	for _, sheet := range xlFile.Sheets {
		log.Debug().Str("sheet", sheet.Name).Msg("Processing sheet")

		// Iterate through all rows
		err := sheet.ForEachRow(func(row *xlsx.Row) error {
			rowIndex := row.GetCoordinate()
			if rowIndex == 0 {
				// Skip header row (assuming first row is header)
				return nil
			}

			// Check if row has at least 8 cells (we need column 3 for English and column 8 for Chinese)
			cellCount := 0
			row.ForEachCell(func(cell *xlsx.Cell) error {
				cellCount++
				return nil
			})

			if cellCount >= 8 {
				var english, chinese string
				cellIndex := 0
				row.ForEachCell(func(cell *xlsx.Cell) error {
					if cellIndex == 2 { // Column 3 (0-indexed) contains English word
						english = cell.String()
					} else if cellIndex == 7 { // Column 8 (0-indexed) contains Chinese explanation
						chinese = cell.String()
					}
					cellIndex++
					return nil
				})

				// Trim whitespace and check if both fields have content
				english = trimString(english)
				chinese = trimString(chinese)

				if english != "" && chinese != "" {
					word := table.Word{
						ID:        utils.GenerateUUID(),
						English:   english,
						Chinese:   chinese,
						CreatedAt: time.Now().UnixMilli(),
						UpdatedAt: time.Now().UnixMilli(),
					}
					words = append(words, word)
				}
			}
			return nil
		})

		if err != nil {
			return nil, fmt.Errorf("error reading rows: %v", err)
		}
	}

	if len(words) == 0 {
		return nil, fmt.Errorf("no valid words found in the Excel file")
	}

	log.Info().Int("count", len(words)).Msg("Successfully read vocabulary words")

	return &dto.VocabularyList{Words: words}, nil
}

// trimString trims whitespace from a string
func trimString(s string) string {
	return strings.TrimSpace(s)
}

// ValidateFile checks if the Excel file exists
func (er *ExcelReader) ValidateFile() error {
	// Check if Excel file exists
	if _, err := os.Stat(er.filePath); err != nil {
		return fmt.Errorf("Excel file not found: %s", er.filePath)
	}

	// Try to open the file to verify it's a valid Excel file
	_, err := xlsx.OpenFile(er.filePath)
	if err != nil {
		return fmt.Errorf("failed to open Excel file: %v", err)
	}

	return nil
}

// GetFileInfo returns information about the Excel file
func (er *ExcelReader) GetFileInfo() (string, error) {
	xlFile, err := xlsx.OpenFile(er.filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open Excel file: %v", err)
	}

	info := fmt.Sprintf("Excel file: %s\n", er.filePath)
	info += fmt.Sprintf("Sheets: %d\n", len(xlFile.Sheets))

	for i, sheet := range xlFile.Sheets {
		rowCount := 0
		sheet.ForEachRow(func(row *xlsx.Row) error {
			rowCount++
			return nil
		})
		info += fmt.Sprintf("  Sheet %d: '%s' - %d rows\n", i+1, sheet.Name, rowCount)
	}

	return info, nil
}