package raft.state;

import raft.core.RaftLog;

public class PersistentState {
    private int currentTerm;
    private long votedFor;
    private final RaftLog raftLog;
    private static final long NO_VOTE = -1;

    public PersistentState() {
        this.currentTerm = 0;
        this.votedFor = -1;
        this.raftLog = new RaftLog();
    }

    public PersistentState(RaftLog raftLog) {
        this.currentTerm = 0;
        this.votedFor = -1;
        this.raftLog = raftLog;
    }

    public synchronized int getCurrentTerm() {
        return currentTerm;
    }

    public synchronized void setCurrentTerm(int currentTerm) {
        this.currentTerm = currentTerm;
    }

    public synchronized long getVotedFor() {
        return votedFor;
    }

    public synchronized void setVotedFor(long votedFor) {
        this.votedFor = votedFor;
    }

    public synchronized RaftLog getRaftLog() {
        return raftLog;
    }

    public synchronized void incrementTerm() {
        currentTerm++;
    }

    public synchronized void resetVote() {
        votedFor = NO_VOTE;
    }

    // stub methods for later
    public void saveToDisk() { /* TODO: write to file */ }
    public void loadFromDisk() { /* TODO: read from file */ }
}
