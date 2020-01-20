# Start Wars Planet API

## How to run the application

1. Clone the application with `git clone https://github.com/wallacebenevides/star-wars-api.git`


2. To run the application and database with docker-compose, please use the following command -

    ```docker-compose up
    ```
> Note: By default the port number its being run on is **8080**.

## Endpoints Description

### Get All Planets

```
    URL - *http://localhost:8080/api/planets*
    Method - GET
```

### Get Planet By ID

```JSON
    URL - *http://localhost:8080/api/planets?{id}*
    Method - GET
```

### Get Planet By Name

```JSON
    URL - *http://localhost:8080/api/planets/findByName/name={name}*
    Method - GET
```

### Create Planet

```JSON
    URL - *http://localhost:8080/api/planets*
    Method - POST
    Body - (content-type = application/json)
    {
    "name": "Haruun Kal",
    "climate": "temperate",
    "terrain": "toxic cloudsea, plateaus, volcanoes",
    "films": 0
}
```

### Delete Planet

```JSON
    URL - *http://localhost:8080/api/planets*
    Method - DELETE
    Body - (content-type = application/json)
    {
    "id": "...",
}
```

## Hope everything works. Thank you.
