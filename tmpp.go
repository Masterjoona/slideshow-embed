package main

/*
func main() {
	client := &http.Client{}
	var data = strings.NewReader(`{"id":"https://www.tiktok.com/@polishfemboy_/photo/7349697632761072929","hash":"ecd021182575b80e1d389ed596d18e91","mode":"audio","locale":"en","loading_indicator_url":"https://ttsave.app/images/slow-down.gif","unlock_url":"https://ttsave.app/en/unlock"}`)
	req, err := http.NewRequest("POST", "https://api.ttsave.app/", data)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("content-type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", bodyText)
}
*/
