package main

import "fmt"

func main() {

	for i := 0; i < 3; i++ {
		go setUpServers(5000 + int32(i)) //sets up servers through setupServers()
	}

}

func setupServers(portNum int32) {

	fmt.Printf("lololol")
}
