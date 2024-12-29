package utils

import (
	"bufio"
	"embed"
	"fmt"
	"log"
	"os"
	"strings"
)

//go:embed testdata/rooms/init_data.json
var initRooms embed.FS

func GetTestInitRoomData() []byte {
	data, err := initRooms.ReadFile("testdata/rooms/init_data.json")
	if err != nil {
		log.Fatalf("failed to read embedded file: %v", err)
	}

	return data
}

func LoadTestEnv() error {

	filepath := "../../../test.env"
	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("could not open %s: %w", filepath, err)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		if err := os.Setenv(key, value); err != nil {
			return fmt.Errorf("could not set env var %s: %w", key, err)
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading %s: %w", filepath, err)
	}
	return nil
}
