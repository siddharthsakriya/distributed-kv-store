package raft.rpc;

public class RequestVoteRpcDto {
    private final int term;
    private final long candidateId;
    private final int lastLogIndex;
    private final int lastLogTerm;

    public RequestVoteRpcDto(int term, long candidateId, int lastLogIndex, int lastLogTerm) {
        this.term = term;
        this.candidateId = candidateId;
        this.lastLogIndex = lastLogIndex;
        this.lastLogTerm = lastLogTerm;
    }

    public int getTerm() {
        return term;
    }

    public long getCandidateId() {
        return candidateId;
    }

    public int getLastLogIndex() {
        return lastLogIndex;
    }

    public int getLastLogTerm() {
        return lastLogTerm;
    }
}
