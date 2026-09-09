package raft

func (n *Node) Start() {
	// start everything
	go n.runElectionTimer()
	// if channel nil, then blocks for ever so need this check to avoid deadlock:/
	if n.applyCh != nil {
		go n.runApplyLoop()
	}
}

func (n *Node) Stop() {
	n.mu.Lock()
	n.stopped = true
	n.mu.Unlock()
	n.applyCond.Broadcast()
}
