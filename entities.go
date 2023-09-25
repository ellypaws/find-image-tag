package main

type Image struct {
	Filename  string
	Extension string
	Directory string
	//Filesize int64
	//Content  []byte
	Caption Caption
}

type Caption struct {
	Filename  string
	Extension string
	Directory string
	//Content      []byte
	//NumberOfTags int
}

type DataSet struct {
	Images      map[string]*Image
	TempCaption map[string]*Caption
}
