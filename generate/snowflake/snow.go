package snowflake

import (
	"github.com/sony/sonyflake"
)

var deft = sonyflake.NewSonyflake(sonyflake.Settings{})

func NextID() (uint64, error) {
	return deft.NextID()
}

func Config(setters ...func(settings *sonyflake.Settings)) {
	setting := sonyflake.Settings{}
	for _, setter := range setters {
		setter(&setting)
	}
	deft = sonyflake.NewSonyflake(setting)
}
