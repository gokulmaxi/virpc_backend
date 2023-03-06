package utilities

import (
	"bufio"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

type StorageInfo struct {
	Totalsize      string
	Used           string
	Available      string
	UsedPercentage string
	MountedOn      string
}

func GetStorage() (StorageInfo, error) {
	cmd := exec.Command("df", "-h", "/")
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("Error:", err)
		return StorageInfo{}, err
	}
	var storageData = StorageInfo{}
	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) >= 2 {
			mount := fields[5]
			if mount == "/" {
				storageData.Totalsize = fields[1]
				storageData.Used = fields[2]
				storageData.Available = fields[3]
				storageData.UsedPercentage = fields[4]
				storageData.MountedOn = fields[5]
				return storageData, nil
			}
		}
	}
	return StorageInfo{}, errors.New("Root fs not  found")
}
