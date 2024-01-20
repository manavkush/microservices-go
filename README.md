**The movieapp contains 3 services that enable the users to rate any movie(actors, etc coming soon)**

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
