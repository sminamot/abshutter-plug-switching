package main

import (
	"context"
	"log"
	"time"

	"github.com/nasa9084/go-switchbot"
)

type client struct {
	*switchbot.Client
}

func newSwitchBotClient(t string) *client {
	return &client{
		switchbot.New(switchBotToken),
	}
}

func (c *client) switching(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	s, err := c.Device().Status(ctx, id)
	if err != nil {
		return err
	}

	var command switchbot.Command
	switch s.Power.ToLower() {
	case switchbot.PowerOff.ToLower():
		command = switchbot.TurnOn()
		log.Println("turn on")
	case switchbot.PowerOn.ToLower():
		command = switchbot.TurnOff()
		log.Println("turn off")
	}

	return c.Device().Command(ctx, id, command)
}
