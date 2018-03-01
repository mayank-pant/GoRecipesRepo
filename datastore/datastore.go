package datastore

import (
	postgre "recipes/datastore/postgre"
)

func GetDatastoreClient(clientName string) (DatastoreClient, error) {
	var client DatastoreClient
	var err error
	switch clientName {
	case POSTGRE:
		client = &postgre.DatabaseClient{}
		err := client.Initialize()

		if err != nil {
			return nil, err
		}
	}

	return client, err
}
