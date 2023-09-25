package main

type Image struct {
	Filename  string
	Extension string
	Directory string
	//Filesize int64
	//Content  []byte
	Caption Caption
}

func (data *DataSet) InitImage() {
	data.Images = make(map[string]*Image)
}

type Caption struct {
	Filename  string
	Extension string
	Directory string
	//Content      []byte
	//NumberOfTags int
}

func (data *DataSet) InitCaption() {
	data.TempCaption = make(map[string]*Caption)
}

type DataSet struct {
	Images      map[string]*Image
	TempCaption map[string]*Caption
}

func InitDataSet() *DataSet {
	var newDataSet DataSet
	newDataSet.InitImage()
	newDataSet.InitCaption()
	return &newDataSet
}
