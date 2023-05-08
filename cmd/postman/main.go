package main

import (
	"os"

	postman "github.com/rbretecher/go-postman-collection"
)

func main() {
	coll := postman.Collection{
		Info: postman.Info{
			Name: "Muerta API",
			Description: postman.Description{
				Content: "RESTful API for a term paper on 'Web application to control the shelf life of products using computer vision'.",
				Version: "1.0.0",
			},
			Version: "1.0.0",
		},
		Variables: []*postman.Variable{
			{
				Key:         "protocol",
				Name:        "Protocol",
				Value:       "http",
				Description: "Protocol for the API",
			},
			{
				Key:         "host",
				Name:        "Host",
				Value:       "localhost",
				Description: "Host for the API",
			},
			{
				Key:         "port",
				Name:        "Port",
				Value:       "8000",
				Description: "Port for the API",
			},
		},
	}
	v1Group := coll.AddItemGroup("v1")
	recipesGroup := v1Group.AddItemGroup("recipes")
	recipesGroup.AddItem(
		postman.CreateItem(postman.Item{
			Name:        "Get all recipes",
			Description: "Returns all recipes",
			Request: &postman.Request{
				URL: &postman.URL{
					Protocol: "{{protocol}}",
					Host: []string{
						"{{host}}",
					},
					Path: []string{
						"api", "v1", "recipes",
					},
					Port: "{{port}}",
					Query: []*postman.QueryParam{
						{Key: "limit", Value: "10", Description: new(string)},
						{Key: "offset", Value: "0", Description: new(string)},
						{Key: "name", Value: "", Description: new(string)},
					},
				},
				Method:      "GET",
				Description: "Returns all recipes",
			},
		}),
	)

	file, err := os.Create("muerta.postman_collection.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	err = coll.Write(file, postman.V210)
	if err != nil {
		panic(err)
	}
}
