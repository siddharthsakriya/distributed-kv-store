package raft.core;

public class LogEntry {
    private int index;
    private int termNumber;
    private String stateMachineCommand;

    public LogEntry(int index, int termNumber, String stateMachineCommand) {
        this.index = index;
        this.termNumber = termNumber;
        this.stateMachineCommand = stateMachineCommand;
    }

    public int getIndex() {
        return index;
    }

    public void setIndex(int index) {
        this.index = index;
    }

    public String getStateMachineCommand() {
        return stateMachineCommand;
    }

    public void setStateMachineCommand(String stateMachineCommand) {
        this.stateMachineCommand = stateMachineCommand;
    }

    public int getTermNumber() {
        return termNumber;
    }

    public void setTermNumber(int termNumber) {
        this.termNumber = termNumber;
    }

    @Override
    public String toString() {
        return "LogEntry{" +
                "index=" + index +
                ", termNumber=" + termNumber +
                ", stateMachineCommand='" + stateMachineCommand + '\'' +
                '}';
    }
}
