package main

import (
	"log"
	"os"
	"ota/armory"
	"ota/cmd"
	"ota/parser"
	"strings"
)

func main() {

	log.Println("Geting Resource Info from config.xml...")
	parser.GetResourceOfConfigFile()

	// download source file from armory
	for _, input := range parser.GlobalPackageInfo.Inputs {
		log.Printf("Input Resource: %s\n", input.Resource)
		if input.Resource == "armory" {
			armory.GetUserTokenOfArmory("wangyifan", "123456")
			for _, pkg := range input.Pkgs {
				log.Printf("  Downloading Package %s, Version: %s from armory...\n", pkg.Name, pkg.Version)
				cmd.GetSourceFileFromArmory(pkg.Name, pkg.Organization, pkg.Version, pkg.Arch, pkg.Platform, pkg.Type)
			}
		}
		if input.Resource == "httpserver" {
			for _, pkg := range input.Pkgs {
				log.Printf("  Downloading Package %s, Url: %s from httpserver...\n", pkg.Name, pkg.Url)
				cmd.GetSourceFileFromHttp(pkg.Url, pkg.Name)
			}
		}
	}

	// create output package
	log.Printf("Output Name: %s\n", parser.GlobalPackageInfo.Output.Name)
	for _, product := range parser.GlobalPackageInfo.Output.Products {
		bundlePath := "output" + "/" + strings.TrimSuffix(parser.GlobalPackageInfo.Output.Name, ".tar.gz") + "/" + strings.TrimSuffix(product.Name, ".tar.gz")
		cmd.CopyFiletoRootfs(product, bundlePath)
	}

	cmd.CreateOutputTarball(parser.GlobalPackageInfo.Output.Name)
	log.Println("Create Package tar success!")

	// cleanup input resource
	for _, input := range parser.GlobalPackageInfo.Inputs {
		if input.Resource == "armory" || input.Resource == "httpserver" {
			for _, pkg := range input.Pkgs {
				os.RemoveAll(pkg.Name)
			}
		}
	}
}
