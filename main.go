// Copyright (c) 2014 Andrea Masi. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE.txt file.

// to-laser is a gcode post processor for jscut that prepares input
// for a Grbl controlled laser cutter.
//
// Usage
//
// Just redirect gcode via pipe:
//
// 	cat gcode.nc | tl > laser.gcode
package main

// TODO implements these as parameters
//
// laser on command
// laser off command
// laser power
// power on delay

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"

	"github.com/eraclitux/cfgp"
)

const (
	delayGcode     = "G4 P0.5"
	laserOnGcode   = "M3"
	laserOffGcode  = "M5 S0"
	plungeComment  = "; plunge"
	retractComment = "; Retract"
)

type Conf struct {
	Power int `cfgp:"power,% of laser power [0-100],"`
}

// parsePower transform laser power from [0-100]% range
// to S[0-100] gcode parameter.
func parsePower(p int) (string, error) {
	if p < 0 || p > 100 {
		return "", fmt.Errorf("invalid power value: %d", p)
	}
	i := strconv.Itoa(p)
	return "S" + i, nil
}

func main() {
	c := Conf{}
	err := cfgp.Parse(&c)
	if err != nil {
		log.Fatal("Unable to parse configuration", err)
	}
	toPlunge := false
	toRetract := false
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		line := s.Text()
		if toPlunge {
			p, err := parsePower(c.Power)
			if err != nil {
				log.Fatalln(err)
			}
			fmt.Println(laserOnGcode, p)
			fmt.Println(delayGcode)
			toPlunge = false
		}
		if toRetract {
			fmt.Println(laserOffGcode)
			toRetract = false
		}
		matched, err := regexp.MatchString(plungeComment, line)
		if err != nil {
			log.Fatalln(err)
		}
		if matched {
			toPlunge = true
			continue
		}
		matched, err = regexp.MatchString(retractComment, line)
		if err != nil {
			log.Fatalln(err)
		}
		if matched {
			toRetract = true
			continue
		}
		// If line contains moves on Z axis skip it.
		matched, err = regexp.MatchString(".*Z.*", line)
		if err != nil {
			log.Fatalln(err)
		}
		if matched {
			continue
		}
		fmt.Println(line)
	}
	// Error parsing stdin.
	if err := s.Err(); err != nil {
		log.Fatalln(err)
	}
	// The End
	fmt.Println("M2")
}
