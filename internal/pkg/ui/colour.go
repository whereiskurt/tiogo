package ui

type Colour struct {
	BOLD     string
	RESET    string
	GREEN    string
	YELLOW   string
	RED      string
	GREY     string
	GRAY     string
	CRUNNING string
	CNEVER   string
	CCRIT    string
	CHIGH    string
	CMED     string
	CLOW     string
	CINFO    string
	CCNT     string
}

func NewColor() (c *Colour) {
	c = new(Colour)
	c.BOLD = ""
	c.RESET = ""
	c.GREEN = ""
	c.YELLOW = ""
	c.RED = ""
	c.GREY = ""
	c.GRAY = ""
	c.CRUNNING = ""
	c.CNEVER = ""
	c.CCRIT = ""
	c.CHIGH = ""
	c.CMED = ""
	c.CLOW = ""
	c.CINFO = ""
	c.CCNT = ""
	return
}
func (c *Colour) Disable() {
	c.BOLD = ""
	c.RESET = ""
	c.GREEN = ""
	c.YELLOW = ""
	c.RED = ""
	c.GREY = ""
	c.GRAY = ""
	c.CRUNNING = ""
	c.CNEVER = ""
	c.CCRIT = ""
	c.CHIGH = ""
	c.CMED = ""
	c.CLOW = ""
	c.CINFO = ""
	c.CCNT = ""
	return
}
