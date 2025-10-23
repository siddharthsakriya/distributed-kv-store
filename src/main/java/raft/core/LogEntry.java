package raft.core;

import java.util.Objects;

public class LogEntry {
    private final int index;
    private final int termNumber;
    private final String stateMachineCommand;

    public LogEntry(int index, int termNumber, String stateMachineCommand) {
        this.index = index;
        this.termNumber = termNumber;
        this.stateMachineCommand = stateMachineCommand;
    }

    public int getIndex() {
        return index;
    }

    public String getStateMachineCommand() {
        return stateMachineCommand;
    }

    public int getTermNumber() {
        return termNumber;
    }

    @Override
    public String toString() {
        return "LogEntry{" +
                "index=" + index +
                ", termNumber=" + termNumber +
                ", stateMachineCommand='" + stateMachineCommand + '\'' +
                '}';
    }

    @Override
    public boolean equals(Object o) {
        if (!(o instanceof LogEntry logEntry)) return false;
        return index == logEntry.index && termNumber == logEntry.termNumber && Objects.equals(stateMachineCommand, logEntry.stateMachineCommand);
    }

    @Override
    public int hashCode() {
        return Objects.hash(index, termNumber, stateMachineCommand);
    }
}
