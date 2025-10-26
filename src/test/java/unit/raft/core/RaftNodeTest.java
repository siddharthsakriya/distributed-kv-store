package unit.raft.core;

import org.junit.jupiter.api.Test;
import raft.core.RaftLog;
import raft.core.RaftNode;
import raft.rpc.RequestVoteRpcDto;
import raft.rpc.ResponseVoteRpcDto;
import raft.state.PersistentState;
import raft.state.RaftRole;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.Mockito.*;

public class RaftNodeTest {
    @Test
    void testBecomeFollower() {

        RaftNode raftNode = new RaftNode(null, new PersistentState());

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
        RaftNode raftNode = new RaftNode(null, new PersistentState());
        assertEquals(RaftRole.FOLLOWER, raftNode.getRaftRole());

        raftNode.becomeCandidate();

        assertEquals(RaftRole.CANDIDATE, raftNode.getRaftRole());
    }

    @Test
    void testBecomeLeader() {
        RaftNode raftNode = new RaftNode(null, new PersistentState());
        assertEquals(RaftRole.FOLLOWER, raftNode.getRaftRole());

        raftNode.becomeLeader();

        assertEquals(RaftRole.LEADER, raftNode.getRaftRole());
    }


    @Test
    void testHandleRequestVoteValid() {

        RaftLog mockRaftLog = mock(RaftLog.class);
        PersistentState persistentState = new PersistentState(mockRaftLog);


        RaftNode raftNode = new RaftNode(null, persistentState);
        RequestVoteRpcDto requestVoteRpcDto = new RequestVoteRpcDto(4, 2, 10, 4);

        int lastTerm = 3;
        when(mockRaftLog.lastTerm()).thenReturn(lastTerm);

        ResponseVoteRpcDto responseVoteRpcDto = raftNode.handleRequestVote(requestVoteRpcDto);

        assertEquals(2, raftNode.getPersistentState().getVotedFor());
        assertEquals(4, responseVoteRpcDto.getTerm());
        assertTrue(responseVoteRpcDto.isVoteGranted());
    }

    @Test
    void testHandleRequestVoteInvalidTerm() {
        PersistentState persistentState = new PersistentState();
        persistentState.setCurrentTerm(10);

        RaftNode raftNode = new RaftNode(null, persistentState);
        RequestVoteRpcDto requestVoteRpcDto = new RequestVoteRpcDto(4, 2, 10, 4);

        ResponseVoteRpcDto responseVoteRpcDto = raftNode.handleRequestVote(requestVoteRpcDto);

        assertEquals(-1, raftNode.getPersistentState().getVotedFor());
        assertEquals(10, responseVoteRpcDto.getTerm());
        assertFalse(responseVoteRpcDto.isVoteGranted());
    }

    @Test
    void testHandleRequestVoteAlreadyVoted() {
        PersistentState persistentState = new PersistentState();
        persistentState.setVotedFor(5);
        persistentState.setCurrentTerm(4);

        RaftNode raftNode = new RaftNode(null, persistentState);
        RequestVoteRpcDto requestVoteRpcDto = new RequestVoteRpcDto(4, 2, 10, 4);

        ResponseVoteRpcDto responseVoteRpcDto = raftNode.handleRequestVote(requestVoteRpcDto);

        assertEquals(5, raftNode.getPersistentState().getVotedFor());
        assertEquals(4, responseVoteRpcDto.getTerm());
        assertFalse(responseVoteRpcDto.isVoteGranted());
    }

    @Test
    void testHandleRequestVoteLogTermMismatch() {

        RaftLog mockRaftLog = mock(RaftLog.class);
        PersistentState persistentState = new PersistentState(mockRaftLog);
        persistentState.setCurrentTerm(6);

        RaftNode raftNode = new RaftNode(null, persistentState);
        RequestVoteRpcDto requestVoteRpcDto = new RequestVoteRpcDto(6, 2, 10, 4);

        when(mockRaftLog.lastTerm()).thenReturn(5);

        ResponseVoteRpcDto responseVoteRpcDto = raftNode.handleRequestVote(requestVoteRpcDto);

        assertEquals(-1, raftNode.getPersistentState().getVotedFor());
        assertEquals(6, responseVoteRpcDto.getTerm());
        assertFalse(responseVoteRpcDto.isVoteGranted());
    }

    @Test
    void testHandleRequestVoteLogIndexMismatch() {
        RaftLog mockRaftLog = mock(RaftLog.class);
        PersistentState persistentState = new PersistentState(mockRaftLog);
        persistentState.setCurrentTerm(6);

        RaftNode raftNode = new RaftNode(null, persistentState);
        RequestVoteRpcDto requestVoteRpcDto = new RequestVoteRpcDto(6, 2, 10, 5);

        when(mockRaftLog.lastTerm()).thenReturn(5);
        when(mockRaftLog.lastIndex()).thenReturn(11);

        ResponseVoteRpcDto responseVoteRpcDto = raftNode.handleRequestVote(requestVoteRpcDto);

        assertEquals(-1, raftNode.getPersistentState().getVotedFor());
        assertEquals(6, responseVoteRpcDto.getTerm());
        assertFalse(responseVoteRpcDto.isVoteGranted());
    }
}
