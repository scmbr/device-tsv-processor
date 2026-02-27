package domain

type UnitMessage struct {
	GUID      string
	UnitGUID  string
	MsgID     string
	Text      string
	Context   string
	Class     string
	Level     int
	Area      string
	Addr      string
	Block     string
	Type      string
	Bit       int
	InvertBit bool
}
