# Inventory Management System with Advanced Analytics

## Project Overview
This project is a comprehensive Inventory Management System built with Go for the backend and PostgreSQL for data storage. It includes user authentication, product management, stock management, supplier management, and sales order management. 


## API Endpoints

### User Authentication
- `POST /api/register`: Register a new user.
- `POST /api/login`: Log in an existing user.

### Product Management
- `POST /products/insert`: Add a new product.
- `GET /products`: Retrieve a list of products.
- `GET /products/{id}`: Retrieve a single product by ID.
- `PUT /products/change-price`: Update a product by price.
- `DELETE /products/remove/{id}`: Delete a product by ID.

### Supplier Management
- `POST /suppliers/insert`: Add a new supplier.
- `GET /suppliers`: Retrieve a list of suppliers.
- `GET /suppliers/{id}`: Retrieve a single supplier by ID.
- `PUT /suppliers/change-email`: Update a supplier by email.
- `PUT /suppliers/change-phone`: Update a supplier by phone number.
- `DELETE /suppliers/remove/{id}`: Delete a supplier by ID.


## Getting Started

### Prerequisites
- Go 1.16 or higher
- PostgreSQL
- Docker and Docker Compose
- Python 3.8 or higher

### Installation

1. **Clone the repository**:
    ```sh
    git clone https://github.com/vennisabarfi/small_business.git
    cd small_business
    ```

2. **Set up environment variables**:
    Create a `.env` file and add your database credentials and other necessary configurations.

3. **Start the development environment**:
    ```sh
    docker-compose up --build
    ```

4. **Run migrations**:
    ```sh
    go run main.go migrate
    ```

5. **Run the server**:
    ```sh
    go run main.go
    ```
    Using fresh
      ```sh
    fresh
    ```


### Usage
- Access the API at `http://localhost:{PORT}`. Specify your port in your .env
- Use tools like Postman to interact with the endpoints.

## Contributing
- Fork the repository.
- Create a new branch.
- Make your changes.
- Submit a pull request.

### Future Integrations
- Data Analysis tools to analyze sales order and inventory data.
- AI tools for more insights.
  

## License
This project is licensed under the MIT License.



