package chain

type ICollect interface {
	Title(c *Collect)
	Collect(txt string, data *InputData) bool
}

type InputData struct {
	Version     string
	Description string
}
