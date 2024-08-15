package index

import "go.mongodb.org/mongo-driver/mongo"

type Index struct {
	Name      string
	Model     mongo.IndexModel
	Available bool
}

type IndexList struct {
	Items []Index
}

func (il *IndexList) GetActiveIndexes() map[string]mongo.IndexModel {
	out := make(map[string]mongo.IndexModel)
	for _, item := range il.Items {
		if item.Available {
			out[item.Name] = item.Model
		}
	}
	return out
}

func (il *IndexList) GetInactiveIndexes() map[string]string {
	out := make(map[string]string)
	for _, item := range il.Items {
		if !item.Available {
			out[item.Name] = item.Name
		}
	}
	return out
}
