package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"ota/common"
	"ota/parser"
	"strings"
)

// generate rootfs with sylixos directioy structure
func otaCreateRootfs(rootFsPath string) error {
	// subdir lists
	subDirs := []string{
		"apps",
		"bin",
		"boot",
		"etc",
		"home",
		"lib",
		"lib/modules",
		"qt",
		"root",
		"sbin",
		"usr",
	}

	// create subdir
	for _, dir := range subDirs {
		if err := common.CreateDir(rootFsPath + dir); err != nil {
			return err
		}
	}

	// create startup.sh and set stack size with 200000 as default
	// if err := common.CreateFile(rootPath+"/etc/startup.sh", "shstack 200000\n"); err != nil {
	// 	return err
	// }

	log.Println("RootFS generated at", rootFsPath)
	return nil
}

func CopyFiletoRootfs(product parser.Product, bundlePath string) {
	log.Printf("	Product Name: %s, Version: %s, Root Directory: %s\n", product.Name, product.Version, product.RootDir)

	_, err := os.Stat(bundlePath)
	if !os.IsNotExist(err) {
		log.Printf("The bundle %s already exists and will be removed...\n", bundlePath)
		err = os.RemoveAll(bundlePath)
		if err != nil {
			log.Fatalf("The bundle removed faild, please retry.\n")
		}
		log.Println("The bundle has been removed.")
	}

	for _, filter := range product.Filters {
		log.Printf("  Filter Depend: %s\n", filter.Depend)

		if product.Type == "image" {
			otaCreateRootfs(bundlePath + product.RootDir)
		} else {
			common.CreateDir(bundlePath + product.RootDir)
		}

		for _, file := range filter.Files {
			err := common.CopyFile(filter.Depend+"/"+file.Src, bundlePath+file.Des)
			if err != nil {
				log.Printf("Copy File failed from %s to %s erro %s", filter.Depend+"/"+file.Src, bundlePath+file.Des, err)
			}
			log.Printf("Copy File from %s to %s success", filter.Depend+"/"+file.Src, bundlePath+file.Des)
		}
		for _, dir := range filter.Dirs {
			err := common.CopyDir(filter.Depend+"/"+dir.Src, bundlePath+dir.Des)
			if err != nil {
				log.Printf("Copy Dir failed from %s to %s erro %s", filter.Depend+"/"+dir.Src, bundlePath+dir.Des, err)
			}
			log.Printf("Copy Dir from %s to %s success", filter.Depend+"/"+dir.Src, bundlePath+dir.Des)
		}
	}
	CreatePackageTar(product.Name, bundlePath)
}

func CreatePackageTar(layImage string, bundlePath string) error {
	if err := common.TarDirectory("output"+"/"+strings.TrimSuffix(parser.GlobalPackageInfo.Output.Name, ".tar.gz")+"/"+layImage, bundlePath); err != nil {
		log.Println("Failed to create Package tar")
		return err
	}
	log.Println("Create Package tar success, rootfs will be deleted...")
	if err := os.RemoveAll(bundlePath); err != nil {
		log.Fatalln("The bundle removed faild, please retry")
	}
	return nil
}

func createFileListJson(jsonPath string) error {
	for _, product := range parser.GlobalPackageInfo.Output.Products {
		md5Value, err := common.CalculateFileMD5(jsonPath + "/" + product.Name)
		if err != nil {
			log.Println("Failed to calculate md5")
			return err
		}
		outputItem := parser.OutputItem{
			Type:         product.Type,
			Architecture: "aarch64",
			Name:         product.Name,
			Version:      product.Version,
			Md5:          fmt.Sprintf("%x", md5Value),
		}

		parser.GlobalPackageInfoStruct.Outputs = append(parser.GlobalPackageInfoStruct.Outputs, outputItem)
	}

	jsonData, err := json.MarshalIndent(parser.GlobalPackageInfoStruct, "", "    ")
	if err != nil {
		fmt.Println("Error parsing jsonData:", err)
		return err
	}

	if err := os.WriteFile(jsonPath+"/"+"filelist.json", jsonData, 0644); err != nil {
		log.Println("Error writing file:", err)
		return err
	}

	log.Println("filelist.json created successfully.")

	return nil
}

func CreateOutputTarball(outputName string) error {

	createFileListJson("output" + "/" + strings.TrimSuffix(outputName, ".tar.gz"))

	common.TarDirectory(outputName, "output")

	os.RemoveAll(strings.TrimSuffix(outputName, ".tar.gz"))
	os.RemoveAll("output")

	log.Println("filelist.json created successfully.")

	return nil
}
