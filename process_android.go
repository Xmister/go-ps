// +build android

package ps

import (
	"fmt"
	"io/ioutil"
	"strings"
)

// Refresh reloads all the data associated with this process.
func (p *UnixProcess) Refresh() error {
	statPath := fmt.Sprintf("/proc/%d/stat", p.pid)
	dataBytes, err := ioutil.ReadFile(statPath)
	if err != nil {
		return err
	}
	cmdPath := fmt.Sprintf("/proc/%d/cmdline", p.pid)
	cmdBytes, err := ioutil.ReadFile(cmdPath)
	if err != nil {
		return err
	}

	// First, parse out the image name
	cmd  := string(cmdBytes)
	args := strings.SplitN(cmd, string(0x0), 2)
	if len(args) > 0 {
		p.binary = args[0]
	}
	data := string(dataBytes)
	binStart := strings.IndexRune(data, '(') + 1
	binEnd := strings.IndexRune(data[binStart:], ')')
	if p.binary == "" {
		p.binary = data[binStart : binStart+binEnd]
	}

	// Move past the image name and start parsing the rest
	data = data[binStart+binEnd+2:]
	_, err = fmt.Sscanf(data,
		"%c %d %d %d",
		&p.state,
		&p.ppid,
		&p.pgrp,
		&p.sid)

	return err
}
