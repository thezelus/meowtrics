package main

import "meowtrics/model"

func StoreEvent(event model.ClientEventData) error {

	if event.GetEventId() == "" {
		return InvalidParametersError
	}

	eventMap[event.GetEventId()] = event
	return nil
}

func RetrieveEvent(eventId string) (*model.ClientEventData, error) {

	if eventId == "" {
		return nil, InvalidParametersError
	}

	if event, ok := eventMap[eventId]; ok {
		return &event, nil
	}

	return nil, RecordNotFoundError
}

/*
func GetDbConnection() (conn *bolt.DB, error){
	filename := viper.GetString("boltDbFilePath")
	mode := 0600
	timeout := 1 * time.Second
	return bolt.Open(filename, mode, &bolt.Options{Timeout: timeout})
}


func StoreEvent(b *bolt.Bucket, event *model.ClientEventData) error {

}

func RetrieveEvent(b *bolt.Bucket, eventId *string) (*ClientEventData, error) {

}
*/
