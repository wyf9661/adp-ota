package parser

import (
	"encoding/xml"
	"io"
	"log"
	"os"
)

// Package 定义了XML文件中<package>标签的内容结构
type Package struct {
	Inputs []Input `xml:"input"`
	Output Output  `xml:"output"`
}

// Input 定义了XML文件中<input>标签的内容结构
type Input struct {
	Resource string        `xml:"resource,attr"`
	Pkgs     []PackageInfo `xml:"pkg"`
}

// PackageInfo 定义了XML文件中<pkg>标签的内容结构
type PackageInfo struct {
	Name         string `xml:"name,attr"`
	Organization string `xml:"organization,attr"`
	Version      string `xml:"version,attr"`
	Type         string `xml:"type,attr"`
	Arch         string `xml:"arch,attr"`
	Platform     string `xml:"platform,attr"`
	Url          string `xml:"url,attr"`
}

// Output 定义了XML文件中<output>标签的内容结构
type Output struct {
	Name     string    `xml:"name,attr"`
	Version  string    `xml:"version,attr"`
	Products []Product `xml:"product"`
}

// Product 定义了XML文件中<product>标签的内容结构
type Product struct {
	Name    string   `xml:"name,attr"`
	Version string   `xml:"version,attr"`
	Type    string   `xml:"type,attr"`
	Filters []Filter `xml:"filter"`
	RootDir string   `xml:"rootdir"`
}

// Filter 定义了XML文件中<filter>标签的内容结构
type Filter struct {
	Depend string `xml:"depend,attr"`
	Files  []File `xml:"file"`
	Dirs   []Dir  `xml:"dir"`
}

// File 定义了XML文件中<file>标签的内容结构
type File struct {
	Src string `xml:"src,attr"`
	Des string `xml:"des,attr"`
}

// Dir 定义了XML文件中<dir>标签的内容结构
type Dir struct {
	Src string `xml:"src,attr"`
	Des string `xml:"des,attr"`
}

// globalConfig 是一个全局变量，用于存储解析后的XML数据
var GlobalPackageInfo Package

func GetResourceOfConfigFile() {
	// 假设XML文件名为config.xml
	xmlFile, err := os.Open("config/config.xml")
	if err != nil {
		log.Println("Error opening file:", err)
		return
	}
	defer xmlFile.Close()

	xmlData, err := io.ReadAll(xmlFile)
	if err != nil {
		log.Println("Error reading file:", err)
		return
	}

	err = xml.Unmarshal(xmlData, &GlobalPackageInfo)
	if err != nil {
		log.Println("Error unmarshalling XML:", err)
		return
	}

	// 打印解析后的数据，以验证是否正确解析
	log.Printf("Parsed XML Data: \n%+v\n", GlobalPackageInfo)
}
