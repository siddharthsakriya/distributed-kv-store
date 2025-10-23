package raft.core;

import raft.rpc.RaftRpcClient;
import raft.rpc.RaftRpcServer;
import raft.state.LeaderState;
import raft.state.PersistentState;
import raft.state.RaftRole;
import raft.state.VolatileState;

public class RaftNode {
    private static int nextId = 1;
    private final long nodeId;
    private RaftRpcServer raftRpcServer;
    private RaftRpcClient raftRpcClient;
    private LeaderState leaderState;
    private VolatileState volatileState;
    private PersistentState persistentState;
    private RaftRole raftRole;
    private RaftConfig raftConfig;

    public RaftNode() {
        this.nodeId = RaftNode.nextId;
        nextId++;
        raftRpcServer = new RaftRpcServer();
        raftRpcClient = new RaftRpcClient();
        leaderState = null;
        volatileState = new VolatileState(0,0,-1);
        persistentState = new PersistentState();
        raftRole = RaftRole.FOLLOWER;
        raftConfig = null;
    }



    public static int getNextId() {
        return nextId;
    }

    public long getNodeId() {
        return nodeId;
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
}
