package main

type Image struct {
	Filename string
	Filesize int64
	Content  []byte
}

type Caption struct {
	Filename     string
	Content      []byte
	NumberOfTags int
}

type DataSet struct {
	Images   map[string]Image
	Captions map[string]Caption
}
