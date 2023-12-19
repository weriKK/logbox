package common

import "log"

type client struct {
	notify chan struct{}
}

type ClientNotifier struct {
	clients []client
}

func NewClientNotifier() *ClientNotifier {
	return &ClientNotifier{
		clients: make([]client, 0),
	}
}

func (cm *ClientNotifier) NotifyAll() {
	log.Println("Notifying all Clients...")

	for _, client := range cm.clients {
		client.notify <- struct{}{}
	}
}

func (cm *ClientNotifier) Register(notifChannel chan struct{}) {
	cm.clients = append(cm.clients, client{notify: notifChannel})
}
