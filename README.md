# NetflixDB

## CLI for querying a Mongo database with a Redis cache

![GitHub top language](https://img.shields.io/github/languages/top/MarioJim/NetflixDB)
![GitHub Workflow Status](https://img.shields.io/github/workflow/status/MarioJim/NetflixDB/Continuous%20Integration)

We used this [Kaggle dataset of Netflix Movies and TV Shows](https://www.kaggle.com/shivamb/netflix-shows) and transformed it into JSON with [a python script](dataset/csv_to_jsondoc.py). Then, in the first execution of our Go program we replicate the database in the MongoDB instance running in localhost.

We can then make 4 types of queries:

1. Search for a movie/actor/TV show
2. Get statistics for movies/TV shows
3. Add a new movie/TV show
4. Update a movie/TV show

Queries 3 and 4 are directly applied to the Mongo database, but queries 1 and 2 are also cached in a Redis instance (with a timeout and a maximum capacity)
