package msgtransfer

import (
	"easy-chat/apps/im/ws/ws"
	"easy-chat/pkg/constants"
	"github.com/zeromicro/go-zero/core/logx"
	"sync"
	"time"
)

type groupMsgRead struct {
	mu sync.Mutex

	conversationId string // 会话ID

	push     *ws.Push      // 用于记录消息
	pushChan chan *ws.Push // 推送消息

	count    int       // 计数
	pushTime time.Time // 上次推送时间

	done chan struct{}
}

func newGroupMsgRead(push *ws.Push, pushChan chan *ws.Push) *groupMsgRead {
	m := &groupMsgRead{
		conversationId: push.ConversationId,
		push:           push,
		pushChan:       pushChan,
		count:          1,
		pushTime:       time.Now(),
		done:           make(chan struct{}),
	}

	go m.transfer()
	return m
}

// mergePush 合并消息
func (g *groupMsgRead) mergePush(push *ws.Push) {
	g.mu.Lock()
	defer g.mu.Unlock()
	// 说明已经被清理，重新设置
	if g.push == nil {
		g.push = push
	}

	g.count++
	for msgId, read := range push.ReadRecords {
		g.push.ReadRecords[msgId] = read
	}
}

func (g *groupMsgRead) transfer() {
	// 超时发送
	timer := time.NewTimer(GroupMsgReadRecordDelayTime / 2)
	defer timer.Stop()
	for {
		select {
		case <-g.done:
			return
		case <-timer.C:
			g.mu.Lock()
			// 获取上一次推送时间
			lastPushTime := g.pushTime
			// *2，防止在推送时刻由于阻塞未推送
			val := GroupMsgReadRecordDelayTime*2 - time.Since(lastPushTime)
			// 得到待推送的数据
			push := g.push
			// 当前没有超时且没有超过最大计数，或数据为空
			if val > 0 && g.count < GroupMsgReadRecordDelayCount || push == nil {
				// 重置定时器
				if val > 0 {
					timer.Reset(val)
				}
				// 未达标
				g.mu.Unlock()
				continue
			}
			// 达标
			g.pushTime = time.Now()
			g.push = nil
			g.count = 0
			timer.Reset(GroupMsgReadRecordDelayTime / 2)
			g.mu.Unlock()
			// 推送
			logx.Infof("merge push delay time condition reached, push: %v ", push)
			g.pushChan <- push
		default:
			g.mu.Lock()
			if g.count >= GroupMsgReadRecordDelayCount {
				// 达标，推送
				// 得到待推送的数据
				push := g.push
				g.push = nil
				g.count = 0
				g.mu.Unlock()
				// 推送
				logx.Infof("merge push max delay count condition reached, push: %v ", push)
				g.pushChan <- push
				continue
			}
			if g.isIdle() {
				g.mu.Unlock()
				// 使得 msgReadTransfer 释放
				g.pushChan <- &ws.Push{
					ChatType:       constants.GroupChatType,
					ConversationId: g.conversationId,
				}
				continue
			}
			g.mu.Unlock()
			// 睡眠等待一段时间，最大为一秒钟
			tempDelay := GroupMsgReadRecordDelayTime / 4
			if tempDelay > time.Second {
				tempDelay = time.Second
			}
			time.Sleep(tempDelay)
		}
	}
}

// IsIdle 判断是否为活跃状态
func (g *groupMsgRead) IsIdle() bool {
	g.mu.Lock()
	defer g.mu.Unlock()

	return g.isIdle()
}

func (g *groupMsgRead) isIdle() bool {
	// 获取上一次推送时间
	lastPushTime := g.pushTime
	// *2，防止在推送时刻由于阻塞未推送
	val := GroupMsgReadRecordDelayTime*2 - time.Since(lastPushTime)
	if val <= 0 && g.push == nil && g.count == 0 {
		return true
	}
	return false
}

func (m *groupMsgRead) Clear() {
	select {
	case <-m.done:
	default:
		close(m.done)
	}

	m.push = nil
}
