package main

import (
	"bufio"
	"context"
	"fmt"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Update : Entry point for updating a document in the database
func Update(ctx context.Context, client *mongo.Client, scanner *bufio.Scanner) error {
	maybeID := ScanStringWithPrompt("Write the ID of the title you want to edit: ", scanner)
	id, err := strconv.Atoi(maybeID)
	if err != nil {
		fmt.Println("Couldn't parse an id from your input, make sure it an 8 digit number")
		return nil
	}

	filter := bson.D{primitive.E{Key: "show_id", Value: id}}
	var searchResult bson.M
	err = GetTitlesColl(client).FindOne(ctx, filter).Decode(&searchResult)
	if err != nil && err == mongo.ErrNoDocuments {
		fmt.Println("A title with that ID wasn't found")
		return nil
	} else if err != nil {
		return err
	}

	didUpdate := false
	updatedFields := bson.D{}
	for key, value := range searchResult {
		if key == "_id" {
			continue
		}
		fmt.Printf("%s: %s\n", key, ValueToString(value))
		if ScanYesNoQuestion("Do you want to edit this field?", 'n', scanner) {
			didUpdate = true
			inpMessage := fmt.Sprintf("New value for '%s': ", key)
			switch value := value.(type) {
			case string:
				updatedFields = append(updatedFields, primitive.E{
					Key:   key,
					Value: ScanStringWithPrompt(inpMessage, scanner),
				})
			case float64:
				updatedFields = append(updatedFields, primitive.E{
					Key:   key,
					Value: ScanNumberWithPrompt(inpMessage, scanner),
				})
			case primitive.A:
				newArray := primitive.A{}
				for elem := range value {
					question := fmt.Sprintf(
						"Do you want to remove element %s?",
						ValueToString(value[elem]),
					)
					if !ScanYesNoQuestion(question, 'n', scanner) {
						newArray = append(newArray, value[elem])
					}
				}
				newArray = append(newArray, ScanStringArray(scanner)...)
				updatedFields = append(updatedFields, primitive.E{
					Key:   key,
					Value: newArray,
				})
			}
		}
	}

	if !didUpdate {
		fmt.Printf("No changes made to document %s!\n", maybeID)
		return nil
	}

	update := bson.D{primitive.E{
		Key:   "$set",
		Value: updatedFields,
	}}
	if _, err = GetTitlesColl(client).UpdateOne(ctx, filter, update); err != nil {
		return err
	}
	fmt.Printf("Correctly updated document %s!\n", maybeID)
	return nil
}
