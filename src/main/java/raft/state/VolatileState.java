package raft.state;

public class VolatileState {
    private final int commitIndex;
    private final int lastApplied;

    public VolatileState(int commitIndex, int lastApplied) {
        this.commitIndex = commitIndex;
        this.lastApplied = lastApplied;
    }

    public int getCommitIndex() {
        return commitIndex;
    }

    public int getLastApplied() {
        return lastApplied;
    }
}
