There are different branches for each step. The branches are incremental i.e. **step 2** will be having all the code of **step 1** plus code for this step.

**The movieapp contains 3 services that enable the users to rate any movie(actors, etc coming soon)**

To run the hashicorp/consul, do the following
```bash
docker run \                                                                                                                                                                 
    -d \
    -p 8500:8500 \
    -p 8600:8600/udp \
    --name=dev-consul \
    hashicorp/consul agent -server -ui \
-node=server-1 -bootstrap-expect=1 -client=0.0.0.0
```

To run the app do the following in different terminals:
```bash
go run metadata/cmd/*.go
go run rating/cmd/*.go
go run movie/cmd/*.go
```

It also gets the aggregated ratings for any record**

Services:

1. Metadata service
* API: Get metadata for a movie
* Database: Movie metadata database
* Interacts with services: None
* Data model type: Movie metadata

2. Rating service
* API: Get the aggregated rating for a record and write a rating.
* Database: Rating database.
* Interacts with services: None.
* Data model type: Rating.

3. Movie service
* API: Get the details for a movie, including the aggregated movie rating and movie metadata.
* Database: None.
* Interacts with services: Movie metadata and rating.
* Data model type: Movie details.

