package accessory

import (
	"github.com/brutella/hc/accessory"
)

type Bridge struct {
	*accessory.Accessory
}

func NewBridge() *Bridge {
	acc := Bridge{}
	info := accessory.Info{
		Name:         "PirageKit",
		Manufacturer: "ipstatic",
		Model:        "Various",
	}
	acc.Accessory = accessory.New(info, accessory.TypeBridge)

	return &acc
}
