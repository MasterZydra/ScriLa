package bashAssembler

import (
	"ScriLa/cmd/scrila/config"
	"fmt"
	"strings"
)

func (self *Assembler) writeLnWithTabsToFile(content string) {
	self.writeLnToFile(fmt.Sprintf("%s%s", self.tabs(), content))
}

func (self *Assembler) writeLnToFile(content string) {
	self.writeToFile(fmt.Sprintf("%s\n", content))
}

func (self *Assembler) writeWithTabsToFile(content string) {
	self.writeToFile(fmt.Sprintf("%s%s", self.tabs(), content))
}

func (self *Assembler) writeToFile(content string) {
	if self.testPrintMode {
		fmt.Print(content)
	}
	if self.testMode || self.testPrintMode {
		return
	}
	if self.outputFile != nil {
		self.outputFile.WriteString(content)
	}
}

func (self *Assembler) writeFileHeader() {
	self.writeLnToFile("#!/bin/bash")
	// Do not write version number in test print mode
	// because this would require a change in the tests with every version change.
	if !self.testPrintMode {
		self.writeLnToFile(fmt.Sprintf("# Created by Scrila Transpiler %s", config.Version))
	}
	self.writeLnToFile("")
}

// Get the tabs for the correct indentation
func (self *Assembler) tabs() string {
	return strings.Repeat("\t", self.indentDepth)
}

func (self *Assembler) decTabs() {
	self.indentDepth--
}

func (self *Assembler) incTabs() {
	self.indentDepth++
}
