import csv
import json


def splitColumnIntoArray(row, column):
    row[column] = list(filter(
        lambda w: len(w) > 0,
        row[column].split(', ')
    ))


dataset = []

with open('./netflix_titles.csv') as csvfile:
    for row in csv.DictReader(csvfile):
        row['show_id'] = int(row['show_id'])
        row['release_year'] = int(row['release_year'])
        splitColumnIntoArray(row, 'cast')
        splitColumnIntoArray(row, 'country')
        splitColumnIntoArray(row, 'director')
        splitColumnIntoArray(row, 'listed_in')

        dataset.append(row)

with open('./netflix_titles.json', 'w') as jsonfile:
    json.dump(dataset, jsonfile, ensure_ascii=False)
