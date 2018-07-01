package main

import (
	"golang.org/x/crypto/ssh"
	"log"
	"bytes"
	"sync"
	"time"
	"io/ioutil"
	"os"
	"bufio"
)

var LOGIN string
var PASSWORD string

func main() {
	log.Println("Program is starting...")
	readCreds()
	devices := readDevices()
	log.Println("Connect to these devices:")
	log.Println(devices)
	commands := readCommands()
	log.Printf("Run these commands:\n%v", commands)

	config := &ssh.ClientConfig{
		User: LOGIN,
		Auth: []ssh.AuthMethod{
			ssh.Password(PASSWORD),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout: time.Second * 5,
	}
	var wg sync.WaitGroup
	start := time.Now()
	for k := range devices {
		wg.Add(1)
		device := devices[k] + ":22"
		go sendCommands(config, device, commands, &wg)
	}
	wg.Wait()
	log.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}

func sendCommands(config *ssh.ClientConfig, device string , commands string, wg *sync.WaitGroup){
	log.Printf("Running commands on device: %v", device)
	defer wg.Done()
	client, err := ssh.Dial("tcp", device, config)
	if err != nil {
		log.Println(err.Error())
		log.Printf("Cannot connect to device: %v", device)
		return
	}
	session, err := client.NewSession()
	if err != nil {
		log.Println(err)
	}
	defer session.Close()
	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run(commands); err != nil {
		log.Println(err)
	}
	outputFileName := device[:len(device) -3] + "-output.txt"
	f, err := os.Create(outputFileName)
	if err != nil {
		log.Print(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()
	if _, err := f.WriteString(b.String()); err != nil {
		log.Print(err)
	}
}

func readCreds () {
	file, err := os.Open("credentials.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for i := 0; i <= 1; i++ {
		scanner.Scan()
		if i == 0 {
			LOGIN = scanner.Text()[10:]
			log.Printf("Login: %v", scanner.Text()[10:])
		}
		if i == 1 {
			PASSWORD = scanner.Text()[10:]
			log.Printf("Password: %v", scanner.Text()[10:])
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func readDevices () (devices []string){
	file, err := os.Open("devices.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		devices = append(devices, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return devices
}

func readCommands () (commands string){
	b, err := ioutil.ReadFile("commands.txt")
	if err != nil {
		log.Print(err)
	}
	commands = string(b)
	return commands
}
