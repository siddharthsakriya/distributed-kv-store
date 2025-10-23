package unit.raft.core;

import org.junit.jupiter.api.Test;
import raft.core.LogEntry;
import raft.core.RaftLog;

import static org.junit.jupiter.api.Assertions.*;

public class RaftLogTest {
    @Test
    void testEmptyLogReturnsZeroValues() {
        RaftLog log = new RaftLog();
        assertEquals(0, log.lastIndex());
        assertEquals(0, log.lastTerm());
        assertNull(log.getEntry(1));
    }

    @Test
    void testAppendAndRetrieve() {
        RaftLog log = new RaftLog();

        log.appendLog(new LogEntry(1, 1, "put 1, 3"));
        log.appendLog(new LogEntry(2, 1, "put 1, 5"));

        assertEquals(2, log.lastIndex());
        assertEquals(1, log.lastTerm());
        assertEquals("put 1, 5", log.getEntry(2).getStateMachineCommand());
    }

    @Test
    void testAppendSequentialIndexing () {
        RaftLog log = new RaftLog();
        assertEquals(0, log.lastIndex());
        assertEquals(0, log.lastTerm());

        log.appendLog(new LogEntry(1, 1, "put 1, 3"));
        log.appendLog(new LogEntry(2, 1, "put 1, 5"));

        assertThrows(IllegalStateException.class,
                () -> log.appendLog(new LogEntry(5, 1, "put 1, 3")));
    }
}
