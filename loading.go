package main

import (
	"fmt"
	"os"
	tafilesread "totala_reader/ta_files_read"
	"totala_reader/ta_files_read/texture"
)

func readAllGAFsFromDirectory(directoryName string) []*texture.GafEntry {
	var allEntries []*texture.GafEntry
	if directoryName[len(directoryName)-1] != "/"[0] {
		directoryName += "/"
	}
	items, _ := os.ReadDir(directoryName)
	for _, item := range items {
		if item.IsDir() {
			// do nothing
		} else {
			openedFileName := directoryName + item.Name()
			fmt.Println("Opening " + openedFileName)
			r := &tafilesread.Reader{}
			r.ReadFromFile(openedFileName)
			readedGAFEntries := texture.ReadTextureFromReader(r)
			for _, e := range readedGAFEntries {
				allEntries = append(allEntries, e)
			}
		}
	}
	return allEntries
}
