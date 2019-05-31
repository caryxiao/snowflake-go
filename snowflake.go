package snowflake_go

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

const (
	workerIdBit     uint8 = 10
	seqNumberBit    uint8 = 12
	maxWorkerId     int64 = -1 ^ (-1 << workerIdBit)   // 最大机器ID
	maxSeqNumberBit int64 = -1 ^ (-1 << seqNumberBit)  //自增序列最大值
	timeLO                = workerIdBit + seqNumberBit //时间偏移量
	workerIdLO            = seqNumberBit               //工作机器ID的偏移量
)

// 雪花算法简介
// 总共64位, 最高位为0 + 41位使用毫秒级时间戳+10位的机器ID+12位的自增序列
type Snowflake struct {
	timestamp int64      // 毫秒级时间戳 41bit
	workerId  int64      // 工作机器的ID 10bit
	seqNumber int64      // 自增序列 12bit 每秒能生成4096个
	mu        sync.Mutex //互斥锁
}

func (s *Snowflake) GetId() (id int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now().UnixNano() / 1e6 // 纳秒转毫秒
	if s.timestamp == now {
		if s.seqNumber >= maxSeqNumberBit {
			for now <= s.timestamp {
				now = time.Now().UnixNano() / 1e6
			}
		}
		s.seqNumber++
	} else {
		s.seqNumber = 0
	}

	s.timestamp = now
	id = s.timestamp<<timeLO | (s.workerId << workerIdLO) | s.seqNumber
	return
}

func New(workerId int64) (s *Snowflake, err error) {
	if workerId < 0 || workerId > maxWorkerId {
		return nil, errors.New(fmt.Sprintf("Worker ID excess of quantity: %d", maxWorkerId))
	}
	return &Snowflake{
		workerId: workerId,
	}, nil
}
