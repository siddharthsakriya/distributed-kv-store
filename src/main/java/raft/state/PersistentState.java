package raft.state;

import raft.core.RaftLog;

public class PersistentState {
    int currentTerm;
    long votedFor;
    RaftLog raftLog;

    public PersistentState(int currentTerm, long votedFor, RaftLog raftLog) {
        this.currentTerm = currentTerm;
        this.votedFor = votedFor;
        this.raftLog = raftLog;
    }

    public int getCurrentTerm() {
        return currentTerm;
    }

    public void setCurrentTerm(int currentTerm) {
        this.currentTerm = currentTerm;
    }

    public long getVotedFor() {
        return votedFor;
    }

    public void setVotedFor(long votedFor) {
        this.votedFor = votedFor;
    }

    public RaftLog getRaftLog() {
        return raftLog;
    }

    public void setRaftLog(RaftLog raftLog) {
        this.raftLog = raftLog;
    }
}
