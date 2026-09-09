package raft

type ApplyMsg struct {
	Index   int
	Command []byte
}

func (n *Node) runApplyLoop() {
	for {
		n.mu.Lock()
		for !(n.commitIndex > n.lastApplied) {
			if n.stopped {
				n.mu.Unlock()
				return
			}
			n.applyCond.Wait()
		}
		msgs := []ApplyMsg{}
		for i := n.lastApplied + 1; i <= n.commitIndex; i++ {
			entry := n.entryAt(i)
			msgs = append(msgs, ApplyMsg{
				Index:   entry.Index,
				Command: entry.Command,
			})
		}
		n.lastApplied = n.commitIndex
		n.mu.Unlock()

		for _, msgToApply := range msgs {
			n.applyCh <- msgToApply
		}
	}
}
