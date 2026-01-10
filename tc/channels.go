package tc

import(
	"sync"
)

type ChannelManager struct {
	mu sync.RWMutex
	Channels map[string]Channel
}

type Channel struct {
	mu sync.RWMutex
	Permissions map[string]ChannelPermission
	Connected map[string]any
}

func CreateChannels() *ChannelManager {
	returnObject := &ChannelManager {
		Channels: make(map[string]Channel),
	}
	return returnObject
}

func(c *Channel) AddUser(username string, connection any) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.Connected[username] = connection
}

func (c *Channel) RemoveByUsername(username string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for k, _ := range c.Connected { 
		if k == username {
			delete(c.Connected, username)
			break
		}
	}

}

func (c *Channel) SendMessage(username string, message string) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	sanitized, err := ValidateMessage(message, 120)
	if err != nil {
		return err
	}

	for _, value := range c.Connected {
		
	}

	return nil
}