package naivebayes

import (
  "log"
  "sync"
  "testing"
)

func ExmapleNaiveBayes() {
  c := NewClassifier()

  c.AddWords(`good`, []string { `楽`, `簡単`, `高給` })
  c.AddWords(`bad`, []string { `きつい`, `汚い`, `危険` })

  res, err := c.GetLogScores([]string { `楽`, `簡単`, `汚い` })
  if err != nil {
    log.Fatalf("Failed to get scores: %s", err)
  }

  log.Printf("%#v\n", res)
}

func TestGetPriorProbabilities(t *testing.T) {
  c := NewClassifier()
  c.AddWords(`good`, []string { `楽`, `簡単`, `高給` })
  c.AddWords(`bad`, []string { `きつい`, `汚い`, `危険` })
  priors, err := c.GetPriorProbabilities()
  if err != nil {
    t.Errorf("Failed to get prior probabilities: %s", err)
  }
  t.Logf("%#v", priors)
}

func TestGetLogScores(t *testing.T) {
  c := NewClassifier()
  c.AddWords(`good`, []string { `楽`, `簡単`, `高給` })
  c.AddWords(`bad`, []string { `きつい`, `汚い`, `危険` })

  res, err := c.GetLogScores([]string { `楽`, `簡単` })
  if err != nil {
    t.Errorf("Failed to GetLogScores: %s", err)
  }

  if res.MaxClasses[0] != "good" {
    t.Errorf("Expected max class to be 'good', got '%s'", res.MaxClasses[0])
  }

  res, err = c.GetLogScores([]string { `きつい`, `汚い`, `楽` })
  if res.MaxClasses[0] != "bad" {
    t.Errorf("Expected max class to be 'bad', got '%s'", res.MaxClasses[0])
  }

  res, err = c.GetLogScores([]string { `きつい`, `高給` })
  if len(res.MaxClasses) != 2 {
    t.Errorf("Expected max class to be 'good' + 'bad', got '%#v'", res)
  }
}

func TestAddFromChannel(t *testing.T) {
  c := NewClassifier()

  ch := make(chan string)
  wg := &sync.WaitGroup {}
  wg.Add(1)
  go func() {
    defer wg.Done()
    c.AddFromChannel(`good`, ch)
  }()
  list := []string { `楽`, `簡単`, `高給` }
  for _, v := range list {
    ch <- v
  }
  close(ch)

  wg.Wait()
  for _, v := range list {
    if k := c.GetWordFrequency(`good`, v); k != 1 {
      t.Errorf("Word frequency for '%s' should be 1, got %d", v, k)
    }
  }
}
