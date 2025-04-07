package cache

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/gabdlr/api-cuit-go/utils"
)

const cacheDir = "./.cache"

func cacheIsOld(fPath string) bool {
	fInfo, _ := os.Stat(fPath)
	fCreationDate := fInfo.ModTime()
	oneYearAgo := (time.Now()).AddDate(-1, 0, 0)
	isOld := fCreationDate.Before(oneYearAgo)
	if isOld {
		os.Remove(fPath)
	}
	return isOld
}

func Search(cuit string) ([]byte, error) {
	cuit = utils.StandardizeCuit(cuit)
	fPath := fmt.Sprintf("%s/%s.json", cacheDir, cuit)
	f, err := os.ReadFile(fPath)
	if err == nil {
		if cacheIsOld(fPath) {
			return []byte{0}, errors.New("cached file expired")
		}
		return f, nil
	}
	return []byte{0}, err
}

func Save(cuit string, cuitInfo []byte) {
	cuit = utils.StandardizeCuit(cuit)
	os.WriteFile(fmt.Sprintf("%s/%s.json", cacheDir, cuit), cuitInfo, 0644)
}
