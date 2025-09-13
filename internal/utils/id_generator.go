package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/google/uuid"
)

// SnowflakeConfig holds configuration for snowflake ID generation
type SnowflakeConfig struct {
	WorkerID     int64
	Epoch        int64 // Custom epoch (default: 2024-01-01)
	WorkerIDBits int64 // Number of bits for worker ID (default: 10)
	SequenceBits int64 // Number of bits for sequence (default: 12)
}

// SnowflakeGenerator generates snowflake IDs
type SnowflakeGenerator struct {
	config      SnowflakeConfig
	sequence    int64
	lastTime    int64
	mutex       sync.Mutex
	maxSequence int64
	workerIDMax int64
	timeShift   int64
	workerShift int64
}

// NewSnowflakeGenerator creates a new snowflake ID generator
func NewSnowflakeGenerator(config SnowflakeConfig) *SnowflakeGenerator {
	if config.WorkerIDBits <= 0 {
		config.WorkerIDBits = 10
	}
	if config.SequenceBits <= 0 {
		config.SequenceBits = 12
	}
	if config.Epoch == 0 {
		config.Epoch = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC).UnixMilli()
	}

	gen := &SnowflakeGenerator{
		config:      config,
		maxSequence: 1<<config.SequenceBits - 1,
		workerIDMax: 1<<config.WorkerIDBits - 1,
	}

	// Calculate bit shifts
	gen.timeShift = config.SequenceBits + config.WorkerIDBits
	gen.workerShift = config.SequenceBits

	// Validate worker ID
	if config.WorkerID < 0 || config.WorkerID > gen.workerIDMax {
		panic(fmt.Sprintf("worker ID must be between 0 and %d", gen.workerIDMax))
	}

	return gen
}

// Generate generates a new snowflake ID
func (g *SnowflakeGenerator) Generate() int64 {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	now := time.Now().UnixMilli()

	if now == g.lastTime {
		g.sequence++
		if g.sequence > g.maxSequence {
			// Sequence overflow, wait for next millisecond
			for now <= g.lastTime {
				now = time.Now().UnixMilli()
			}
			g.sequence = 0
		}
	} else {
		g.sequence = 0
	}

	g.lastTime = now

	// Calculate snowflake ID
	// Format: timestamp | worker_id | sequence
	id := ((now - g.config.Epoch) << g.timeShift) |
		(g.config.WorkerID << g.workerShift) |
		g.sequence

	return id
}

// ParseSnowflakeID parses a snowflake ID back into its components
func (g *SnowflakeGenerator) ParseSnowflakeID(id int64) (timestamp, workerID, sequence int64) {
	sequence = id & ((1 << g.config.SequenceBits) - 1)
	workerID = (id >> g.config.SequenceBits) & ((1 << g.config.WorkerIDBits) - 1)
	timestamp = (id >> (g.config.SequenceBits + g.config.WorkerIDBits)) + g.config.Epoch
	return timestamp, workerID, sequence
}

// GetTimeFromID extracts the timestamp from a snowflake ID
func (g *SnowflakeGenerator) GetTimeFromID(id int64) time.Time {
	timestamp, _, _ := g.ParseSnowflakeID(id)
	return time.UnixMilli(timestamp)
}

// Global snowflake generator instance
var defaultSnowflake *SnowflakeGenerator

func init() {
	defaultSnowflake = NewSnowflakeGenerator(SnowflakeConfig{
		WorkerID: 1, // Default worker ID, should be configured per instance
	})
}

// GenerateUUID generates a new UUID
func GenerateUUID() string {
	return uuid.New().String()
}

// GenerateSnowflakeID generates a new snowflake ID using default generator
func GenerateSnowflakeID() int64 {
	return defaultSnowflake.Generate()
}

// SnowflakeIDToString converts a snowflake ID to string
func SnowflakeIDToString(id int64) string {
	return fmt.Sprintf("%d", id)
}

// StringToSnowflakeID converts a string to snowflake ID
func StringToSnowflakeID(s string) (int64, error) {
	var id int64
	_, err := fmt.Sscanf(s, "%d", &id)
	return id, err
}

// GenerateRandomString generates a random string of specified length
func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			// Fallback to simpler random generation
			b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
		} else {
			b[i] = charset[n.Int64()]
		}
	}
	return string(b)
}

// IDGenerator interface for different ID generation strategies
type IDGenerator interface {
	GenerateID() string
}

// UUIDGenerator implements IDGenerator using UUID
type UUIDGenerator struct{}

func (g *UUIDGenerator) GenerateID() string {
	return GenerateUUID()
}

// SnowflakeStringGenerator implements IDGenerator using snowflake as string
type SnowflakeStringGenerator struct{}

func (g *SnowflakeStringGenerator) GenerateID() string {
	return SnowflakeIDToString(GenerateSnowflakeID())
}

// GetDefaultIDGenerator returns the default ID generator (UUID)
func GetDefaultIDGenerator() IDGenerator {
	return &UUIDGenerator{}
}

// GetSnowflakeIDGenerator returns a snowflake-based ID generator
func GetSnowflakeIDGenerator() IDGenerator {
	return &SnowflakeStringGenerator{}
}