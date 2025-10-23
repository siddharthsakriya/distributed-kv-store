package raft.core;

import java.util.List;

public class RaftConfig {
    private long selfId;
    private List<Long> peerIds;
    private int electionTimeoutMin;
    private int electionTimeoutMax;
    private int heartBeatInterval;

    public RaftConfig(long selfId, List<Long> peerIds, int electionTimeoutMin, int electionTimeoutMax, int heartBeatInterval) {
        this.selfId = selfId;
        this.peerIds = peerIds;
        this.electionTimeoutMin = electionTimeoutMin;
        this.electionTimeoutMax = electionTimeoutMax;
        this.heartBeatInterval = heartBeatInterval;
    }

    public long getSelfId() {
        return selfId;
    }

    public void setSelfId(long selfId) {
        this.selfId = selfId;
    }

    public List<Long> getPeerIds() {
        return peerIds;
    }

    public void setPeerIds(List<Long> peerIds) {
        this.peerIds = peerIds;
    }

    public int getElectionTimeoutMin() {
        return electionTimeoutMin;
    }

    public void setElectionTimeoutMin(int electionTimeoutMin) {
        this.electionTimeoutMin = electionTimeoutMin;
    }

    public int getElectionTimeoutMax() {
        return electionTimeoutMax;
    }

    public void setElectionTimeoutMax(int electionTimeoutMax) {
        this.electionTimeoutMax = electionTimeoutMax;
    }

    public int getHeartBeatInterval() {
        return heartBeatInterval;
    }

    public void setHeartBeatInterval(int heartBeatInterval) {
        this.heartBeatInterval = heartBeatInterval;
    }
}
