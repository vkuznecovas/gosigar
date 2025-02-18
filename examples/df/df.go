// Copyright (c) 2012 VMware, Inc.

package main

import (
	"fmt"
	"os"

	"github.com/vkuznecovas/gosigar"
)

const output_format = "%-15s %4s %4s %5s %4s %-15s\n"

func main() {
	fslist := gosigar.FileSystemList{}
	err := fslist.Get()
	if err != nil {
		fmt.Printf("Failed to get list of filesystems: %v", err)
		os.Exit(-1)
	}

	fmt.Fprintf(os.Stdout, output_format,
		"Filesystem", "Size", "Used", "Avail", "Use%", "Mounted on")

	for _, fs := range fslist.List {
		dir_name := fs.DirName

		usage := gosigar.FileSystemUsage{}

		usage.Get(dir_name)

		fmt.Fprintf(os.Stdout, output_format,
			fs.DevName,
			gosigar.FormatSize(usage.Total),
			gosigar.FormatSize(usage.Used),
			gosigar.FormatSize(usage.Avail),
			gosigar.FormatPercent(usage.UsePercent()),
			dir_name)
	}
}
