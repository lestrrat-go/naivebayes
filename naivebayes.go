package naivebayes

import (
  "errors"
  "math"
)

const DEFAULT_PROBABILITY = 0.00000000001


type Classifier interface {
  AddWords(string, []string)
  AddFromChannel(string, <-chan string)
  ClassStorageFor(string) Storage
  Classes() []string
}

type StorageItem interface {
  Word() string
  Freq() int
}

type Storage interface {
  Get(string) (int, bool)
  GetWordProbability(string) float64
  Iter() <-chan StorageItem
  Store(string) error
  Total() int
}

type Item struct {
  class string
  words []string
}

type NaiveBayes struct {
  storages  map[string]Storage
}

func NewClassifier() *NaiveBayes {
  c := &NaiveBayes {
    storages: make(map[string]Storage),
  }

  return c
}

func (c *NaiveBayes) ClassStorageFor(klass string) Storage {
  storage, ok := c.storages[klass]
  if ! ok {
    // use reflect to automate this
    storage = NewInMemoryStorage()
    c.storages[klass] = storage
  }
  return storage
}

func (c *NaiveBayes) AddWords(klass string, words []string) {
  klassStorage := c.ClassStorageFor(klass)
  for _, word := range words {
    klassStorage.Store(word)
  }
}

func (c *NaiveBayes) AddFromChannel(klass string, ch <-chan string) {
  klassStorage := c.ClassStorageFor(klass)
  for word := range ch {
    klassStorage.Store(word)
  }
}

func (c *NaiveBayes) GetPriorProbabilities() (map[string]float64, error) {
  // XXX cache this?
  priors := make(map[string]float64)
  sum    := 0.0
  for klass, storage := range c.storages {
    total := float64(storage.Total())
    sum   += total
    priors[klass] = total
  }

  if sum <= 0 {
    return nil, errors.New("Sum of all words are 0. No words registered?")
  }

  for klass, count := range priors {
    priors[klass] = count / sum
  }

  return priors, nil
}

type Result struct {
  Scores      map[string]float64
  MaxScore    float64
  MaxClasses  []string
}

func (c *NaiveBayes) GetLogScores(document []string) (*Result, error) {
  priors, err := c.GetPriorProbabilities()
  if err != nil {
    return nil, err
  }

  max := math.Inf(-1)
  scores := make(map[string]float64)
  for klass, storage := range c.storages {
    score := math.Log(priors[klass])
    for _, word := range document {
      score += math.Log(storage.GetWordProbability(word))
    }
    scores[klass] = score
    if max < score {
      max = score
    }
  }

  maxClasses := []string {}
  for klass, score := range scores {
    if score == max {
      maxClasses = append(maxClasses, klass)
    }
  }

  return &Result { scores, max, maxClasses }, nil
}

func (c *NaiveBayes) GetProbabilities(document []string) (*Result, error) {
  priors, err := c.GetPriorProbabilities()
  if err != nil {
    return nil, err
  }

  max := 0.0
  scores := make(map[string]float64)
  for klass, storage := range c.storages {
    score := priors[klass]
    for _, word := range document {
      p := storage.GetWordProbability(word)
      score *= p
    }
    scores[klass] = score
    if max < score {
      max = score
    }
  }

  maxClasses := []string {}
  for klass, score := range scores {
    if score == max {
      maxClasses = append(maxClasses, klass)
    }
  }

  return &Result { scores, max, maxClasses }, nil
}

func (c *NaiveBayes) GetSafeProbabilities(document []string) (*Result, error) {
  priors, err := c.GetPriorProbabilities()
  if err != nil {
    return nil, err
  }

  max := 0.0
  maxLog := math.Inf(-1)
  scores := make(map[string]float64)
  logScores := make(map[string]float64)
  for klass, storage := range c.storages {
    score := priors[klass]
    logScore := math.Log(score)
    for _, word := range document {
      p := math.Log(storage.GetWordProbability(word))
      score *= p
      logScore += math.Log(p)
    }
    scores[klass] = score
    if max < score {
      max = score
    }

    logScores[klass] = score
    if maxLog < score {
      maxLog = score
    }
  }

  maxClasses := []string {}
  for klass, score := range scores {
    if score == max {
      maxClasses = append(maxClasses, klass)
    }
  }

  maxLogClasses := []string {}
  for klass, logScore := range logScores {
    if logScore == max {
      maxLogClasses = append(maxClasses, klass)
    }
  }

  if len(maxClasses) != len(maxLogClasses) {
    return nil, errors.New("Possible underflow detected")
  }

  for k, v := range maxClasses {
    if len(maxLogClasses) >= k {
      return nil, errors.New("Possible underflow detected")
    }
    if vlog := maxLogClasses[k]; vlog != v {
      return nil, errors.New("Possible underflow detected")
    }
  }

  return &Result { scores, max, maxClasses }, nil
}

func (c *NaiveBayes) Classes() []string {
  ret := make([]string, len(c.storages))
  i   := 0
  for k, _ := range c.storages {
    ret[i] = k
    i++
  }
  return ret
}

func (c *NaiveBayes) GetWordFrequency(klass, word string) int {
  v, ok := c.ClassStorageFor(klass).Get(word)
  if !ok {
    v = 0
  }
  return v
}

func (c *NaiveBayes) GetWordFrequencies(word string) map[string]int {
  ret := make(map[string]int)
  for _, klass := range c.Classes() {
    ret[klass] = c.GetWordFrequency(klass, word)
  }
  return ret
}
