package raft.state;

import java.util.HashMap;
import java.util.List;

public class LeaderState {
    private final HashMap<Long, Integer> nextIndex;
    private final HashMap<Long, Integer> matchIndex;

    public LeaderState(List<Long> followerIds, int lastLogIndex) {
        nextIndex = new HashMap<>();
        matchIndex = new HashMap<>();

        for (long followerId : followerIds) {
            nextIndex.put(followerId, lastLogIndex + 1);
            matchIndex.put(followerId, 0);
        }
    }

    public void updateNextIndex(long followerId, int index) {
        nextIndex.put(followerId, index);
    }

    public void updateMatchIndex(long followerId, int index) {
        matchIndex.put(followerId, index);
    }

    public int getNextIndexById(long followerId) {
        return nextIndex.get(followerId);
    }

    public int getMatchIndexById(long followerId) {
        return matchIndex.get(followerId);
    }
}
