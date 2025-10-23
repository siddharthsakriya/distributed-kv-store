package raft.core;

import java.util.ArrayList;
import java.util.List;

public class RaftLog {
    private final List<LogEntry> logEntries;

    public RaftLog() {
        this.logEntries = new ArrayList<>();
    }

    public void appendLog(LogEntry logEntry) {
        int expectedIndex = lastIndex() + 1;
        if (logEntry.getIndex() != expectedIndex) {
            throw new IllegalStateException(
              "Invalid log append: expected index " + expectedIndex + " but got " + logEntry.getIndex()
            );
        }
        logEntries.add(logEntry);
    }

    public LogEntry getEntry(int index) {
        int zeroIndex = index - 1;
        return (zeroIndex >= 0 && zeroIndex < logEntries.size()) ? logEntries.get(zeroIndex) : null;
    }

    public int lastIndex() {
        if (!logEntries.isEmpty()) {
            return logEntries.getLast().getIndex();
        }
        return 0;
    }

    public int lastTerm() {
        if (!logEntries.isEmpty()) {
            return logEntries.getLast().getTermNumber();
        }
        return 0;
    }
}

