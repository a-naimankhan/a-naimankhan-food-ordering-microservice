package main

import "log"

func main() {
	configs, err := configs.GetConfigs()
	if err != nil {
		log.Fatal(err)
	}

	startServer(configs)

}

func startServer(configs *configs.Config) {
	orderRepo ,
}
