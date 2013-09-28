package main

import (
	"sync"
)

// Hold on to the clients - lock so that no weird operations occur where we try to
// update a client right after deleting it from the list or something.
type ClientHolder struct {
	holder []*ChatClient
	lock   *sync.RWMutex
}

// Make sure that a new clientholder's slice is setup
func (this *ClientHolder) init() {
	this.holder = make([]*ChatClient, 0)
	this.lock = new(sync.RWMutex)

}

// Gets the underlying slice.  Not sure if there's a point to locking
// this but will try it out.
func (this *ClientHolder) getClients() []*ChatClient {
	this.lock.RLock()
	defer this.lock.RUnlock()
	return this.holder
}

// Adds a new client to the slice
func (this *ClientHolder) addClient(client *ChatClient) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.holder = append(this.holder, client)
}

// Removes a client (by it's ID) from the slice
func (this *ClientHolder) removeClient(id uint64) {
	newHolder := make([]*ChatClient, 0)
	this.lock.RLock()
	for _, c := range this.holder {
		if c.id != id {
			newHolder = append(newHolder, c)
		}
	}
	this.lock.RUnlock()
	this.lock.Lock()
	this.holder = newHolder
	this.lock.Unlock()
}

// Search the holder for a client with the passed in id.  Returns
// nil if not found.
func (this *ClientHolder) getClientById(id uint64) *ChatClient {
	var client *ChatClient
	this.lock.RLock()
	defer this.lock.RUnlock()
	for _, c := range this.holder {
		if c.id == id {
			client = c
			break
		}
	}
	return client
}

// Search the holder for a client with the passed in nick.  Returns
// nil if not found.
func (this *ClientHolder) getClientByNick(nick string) *ChatClient {
	var client *ChatClient
	this.lock.RLock()
	defer this.lock.RUnlock()
	for _, c := range this.holder {
		if c.nick == nick {
			client = c
			break
		}
	}
	return client
}
