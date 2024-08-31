# Groupie Trackers

Groupie Trackers is a Go-based backend application that interacts with a RESTful API to fetch and manipulate data about musical artists, their concert locations, dates, and relationships. This project aims to create a user-friendly website that visualizes this data effectively.

## Objectives

The application connects to an API with four main parts:
1. **Artists**: Information about bands and artists, including their names, images, formation year, first album date, and members.
2. **Locations**: Locations of their past and/or upcoming concerts.
3. **Dates**: Dates of their past and/or upcoming concerts.
4. **Relation**: Links between artists, dates, and locations.

The goal is to build a website that displays this information using various data visualizations such as blocks, cards, tables, lists, and graphics.

## Features

- Fetch and display artist information, concert dates, and locations.
- Visualize data through different UI components.
- Handle client-server communication effectively.
- Implement features based on client-triggered actions.

## Technology Stack

- **Backend**: Go (Golang)
- **Frontend**: HTML/CSS for the user interface
- **API**: RESTful API for data retrieval

## Installation

To run this project locally, follow these steps:

1. **Clone the Repository**:

   ```
   git clone https://github.com/yourusername/groupie-trackers.git
   cd groupie-trackers
   ```

2. **Build the Application**:

    Ensure you have Go installed. Build the application using:
    ```
    go build -o groupie-trackers
    ```
3. **Run the Application**:
    Start the server:

    ```
    ./groupie-trackers
    ```
    The server will start on port 8080. Open http://localhost:8080 in your web browser to access the site.

## Endpoints
Home Page
- URL: /
- Method: GET
- Description: Displays a list of all artists with their basic information.

Artist Details
- URL: /Artist
- Method: GET
- Description: Displays detailed information about a specific artist.
- Query Parameter: id (integer) - The ID of the artist to display.

## Error Handling
- All pages will show a 404 Page Not Found error for invalid routes.
- Invalid method requests will return a 405 Method Not Allowed error.
- Internal server errors will return a 500 Internal Server Error message.

## Code Structure
- main.go: Main application file, handles routing and server setup.
- fetchData.go: Contains functions for fetching data from the API.
- handlers.go: Includes HTTP request handlers for different routes.
- templates/: Contains HTML templates for rendering data.

## *Contributing*
Contributions are welcome! Please fork the repository, make your changes, and create a pull request. Ensure your code follows the project's coding standards and passes all tests.


## *Authors*
`Ismail Bentour`


