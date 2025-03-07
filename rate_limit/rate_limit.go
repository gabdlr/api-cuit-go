package rate_limit

import (
	"encoding/gob"
	"os"
	"strings"
	"time"
)

const TIMEFRAME = 60

func TimeLeft(addr string) int64 {
	addr = (strings.Split(addr, ":"))[0]
	timeLeft := int64(0)
	file, err := os.OpenFile("addr_table.gob", os.O_RDWR|os.O_CREATE, 0644)

	if err == nil {
		defer file.Close()
		loadedData := make(map[string]int64)
		gob.NewDecoder(file).Decode(&loadedData)

		if loadedData[addr] == 0 {
			loadedData[addr] = time.Now().Unix()
			saveDate(file, loadedData)

		} else {
			t := time.Now().Unix() - loadedData[addr]
			if t >= TIMEFRAME {
				loadedData[addr] = time.Now().Unix()
				saveDate(file, loadedData)
			} else {
				timeLeft = TIMEFRAME - t
			}
		}

	}
	return timeLeft
}

func saveDate(file *os.File, data map[string]int64) {
	file.Seek(0, 0)
	gob.NewEncoder(file).Encode(data)
}
