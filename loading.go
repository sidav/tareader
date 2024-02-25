package main

import (
	"os"
	tafilesread "totala_reader/ta_files_read"
	"totala_reader/ta_files_read/texture"
)

func readAllGAFsFromDirectory(directoryName string) []*texture.GafEntry {
	pp("Reading all GAF entries from dir %s", directoryName)
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
			r := &tafilesread.Reader{}
			r.ReadFromFile(openedFileName)
			readedGAFEntries := texture.ReadTextureFromReader(r, false)
			for _, e := range readedGAFEntries {
				allEntries = append(allEntries, e)
			}
		}
	}
	return allEntries
}
