package main

import (
	"log"
	"os"
	"sync"

	evdev "github.com/gvalkov/golang-evdev"
)

var (
	phys           string
	switchBotToken string
	plugID         string
)

func watch(wg *sync.WaitGroup, dev *evdev.InputDevice) {
	defer wg.Done()

	c := newSwitchBotClient(switchBotToken)

	for {
		ev, err := dev.ReadOne()
		if err != nil {
			log.Println(err)
			continue
		}
		if ev.Type != evdev.EV_KEY {
			continue
		}
		if ev.Value == 0 {
			switch ev.Code {
			case evdev.KEY_VOLUMEUP:
				// upper button
				if err := c.switching(plugID); err != nil {
					log.Println(err)
				}
			case evdev.KEY_ENTER:
				// lower button
			}
		}
	}
}

func main() {
	var ok bool
	phys, ok = os.LookupEnv("DEVICE")
	if !ok {
		log.Fatal("must specify DEVICE")
	}

	switchBotToken, ok = os.LookupEnv("TOKEN")
	if !ok {
		log.Fatal("must specify TOKEN")
	}

	plugID, ok = os.LookupEnv("PLUG_ID")
	if !ok {
		log.Fatal("must specify PLUG_ID")
	}

	devs, err := evdev.ListInputDevices()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		for _, dev := range devs {
			dev.File.Close()
		}
	}()

	var wg sync.WaitGroup
	n := 0

	for _, dev := range devs {
		if dev.Phys == phys {
			wg.Add(1)
			n++
			go watch(&wg, dev)
		}
	}

	if n == 0 {
		log.Fatal("cannot open uinput device")
	}

	wg.Wait()
}
