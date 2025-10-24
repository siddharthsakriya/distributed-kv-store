**Main Gist**

- Consensus algorithm for managing a replicated log

**Replicated State Machines**

- Replicated log consisting of state machine commands from clients, each state machine is deterministic.
- Consensus Algo Properties:
    - Safety: Safety under all non-byzantine conditions, including network delay, packet loses etc.
    - Available: Fully functional assuming the majority of servers are operational and can communicate with each other and clients. So 5 nodes can tolerate failure of 2 nodes at any given time.
    - Timing: Timing uses will only cause availability problems.

**Raft Algo**

1. Elect a distinguished leader, who has full responsibility of managing the the replicated log  (source of truth).
    1. Responsibilities: Accept log entries from clients, replicates them on other servers and tells servers when it is safe to **apply** an entry to their state machine.
2. If a leader goes down or fails, then a new leader is elected in an **election.**

**Raft Basics**

- At any given time a node is one of these 3 states: *leader, follower or candidate.*
    - **Followers** are passive, they ONLY respond to requests from the leader or candidate nodes.
    - **Leaders** handle all client requests (if follower is contacted, then they redirect to leader). Normally, there is only one leader server, and then the remaining servers are all followers.
    - **Candidate:** Used for electing a new leader.
- Raft divides time into **terms** of arbitrary length, and these are devised of consecutive integers. Each one begins with an election, in which one or more candidate attempt to become leader.
- Some data a node would have:
    - **Current Term:** in communication, if one nodes term is < than the other ones, it updates it to the higher term. If a **candidate or leader** detects they are stale, they revert to **follower** state.

The algo is broken down into 3 independent sub problems:

- **Leader Election**
    - **Heartbeat Mechanism:**
        - Come in the form of logs + empty logs, and if no heartbeat is sent over a certain period of time, we have an **election timeout** whereby no leader is available and an election begins to choose a new leader.
    - **Election Beginning:**
        - **Follower** increments current term and becomes **candidate**. It votes for itself and issues request to vote RPCs in parallel to all the other servers in the cluster.
    - **Election ending:**
        - Election continues until one of the 3 scenarios happens:
            - **A:** Candidate wins the election
                - Occurs when it receives a majority vote, so if there are 5 servers, it receives 3. Each server votes for at most 1 candidate in a FCFS basis. This ensures only 1 winner for a term. After winning, the new leader sends a heartbeat to establish its authority.
            - **B:** Another server establishes itself as leader
                - During election, candidate could receive a heart beat from the leader, however if the term is smaller then the RPC is discarded.
            - **C:** A period goes by with no winner
                - No one wins, happens if many followers turn into candidates at the same time. If they time out then new elections would start eventually resolving it, but we need to be careful. We have **randomized retry** for this reason.



- **Log Replication**
- **Safety**