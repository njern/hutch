package main

type LinkScore struct {
	link  string
	score int
}

type ByScore []LinkScore

func (v ByScore) Len() int           { return len(v) }
func (v ByScore) Swap(i, j int)      { v[i], v[j] = v[j], v[i] }
func (v ByScore) Less(i, j int) bool { return v[i].score < v[j].score }
