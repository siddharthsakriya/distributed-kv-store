package raft.state;

public class VolatileState {
    private int commitIndex;
    private int lastApplied;
    private long leaderId;

    public VolatileState(int commitIndex, int lastApplied, long leaderId) {
        this.commitIndex = commitIndex;
        this.lastApplied = lastApplied;
        this.leaderId = leaderId;
    }

    public int getCommitIndex() {
        return commitIndex;
    }

    public int getLastApplied() {
        return lastApplied;
    }

    public long getLeaderId() {
        return leaderId;
    }

    public void setLeaderId(long leaderId) {
        this.leaderId = leaderId;
    }

    public void setCommitIndex(int commitIndex) {
        this.commitIndex = commitIndex;
    }

    public void setLastApplied(int lastApplied) {
        this.lastApplied = lastApplied;
    }
}
