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

    @Test
    void testTermAt () {
        RaftLog log = new RaftLog();

        log.appendLog(new LogEntry(1, 1, "put 1, 3"));
        log.appendLog(new LogEntry(2, 1, "put 1, 5"));
        log.appendLog(new LogEntry(3, 1, "put 1, 6"));
        log.appendLog(new LogEntry(4, 2, "put 2, 6"));

        assertEquals(1, log.termAt(3));
        assertEquals(2, log.termAt(4));
    }

    @Test
    void testTermAtDoesntExist() {
        RaftLog log = new RaftLog();

        log.appendLog(new LogEntry(1, 1, "put 1, 3"));
        log.appendLog(new LogEntry(2, 1, "put 1, 5"));
        log.appendLog(new LogEntry(3, 1, "put 1, 6"));
        log.appendLog(new LogEntry(4, 2, "put 2, 6"));

        assertEquals(0, log.termAt(5));
        assertEquals(0, log.termAt(-1));
    }

    @Test
    void testContainsEntry() {
        RaftLog log = new RaftLog();

        log.appendLog(new LogEntry(1, 1, "put 1, 3"));
        log.appendLog(new LogEntry(2, 1, "put 1, 5"));
        log.appendLog(new LogEntry(3, 1, "put 1, 6"));
        log.appendLog(new LogEntry(4, 2, "put 2, 6"));

        assertTrue(log.containsEntry(4, 2));
        assertTrue(log.containsEntry(2, 1));
    }

    @Test
    void testContainsEntryWhenFalse() {
        RaftLog log = new RaftLog();

        log.appendLog(new LogEntry(1, 1, "put 1, 3"));
        log.appendLog(new LogEntry(2, 1, "put 1, 5"));
        log.appendLog(new LogEntry(3, 1, "put 1, 6"));
        log.appendLog(new LogEntry(4, 2, "put 2, 6"));

        assertFalse(log.containsEntry(1, 2));
        assertFalse(log.containsEntry(8, 1));
    }
}
