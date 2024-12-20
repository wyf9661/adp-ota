package armory

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"ota/common"
	"path/filepath"
)

type LoginResponse struct {
	Status      int         `json:"status"`
	Message     string      `json:"message"`
	FieldErrors interface{} `json:"fieldErrors"`
	Data        LoginData   `json:"data"`
}
type LoginData struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	Jti          string `json:"jti"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

var GlobalLoginResp LoginResponse

func GetUserTokenOfArmory(userName string, passWord string) string {
	data := map[string]string{
		"account":  userName,
		"password": passWord,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return ""
	}

	url := "http://10.7.1.31/usercenter/v1/auth/login"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return ""
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return ""
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return ""
	}
	log.Printf("Token Response status %s from armory server status.", resp.Status)

	err = json.Unmarshal(body, &GlobalLoginResp)
	if err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}

	log.Println("Get armory token success!")
	return GlobalLoginResp.Data.AccessToken
}

func DownloadFileFromArmory(srcUrl string, filePath string) error {

	client := &http.Client{}
	req, err := http.NewRequest("GET", srcUrl, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return err
	}

	req.AddCookie(&http.Cookie{Name: "USER_TOKEN", Value: GlobalLoginResp.Data.AccessToken})
	newresp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return err
	}

	defer newresp.Body.Close()
	fmt.Println("Response status:", newresp.Status)

	dir := filepath.Dir(filePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			fmt.Printf("Error creating directory: %v\n", err)
			return err
		}
	}

	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return err
	}
	defer file.Close()

	wc := &common.WriteCounter{}
	writer := bufio.NewWriter(file)
	_, err = io.Copy(writer, io.TeeReader(newresp.Body, wc))
	if err != nil {
		fmt.Println("Error writing file:", err)
		return err
	}
	writer.Flush()

	fmt.Println("File saved successfully:", filePath)
	return nil
}
