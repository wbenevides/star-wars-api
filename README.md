# Start Wars Planet API

## How to run the application

1. Clone the application with `git clone https://github.com/wallacebenevides/star-wars-api.git`


2. To run the application with Makefile, please use the following command -

    ```
        make
    ```

> Note: By default the port number its being run on is **8080**.

## Endpoints Description

### Get All Planets

```
    URL - *localhost:8080/api/planets*
    Method - GET
```

### Get Planet By ID

```JSON
    URL - *localhost:8080/api/planets/{id}*
    Method - GET
```

### Get Planet By Name

```JSON
    URL - *localhost:8080/api/planets/findByName?name={name}*
    Method - GET
```

### Create Planet

```JSON
    URL - *localhost:8080/api/planets*
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
    URL - *localhost:8080/api/planets*
    Method - DELETE
    Body - (content-type = application/json)
    {
    "id": "..."
}
```

## Hope everything works. Thank you.
