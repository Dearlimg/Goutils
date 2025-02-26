package sonyflake

import (
	"fmt"
	"github.com/sony/sonyflake"
	"time"
)

var (
	sonyFlake     *sonyflake.Sonyflake
	sonyMachineID uint16
)

func getMachineID() (uint16, error) {
	return sonyMachineID, nil
}

func Init(startTime string, machineID uint16) (err error) {
	sonyMachineID = machineID
	st, err := time.Parse("2006-01-02", startTime)
	if err != nil {
		return err
	}
	settings := sonyflake.Settings{
		StartTime: st,
		MachineID: getMachineID,
	}
	sonyFlake = sonyflake.NewSonyflake(settings)
	return nil
}

// GetID 生成ID
func GetID() (id uint64, err error) {
	if sonyFlake == nil {
		err = fmt.Errorf("sony failed, err:%v\n", err)
		return
	}
	id, err = sonyFlake.NextID()
	return
}
