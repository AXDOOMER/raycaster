// File: record.go (Raycaster)
// Copyright (C) 2021 Alexandre-Xavier Labonté-Lamoureux
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package game

import (
	"compress/gzip"
	"encoding/gob"
	"fmt"
	"os"
)

var commands []Keyboard
var demoFileName string
var isRecording bool = false

func writeDemoCommand(kb Keyboard) {
	if isRecording {
		commands = append(commands, kb)
	}
}

func writeDemo(filePath string) {
	file, err := os.Create(filePath)
	if err == nil {
		writer := gzip.NewWriter(file)
		encoder := gob.NewEncoder(writer)
		encoder.Encode(commands)
		writer.Close()
	} else {
		panic(err)
	}

	file.Close()
}

func writeDemoFile() {
	if isRecording {
		writeDemo(demoFileName)
	}
}

func readDemoFile(filePath string) {
	file, err := os.Open(filePath)
	if err == nil {
		reader, err := gzip.NewReader(file)
		if err == nil {
			decoder := gob.NewDecoder(reader)
			err = decoder.Decode(&commands)
			if err != nil {
				panic(err)
			}
			fmt.Println("Will play demo file", filePath)
		} else {
			panic(err)
		}
		reader.Close()
	} else {
		// When the file doesn't exist, we're in record mode
		fmt.Println("A demo will be recorded to file", filePath)
		isRecording = true
	}

	file.Close()
}

func readDemoCommand(index int, kb *Keyboard, running bool) bool {
	if !isRecording && len(demoFileName) > 0 {
		if index < len(commands) {
			*kb = commands[index]
			return running
		} else {
			// finished
			return false
		}
	}
	return running
}
