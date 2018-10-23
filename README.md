# Picopub

Go Micropub implementation

Go HTTP handler, connect to datastore via Interface?

* Create
* Update
* Delete
* Media Upload

* You can create a post that's written to disk
    * Parse HTTP request -> microformat object(s)
    * Validate object, convert to stricter type
    * Write to disk
* You can upload media for that post
* The media can be stored elsewhere (ie, S3)
* You can update an existing post
* You can delete a post

* auth (maybe left as an exercise for the storage backend?)
* discovery