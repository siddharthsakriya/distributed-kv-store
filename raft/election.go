package raft

func isLogUpToDate(candidateLastTerm, candidateLastIndex, myLastTerm, myLastIndex int) bool {
	if candidateLastTerm != myLastTerm {
		return candidateLastTerm > myLastTerm
	}
	return candidateLastIndex >= myLastIndex
}
