package hedone

import (
	"sync"
	"time"
)

type WordProgress struct {
	PlayCount  int
	LastPlayed time.Time
}

type SessionProgress struct {
	// Key: "theme+set+segment"
	// Value: map of Greek word -> WordProgress
	Progress map[string]map[string]*WordProgress
}

type ProgressTracker struct {
	sync.RWMutex
	Data map[string]*SessionProgress
}

func (p *ProgressTracker) GetPlayableWords(sessionId, segmentKey string, maxPlays int) (unplayed, playable []string) {
	p.RLock()
	defer p.RUnlock()

	session, exists := p.Data[sessionId]
	if !exists {
		return nil, nil
	}

	wordMap := session.Progress[segmentKey]
	for word, progress := range wordMap {
		if progress.PlayCount == 0 {
			unplayed = append(unplayed, word)
		} else if progress.PlayCount < maxPlays {
			playable = append(playable, word)
		}
	}

	return unplayed, playable
}

func (p *ProgressTracker) RecordWordPlay(sessionId, segmentKey, greekWord string) {
	p.Lock()
	defer p.Unlock()

	if _, exists := p.Data[sessionId]; !exists {
		p.Data[sessionId] = &SessionProgress{Progress: make(map[string]map[string]*WordProgress)}
	}

	if _, exists := p.Data[sessionId].Progress[segmentKey]; !exists {
		p.Data[sessionId].Progress[segmentKey] = make(map[string]*WordProgress)
	}

	entry, exists := p.Data[sessionId].Progress[segmentKey][greekWord]
	if !exists {
		p.Data[sessionId].Progress[segmentKey][greekWord] = &WordProgress{PlayCount: 1, LastPlayed: time.Now()}
	} else {
		entry.PlayCount++
		entry.LastPlayed = time.Now()
	}
}

func (p *ProgressTracker) InitWordsForSegment(sessionId, segmentKey string, greekWords []string) {
	p.Lock()
	defer p.Unlock()

	if _, exists := p.Data[sessionId]; !exists {
		p.Data[sessionId] = &SessionProgress{Progress: make(map[string]map[string]*WordProgress)}
	}

	if _, exists := p.Data[sessionId].Progress[segmentKey]; !exists {
		p.Data[sessionId].Progress[segmentKey] = make(map[string]*WordProgress)
	}

	for _, word := range greekWords {
		if _, exists := p.Data[sessionId].Progress[segmentKey][word]; !exists {
			p.Data[sessionId].Progress[segmentKey][word] = &WordProgress{PlayCount: 0}
		}
	}
}

func (p *ProgressTracker) Exists(sessionId, segmentKey string) bool {
	p.RLock()
	defer p.RUnlock()

	session, ok := p.Data[sessionId]
	if !ok {
		return false
	}

	_, exists := session.Progress[segmentKey]
	return exists
}
