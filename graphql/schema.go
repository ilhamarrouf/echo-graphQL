package graphql

import (
	"github.com/graphql-go/graphql"
	"github.com/ilhamarrouf/echo-graphql/db"
	"github.com/ilhamarrouf/echo-graphql/models"
	"strconv"
	"fmt"
	"log"
)

var userType = graphql.NewObject(
	graphql.ObjectConfig {
		Name: "User",
		Fields: graphql.Fields {
			"id": &graphql.Field {
				Type: graphql.String,
			},
			"name": &graphql.Field {
				Type: graphql.String,
			},
			"hobby": &graphql.Field {
				Type: graphql.String,
			},
		},
	},
)

var queryType = graphql.NewObject(
	graphql.ObjectConfig {
		Name: "Query",
		Fields: graphql.Fields {
			"User": &graphql.Field {
				Type: userType,
				Args: graphql.FieldConfigArgument {
					"id": &graphql.ArgumentConfig {
						Type: graphql.String,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					idQuery, err := strconv.ParseInt(p.Args["id"].(string), 10, 64)
					if err == nil {
						db := db.CreateConnection()
						db.SingularTable(true)
						user := models.User{}
						user.Id = idQuery
						db.First(&user)
						log.Print(idQuery)

						return &user, nil
					}

					return nil, nil
				},
			},
		},
	},
)

func ExecuteQuery(query string) *graphql.Result {
	var schema, _= graphql.NewSchema(
		graphql.SchemaConfig {
			Query: queryType,
		},
	)

	result := graphql.Do(graphql.Params {
		Schema: schema,
		RequestString: query,
	})

	if len(result.Errors) > 0 {
		fmt.Printf("Wrong result, unexpected errors: %v", result.Errors)
	}

	return result
}