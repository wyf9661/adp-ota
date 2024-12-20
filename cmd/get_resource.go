package cmd

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"ota/armory"
	"ota/common"
	"path/filepath"
)

func GetSourceFileFromArmory(name string, organization string, version string, arch string, platform string, filetype string) {

	url := "http://10.7.1.31/armory/v1/packages/@" + organization + "/" + name + "/v/" + version + "/p/" + platform + "/a/" + arch

	fmt.Println("	url:", url)

	fileName := name + "_" + version + "." + filetype
	filePath := name

	armory.DownloadFileFromArmory(url, filePath+"/"+fileName)
	fmt.Println("	Download", fileName, "from armory success ,ready to decompress file...")

	switch filepath.Ext(fileName) {
	case ".zip":
		log.Println("unzip file:", fileName)
		err := common.Unzip(filePath+"/"+fileName, filePath)
		if err != nil {
			return
		}
		log.Println("unzip file to", filePath, "success!")
		os.RemoveAll(fileName)
	case ".gz":
		log.Println("decompress file with tar.gz:", fileName)
		err := common.UntarGz(filePath+"/"+fileName, filePath)
		if err != nil {
			return
		}
		log.Println("decompress file with tar.gz to", filePath, "success!")
	}

	err := filepath.Walk(filePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(path) == ".zip" {
			err = common.Unzip(path, filePath)
			if err != nil {
				return err
			}
			log.Println("unzip file to", filePath, "success!")
		} else if !info.IsDir() && filepath.Ext(path) == ".gz" {
			err = common.UntarGz(path, filePath)
			if err != nil {
				return err
			}
			log.Println("decompress file with tar to", filePath, "success!")
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Error walking the path %s: %v\n", filePath, err)
	}
}

func GetSourceFileFromHttp(url string, targetDir string) error {

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned %d status", resp.StatusCode)
	}

	filename := filepath.Base(url)

	fileName := filepath.Join(targetDir, filename)

	err = os.MkdirAll(targetDir, 0755)
	if err != nil {
		return err
	}

	out, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	switch filepath.Ext(filename) {
	case ".zip":
		log.Println("unzip file:", fileName)
		err = common.Unzip(fileName, targetDir)
		if err != nil {
			return err
		}
		log.Println("unzip file to", targetDir, "success!")
		os.RemoveAll(fileName)
	case ".gz":
		log.Println("decompress file with tar.gz:", fileName)
		err = common.UntarGz(fileName, targetDir)
		if err != nil {
			return err
		}
		log.Println("decompress file with tar.gz to", targetDir, "success!")
	}

	err = filepath.Walk(targetDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			ext := filepath.Ext(path)
			switch ext {
			case ".zip":
				log.Println("unzip file:", path)
				err = common.Unzip(path, targetDir)
				if err != nil {
					return err
				}
				log.Println("unzip file to", targetDir, "success!")
			case ".gz":
				log.Println("decompress file with tar.gz:", path)
				err = common.UntarGz(path, targetDir)
				if err != nil {
					return err
				}
				log.Println("decompress file with tar.gz to", targetDir, "success!")
			}
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Error walking the path %s: %v\n", targetDir, err)
		return err
	}

	return nil
}
