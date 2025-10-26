package raft.core;

import raft.rpc.RaftRpcClient;
import raft.rpc.RaftRpcServer;
import raft.rpc.RequestVoteRpcDto;
import raft.rpc.ResponseVoteRpcDto;
import raft.state.LeaderState;
import raft.state.PersistentState;
import raft.state.RaftRole;
import raft.state.VolatileState;

public class RaftNode {
    private RaftRpcServer raftRpcServer;
    private RaftRpcClient raftRpcClient;
    private LeaderState leaderState;
    private VolatileState volatileState;
    private PersistentState persistentState;
    private RaftRole raftRole;
    private RaftConfig raftConfig;

    public RaftNode(RaftConfig raftConfig, PersistentState persistentState) {
        raftRpcServer = new RaftRpcServer();
        raftRpcClient = new RaftRpcClient();
        leaderState = null;
        volatileState = new VolatileState(0,0,-1);
        raftRole = RaftRole.FOLLOWER;
        this.persistentState = persistentState;
        this.raftConfig = raftConfig;
    }

    public RaftRpcServer getRaftRpcServer() {
        return raftRpcServer;
    }

    public RaftRpcClient getRaftRpcClient() {
        return raftRpcClient;
    }

    public LeaderState getLeaderState() {
        return leaderState;
    }

    public VolatileState getVolatileState() {
        return volatileState;
    }

    public PersistentState getPersistentState() {
        return persistentState;
    }

    public RaftRole getRaftRole() {
        return raftRole;
    }

    public RaftConfig getRaftConfig() {
        return raftConfig;
    }

    private void setRaftRole(RaftRole raftRole) {
        this.raftRole = raftRole;
    }

    public void becomeFollower(int term) {
        if (term > persistentState.getCurrentTerm()) {
            persistentState.setCurrentTerm(term);
            persistentState.resetVote();
        }
        setRaftRole(RaftRole.FOLLOWER);
        leaderState = null;
    }

    public void becomeCandidate() {
        setRaftRole(RaftRole.CANDIDATE);
        //TODO: need to add some logic to increment term, vote for self and start voting
    }

    public void becomeLeader() {
        setRaftRole(RaftRole.LEADER);
    }

    public ResponseVoteRpcDto handleRequestVote(RequestVoteRpcDto requestVoteRpcDto) {
        if (persistentState.getCurrentTerm() < requestVoteRpcDto.getTerm()) {
            becomeFollower(requestVoteRpcDto.getTerm());
        }

        if (requestVoteRpcDto.getTerm() < persistentState.getCurrentTerm()) {
            return new ResponseVoteRpcDto(persistentState.getCurrentTerm(), false);
        }

        if (persistentState.getVotedFor() != -1 && persistentState.getVotedFor() != requestVoteRpcDto.getCandidateId()) {
            return new ResponseVoteRpcDto(persistentState.getCurrentTerm(), false);
        }

        if (persistentState.getRaftLog().lastTerm() > requestVoteRpcDto.getLastLogTerm() || (
                persistentState.getRaftLog().lastTerm() == requestVoteRpcDto.getLastLogTerm() &&
                        persistentState.getRaftLog().lastIndex() > requestVoteRpcDto.getLastLogIndex()
                )) {
            return new ResponseVoteRpcDto(persistentState.getCurrentTerm(), false);
        }

        persistentState.setVotedFor(requestVoteRpcDto.getCandidateId());

        return new ResponseVoteRpcDto(persistentState.getCurrentTerm(), true);
    }
}
