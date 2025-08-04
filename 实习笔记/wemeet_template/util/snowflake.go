package util

import (
	"errors"
	"github.com/sony/sonyflake"
	"net"
	"time"
)

var Sf *sonyflake.Sonyflake
var localIps net.IP

//初始化
func init() {
	InitSonyflake()
}

//InitSonyflake 初始化
func InitSonyflake() {
	st := sonyflake.Settings{}
	st.MachineID = getMachineID

	Sf = sonyflake.NewSonyflake(st)
	if Sf == nil {
		panic("sonyflake not created")
	}
}

//GenId 获取任务ID
func GenId() int64 {
	unique_id, err := Sf.NextID()
	if err != nil {
		unique_id = uint64(time.Now().UnixNano())
	}
	return int64(unique_id)
}

//GenTimeByID 通过unique_id 解记录生成时间，精度秒
func GenTimeByID(id int64) uint64 {
	setMap := sonyflake.Decompose(uint64(id))
	startTime := setMap["time"] //10毫秒为单位
	a := startTime / 100        //转化为秒
	//StartTime是将Sonyflake时间定义为经过时间的时间,Sonyflake的开始时间设置为“2014-09-01 00:00:00+0000 UTC
	b := time.Date(2014, 9, 1, 0, 0, 0, 0, time.UTC).UTC().Unix()
	c := a + uint64(b)
	return c
}

//isLocalIPv4
func isLocalIPv4(ip net.IP) bool {
	return ip == nil || ip[0] == 127
}

//privateIPv4
func privateIPv4() (net.IP, error) {
	as, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	for _, a := range as {
		ipNet, ok := a.(*net.IPNet)
		if !ok || ipNet.IP.IsLoopback() {
			continue
		}

		ip := ipNet.IP.To4()
		if !isLocalIPv4(ip) {
			return ip, nil
		}
	}
	return nil, errors.New("no private ip address")
}

//getMachineID
func getMachineID() (uint16, error) {
	var instanceId uint16 = 0

	ip, err := privateIPv4()
	if err != nil {
		return 0, err
	}
	localIps = ip

	instanceId = uint16(ip[2])<<8 + uint16(ip[3])

	return instanceId, nil
}
