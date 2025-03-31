package triakonta

import (
	"sync"
	"time"
)

type WordProgress struct {
	PlayCount      int
	CorrectCount   int
	IncorrectCount int
	Translation    string
	LastPlayed     time.Time
}

type SessionProgress struct {
	// Key: "theme+set+segment"
	// Value: map of Greek word -> WordProgress
	Progress     map[string]map[string]*WordProgress
	PreviousRuns map[string][]map[string]*WordProgress
}

type ProgressTracker struct {
	sync.RWMutex
	Data map[string]*SessionProgress
}

func (p *ProgressTracker) GetPlayableWords(sessionId, segmentKey string, doneAfter int) (unplayed, unmastered []string) {
	p.RLock()
	defer p.RUnlock()

	session, exists := p.Data[sessionId]
	if !exists {
		return nil, nil
	}

	wordMap := session.Progress[segmentKey]
	for word, progress := range wordMap {
		if progress.CorrectCount >= doneAfter {
			continue // Skip mastered
		}

		if progress.PlayCount == 0 {
			unplayed = append(unplayed, word) // Phase 1
		} else if progress.CorrectCount == 0 {
			unmastered = append(unmastered, word) // Phase 2
		}
	}

	return unplayed, unmastered
}

func (p *ProgressTracker) RecordWordPlay(sessionId, segmentKey, greekWord, translation string) {
	p.Lock()
	defer p.Unlock()

	if _, exists := p.Data[sessionId]; !exists {
		p.Data[sessionId] = &SessionProgress{Progress: make(map[string]map[string]*WordProgress), PreviousRuns: make(map[string][]map[string]*WordProgress)}
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
		entry.Translation = translation
	}
}

func (p *ProgressTracker) RecordAnswerResult(sessionId, segmentKey, greekWord string, correct bool) {
	p.Lock()
	defer p.Unlock()

	session, exists := p.Data[sessionId]
	if !exists {
		return
	}

	segmentProgress, exists := session.Progress[segmentKey]
	if !exists {
		return
	}

	entry, exists := segmentProgress[greekWord]
	if !exists {
		return
	}

	if correct {
		entry.CorrectCount++
	} else {
		entry.IncorrectCount++
	}
}

func (p *ProgressTracker) InitWordsForSegment(sessionId, segmentKey string, greekWords []string) {
	p.Lock()
	defer p.Unlock()

	if _, exists := p.Data[sessionId]; !exists {
		p.Data[sessionId] = &SessionProgress{Progress: make(map[string]map[string]*WordProgress), PreviousRuns: make(map[string][]map[string]*WordProgress)}
	}

	if _, exists := p.Data[sessionId].Progress[segmentKey]; !exists {
		p.Data[sessionId].Progress[segmentKey] = make(map[string]*WordProgress)
	}

	for _, word := range greekWords {
		if _, exists := p.Data[sessionId].Progress[segmentKey][word]; !exists {
			p.Data[sessionId].Progress[segmentKey][word] = &WordProgress{PlayCount: 0, CorrectCount: 0, IncorrectCount: 0, Translation: ""}
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

func (p *ProgressTracker) GetRetryableWords(sessionId, segmentKey string, doneAfter int) []string {
	p.RLock()
	defer p.RUnlock()

	var retryable []string
	session, exists := p.Data[sessionId]
	if !exists {
		return retryable
	}

	wordMap := session.Progress[segmentKey]
	for word, progress := range wordMap {
		if progress.CorrectCount < doneAfter {
			retryable = append(retryable, word)
		}
	}
	return retryable
}

func (wp *WordProgress) Accuracy() float64 {
	if wp.PlayCount == 0 {
		return 0.0
	}
	return float64(wp.CorrectCount) / float64(wp.PlayCount)
}

func (p *ProgressTracker) ResetSegment(sessionId, segmentKey string) {
	p.Lock()
	defer p.Unlock()

	session, exists := p.Data[sessionId]
	if !exists {
		return
	}

	current := session.Progress[segmentKey]
	if len(current) == 0 {
		return
	}

	// Create archive map (shallow copy)
	archived := make(map[string]*WordProgress, len(current))
	for k, v := range current {
		copied := *v // deep copy to avoid mutation
		archived[k] = &copied
	}

	if session.PreviousRuns == nil {
		session.PreviousRuns = make(map[string][]map[string]*WordProgress)
	}
	session.PreviousRuns[segmentKey] = append(session.PreviousRuns[segmentKey], archived)

	// Now reset
	session.Progress[segmentKey] = make(map[string]*WordProgress)
}

func (p *ProgressTracker) ClearSegment(sessionId, segmentKey string) {
	p.Lock()
	defer p.Unlock()

	session, exists := p.Data[sessionId]
	if !exists {
		return
	}

	session.Progress[segmentKey] = make(map[string]*WordProgress)
}

func (p *ProgressTracker) GetProgressForSegment(sessionId, segmentKey string, doneAfter int) (map[string]*WordProgress, bool) {
	p.RLock()
	defer p.RUnlock()

	if session, ok := p.Data[sessionId]; ok {
		finished := p.isSegmentFinished(sessionId, segmentKey, doneAfter)
		return session.Progress[segmentKey], finished
	}
	return nil, false
}

func (p *ProgressTracker) isSegmentFinished(sessionId, segmentKey string, doneAfter int) bool {
	session, exists := p.Data[sessionId]
	if !exists {
		return false
	}

	wordMap := session.Progress[segmentKey]
	if len(wordMap) == 0 {
		return false
	}

	for _, progress := range wordMap {
		if progress.CorrectCount < doneAfter {
			return false
		}
	}

	return true
}
