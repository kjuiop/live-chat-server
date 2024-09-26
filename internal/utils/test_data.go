package utils

import (
	"embed"
	"log"
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
