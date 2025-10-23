package unit.raft.core;

import org.junit.jupiter.api.Test;
import raft.core.RaftNode;
import raft.state.PersistentState;
import raft.state.RaftRole;

import static org.junit.jupiter.api.Assertions.*;

public class RaftNodeTest {
    @Test
    void testBecomeFollower() {

        RaftNode raftNode = new RaftNode(null);

        raftNode.getPersistentState().setCurrentTerm(10);
        raftNode.getPersistentState().setVotedFor(8);

        assertEquals(10, raftNode.getPersistentState().getCurrentTerm());
        assertEquals(8, raftNode.getPersistentState().getVotedFor());

        raftNode.becomeFollower(12);

        assertEquals(12, raftNode.getPersistentState().getCurrentTerm());
        assertEquals(-1, raftNode.getPersistentState().getVotedFor());
        assertEquals(RaftRole.FOLLOWER, raftNode.getRaftRole());
    }

    @Test
    void testBecomeCandidate() {
        RaftNode raftNode = new RaftNode(null);
        assertEquals(RaftRole.FOLLOWER, raftNode.getRaftRole());

        raftNode.becomeCandidate();

        assertEquals(RaftRole.CANDIDATE, raftNode.getRaftRole());
    }

    @Test
    void testBecomeLeader() {
        RaftNode raftNode = new RaftNode(null);
        assertEquals(RaftRole.FOLLOWER, raftNode.getRaftRole());

        raftNode.becomeLeader();

        assertEquals(RaftRole.LEADER, raftNode.getRaftRole());
    }
}
