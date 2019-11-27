package upgrade

// NewChannel creates a new channel for the provided URL
func NewChannel(uri string) (Channel, error) {
	channel := &githubChannel{
		uri: uri,
	}
	return channel, channel.resolve()
}
