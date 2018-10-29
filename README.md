# IGC track viewer API - Extended

This api is an online service that will allow users to browse information about IGC files. 
IGC is an international file format for soaring track files that are used by paragliders and gliders. 
The program will store IGC files metadata in a NoSQL Database (persistent storage). The system will generate events and it will monitor for new events happening from the outside services.
These files can be browsed later by entering the specific URL paths defined later in this README.

For the development of the IGC processing, I have been using the open source IGC library for Go: <a href="https://github.com/marni/goigc">goigc</a>



## Getting Started

The project's name is **paragliding**. The root of the Igc Api is /paragliding/ which if you enter that path you get redirected to the /paragliding/api path. From that root you can write the other paths to get results.


### GET     /paragliding/api

For this get request you get meta information about the API

    {
      "uptime": <uptime>
      "info": "Service for IGC tracks."
      "version": "v1"
    }

where: <uptime> is the current uptime of the service formatted according to Duration format as specified by ISO 8601. 


### POST    /paragliding/api/track

You can do track registration with this request by putting in the request body the URL of the Igc File

    {
      "url": "<url>"
    }


As a response you get the ID for that track:

    {
      "id": "<id>"
    }



### GET     /paragliding/api/track

This request returns the array of all tracks ids as an array:
    
    [<id1>, <id2>, ...]
        


### GET     /paragliding/api/track/\<id\>

Returns the meta information about a given track with the provided <id>, or NOT FOUND response code with an empty body.
    
    {
        "H_date": <date from File Header, H-record>,
        "pilot": <pilot>,
        "glider": <glider>,
        "glider_id": <glider_id>,
        "track_length": <calculated total track length>
    }


### GET     /paragliding/api/track/\<id\>/\<field\>

Returns the single detailed meta information about a given track with the provided <id>, or NOT FOUND response code with an empty body. The response should always be a string, with the exception of the calculated track length, that should be a number.

Response will be: 

    <pilot> for pilot
    
    <glider> for glider
    
    <glider_id> for glider_id
    
    <calculated total track length> for track_length
    
    <H_date> for H_date
    


## Resources

* [Go IGC library](https://github.com/marni/goigc)
* [official MongoDB Go driver](https://github.com/mongodb/mongo-go-driver)


## Deployment

This api has been deployed in Heroku: https://igcapi.herokuapp.com/paragliding/.
It has also been deployed in OpenStack and it has these floating IPs: 10.212.137.38.


## Built With

Go Language
