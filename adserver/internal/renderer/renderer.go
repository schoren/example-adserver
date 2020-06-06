package renderer

// Renderer can return a string representation of an ad, ready to be served
type Renderer interface {
	Render() string
}
