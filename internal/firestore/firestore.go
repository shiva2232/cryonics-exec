package firestore

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

const firestoreBase = "https://firestore.googleapis.com/v1/projects/%s/databases/(default)/documents/%s"

func AddDocument(projectID, collection, idToken string, data map[string]interface{}) error {
	client := resty.New()

	url := fmt.Sprintf(firestoreBase, projectID, collection)
	resp, err := client.R().
		SetAuthToken(idToken).
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]interface{}{
			"fields": map[string]interface{}{
				"name":  map[string]string{"stringValue": data["name"].(string)},
				"email": map[string]string{"stringValue": data["email"].(string)},
			},
		}).
		Post(url)
	if err != nil {
		return err
	}
	if resp.StatusCode() >= 300 {
		return fmt.Errorf("firestore error: %s", resp.String())
	}
	fmt.Println("Firestore write success:", resp.String())
	return nil
}
