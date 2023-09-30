package entities

import "sync"

type Image struct {
	Filename  string
	Extension string
	Directory string
	Filesize  int64
	Content   []byte
	Caption   Caption
}

func (image *Image) InitCaption() {
	image.Caption = Caption{}
}

type Caption struct {
	Filename     string
	Extension    string
	Directory    string
	Content      []byte
	NumberOfTags int
}

type DataSet struct {
	Images      map[string]*Image
	TempCaption map[string]*Caption
	imagesLock  sync.RWMutex
	captionLock sync.RWMutex
}

func (data *DataSet) InitImage() {
	data.Images = make(map[string]*Image)
}

func (data *DataSet) InitTempCaption() {
	data.TempCaption = make(map[string]*Caption)
}

func InitDataSet() *DataSet {
	var newDataSet DataSet
	newDataSet.InitImage()
	newDataSet.InitTempCaption()
	return &newDataSet
}
