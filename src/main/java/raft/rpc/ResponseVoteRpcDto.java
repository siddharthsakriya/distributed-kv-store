package raft.rpc;

public class ResponseVoteRpcDto {
    private final int term;
    private final boolean voteGranted;

    public ResponseVoteRpcDto(int term, boolean voteGranted) {
        this.term = term;
        this.voteGranted = voteGranted;
    }

    public int getTerm() {
        return term;
    }

    public boolean isVoteGranted() {
        return voteGranted;
    }
}
