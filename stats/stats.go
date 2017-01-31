package stats

import (
	"sync"
	"time"
)

type StatResponse struct {
	MsgTotal        int     `json:"msg_total"`
	MsgReqTotal     int     `json:"msg_req_total"`
	MsgAckTotal     int     `json:"msg_ack_total"`
	MsgNakTotal     int     `json:"msg_nak_total"`
	RequestRate1s   float64 `json:"request_rate1s"`
	RequestRate10s  float64 `json:"request_rate10s"`
	ResponseRate1s  float64 `json:"response_rate1s"`
	ResponseRate10s float64 `json:"response_rate10s"`
}

type StatRecord struct {
	MsgTotal        int
	MsgReqTotal     int
	MsgAckTotal     int
	MsgNakTotal     int
	RequestRecords  []int
	ResponseRecords []int
}

var mutex = &sync.Mutex{}

func NewStatRecorder() *StatRecord {
	sr := &StatRecord{}

	sr.RequestRecords = append(sr.RequestRecords, 0)
	sr.ResponseRecords = append(sr.ResponseRecords, 0)
	ticker := time.NewTicker(time.Second)
	go func() {
		for range ticker.C {
			mutex.Lock()
			sr.RequestRecords = append(sr.RequestRecords, 0)
			sr.ResponseRecords = append(sr.ResponseRecords, 0)
			mutex.Unlock()
		}
	}()

	return sr
}

func (sr *StatRecord) RecordReq() {
	mutex.Lock()
	sr.MsgTotal = sr.MsgTotal + 1
	sr.MsgReqTotal = sr.MsgReqTotal + 1
	sr.RequestRecords[len(sr.RequestRecords)-1] = sr.RequestRecords[len(sr.RequestRecords)-1] + 1
	mutex.Unlock()
}

func (sr *StatRecord) RecordAck() {
	mutex.Lock()
	sr.MsgTotal = sr.MsgTotal + 1
	sr.MsgAckTotal = sr.MsgAckTotal + 1
	sr.ResponseRecords[len(sr.ResponseRecords)-1] = sr.ResponseRecords[len(sr.ResponseRecords)-1] + 1
	mutex.Unlock()
}

func (sr *StatRecord) RecordNak() {
	mutex.Lock()
	sr.MsgTotal = sr.MsgTotal + 1
	sr.MsgNakTotal = sr.MsgNakTotal + 1
	sr.ResponseRecords[len(sr.ResponseRecords)-1] = sr.ResponseRecords[len(sr.ResponseRecords)-1] + 1
	mutex.Unlock()
}

func (sr *StatRecord) StatResponse() StatResponse {
	mutex.Lock()
	res := StatResponse{
		MsgTotal:        sr.MsgTotal,
		MsgReqTotal:     sr.MsgReqTotal,
		MsgAckTotal:     sr.MsgAckTotal,
		MsgNakTotal:     sr.MsgNakTotal,
		RequestRate1s:   avg(sr.RequestRecords),
		RequestRate10s:  avg10(sr.RequestRecords),
		ResponseRate1s:  avg(sr.ResponseRecords),
		ResponseRate10s: avg10(sr.ResponseRecords),
	}
	mutex.Unlock()
	return res
}
func avg(xs []int) float64 {
	total := float64(0)
	for _, x := range xs {
		total += float64(x)
	}
	return total / float64(len(xs))
}

func avgFloat64(xs []float64) float64 {
	total := float64(0)
	for _, x := range xs {
		total += x
	}
	return total / float64(len(xs))
}

func avg10(xs []int) float64 {
	if len(xs) < 10 {
		return avg(xs)
	}

	rollingAvg := []float64{}

	for i := 1; i <= (len(xs) - 10); i++ {
		w := xs[i : i+10]
		rollingAvg = append(rollingAvg, avg(w))
	}

	return avgFloat64(rollingAvg)
}
