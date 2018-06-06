package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	req, err := http.NewRequest(http.MethodGet, "https://api.github.com/users/defunkt", nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	fmt.Printf("%s\n", body)

	user := struct {
		Login     string `json:"login"`
		ID        int    `json:"id"`
		SiteAdmin bool   `json:"site_admin"`
		Bio       string `json:"bio"`
	}{}

	if err := json.Unmarshal(body, &user); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Login: ", user.Login)
	fmt.Println("ID: ", user.ID)
	fmt.Println("SiteAdmin: ", user.SiteAdmin)
	fmt.Println("Bio: ", user.Bio)
}
