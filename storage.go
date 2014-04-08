package naivebayes

type InMemoryStorage struct {
  freq map[string]int
  total int
}

type InMemoryItem struct {
  word string
  freq int
}

func (i *InMemoryItem) Word() string {
  return i.word
}

func (i *InMemoryItem) Freq() int {
  return i.freq
}

func NewInMemoryStorage() *InMemoryStorage {
  return &InMemoryStorage {
    freq: make(map[string]int),
    total: 0,
  }
}

func (s *InMemoryStorage) Total() int {
  return s.total
}

func (s *InMemoryStorage) Store(word string) error {
  s.freq[word]++
  s.total++
  return nil
}

func (s *InMemoryStorage) Get(word string) (int, bool) {
  v, ok := s.freq[word]
  return v, ok
}

func (s *InMemoryStorage) Iter() <-chan StorageItem {
  ch := make(chan StorageItem)
  go func () {
    for word, freq := range s.freq {
      ch <- &InMemoryItem { word, freq }
    }
    close(ch)
  }()
  return ch
}

func (s *InMemoryStorage) GetWordProbability(word string) float64 {
  value, ok := s.Get(word)
  if ! ok {
    return DefaultProbability
  }
  return float64(value) / float64(s.total)
}
