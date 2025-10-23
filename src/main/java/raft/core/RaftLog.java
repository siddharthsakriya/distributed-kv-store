package raft.core;

import java.util.ArrayList;
import java.util.List;

public class RaftLog {
    private final List<LogEntry> logEntries;

    public RaftLog() {
        this.logEntries = new ArrayList<>();
    }

    public synchronized void appendLog(LogEntry logEntry) {
        int expectedIndex = lastIndex() + 1;
        if (logEntry.getIndex() != expectedIndex) {
            throw new IllegalStateException(
              "Invalid log append: expected index " + expectedIndex + " but got " + logEntry.getIndex()
            );
        }
        logEntries.add(logEntry);
    }

    public synchronized LogEntry getEntry(int index) {
        int zeroIndex = index - 1;
        return (zeroIndex >= 0 && zeroIndex < logEntries.size()) ? logEntries.get(zeroIndex) : null;
    }

    public synchronized int lastIndex() {
        return logEntries.isEmpty() ? 0 : logEntries.getLast().getIndex();
    }

    public synchronized int lastTerm() {
        return logEntries.isEmpty() ? 0 : logEntries.getLast().getTermNumber();
    }

    public synchronized int termAt(int index) {
        int zeroIndex = index - 1;
        if (logEntries.isEmpty() || zeroIndex < 0 || zeroIndex >= logEntries.size()) return 0;
        return logEntries.get(zeroIndex).getTermNumber();
    }

    public synchronized boolean containsEntry(int prevIndex, int prevTerm) {
        if (prevIndex == 0) return true;
        int zeroIndex = prevIndex - 1;
        if (logEntries.isEmpty()|| zeroIndex < 0 || zeroIndex >= logEntries.size()) return false;
        return logEntries.get(zeroIndex).getTermNumber() == prevTerm;
    }
}

