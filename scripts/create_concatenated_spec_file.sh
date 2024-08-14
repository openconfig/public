#!/bin/bash
#
# This script creates a concatenated version of all the specs, and also changes
# the format of the YAML to be keyed on the names such that duplicate model
# entries can be detected by yaml-lint.
FILES="**/*/*/.spec.yml"
stat $FILES && cat $FILES | sed 's/^- name: \(.*\)/\1:/' > concatenated_spec_file.yml
import sqlite3

def create_database():
    conn = sqlite3.connect('homes.db')
    cursor = conn.cursor()

    cursor.execute('''
    CREATE TABLE IF NOT EXISTS homes (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        price INTEGER,
        bedrooms INTEGER,
        bathrooms FLOAT,
        square_feet INTEGER,
        city TEXT,
        neighborhood TEXT,
        year_built INTEGER
    )
    ''')

    # Sample data
    homes = [
        (500000, 3, 2.5, 2000, "New York", "Brooklyn", 1990),
        (750000, 4, 3, 2500, "Los Angeles", "Hollywood", 1985),
        (300000, 2, 1, 1000, "Chicago", "Lincoln Park", 1970),
        (1000000, 5, 4.5, 3500, "San Francisco", "Nob Hill", 2000),
        (450000, 3, 2, 1800, "Boston", "Back Bay", 1960)
    ]

    cursor.executemany('''
    INSERT INTO homes (price, bedrooms, bathrooms, square_feet, city, neighborhood, year_built)
    VALUES (?, ?, ?, ?, ?, ?, ?)
    ''', homes)

    conn.commit()
    conn.close()

if __name__ == "__main__":
    create_database()
    print("Database created and populated with sample data.")
    def perform_search(params):
    conn = get_db_connection()
    cursor = conn.cursor()

    query = """
    SELECT * FROM homes
    WHERE price BETWEEN ? AND ?
    AND bedrooms >= ?
    AND bathrooms >= ?
    AND square_feet >= ?
    AND city = ?
    AND year_built >= ?
    """

    cursor.execute(query, (
        params['min_price'],
        params['max_price'],
        params['min_bedrooms'],
        params['min_bathrooms'],
        params['min_square_feet'],
        params['city'],
        params['min_year_built']
    ))

    results = cursor.fetchall()
    conn.close()

    return [dict(row) for row in results]
    <!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Home Finder</title>
    <script src="https://cdn.jsdelivr.net/npm/axios/dist/axios.min.js"></script>
</head>
<body>
    <h1>Home Finder</h1>
    <form id="search-form">
        <label for="min-price">Min Price:</label>
        <input type="number" id="min-price" name="min_price" required><br>

        <label for="max-price">Max Price:</label>
        <input type="number" id="max-price" name="max_price" required><br>

        <label for="min-bedrooms">Min Bedrooms:</label>
        <input type="number" id="min-bedrooms" name="min_bedrooms" required><br>

        <label for="min-bathrooms">Min Bathrooms:</label>
        <input type="number" id="min-bathrooms" name="min_bathrooms" step="0.5" required><br>

        <label for="min-square-feet">Min Square Feet:</label>
        <input type="number" id="min-square-feet" name="min_square_feet" required><br>

        <label for="city">City:</label>
        <input type="text" id="city" name="city" required><br>

        <label for="min-year-built">Min Year Built:</label>
        <input type="number" id="min-year-built" name="min_year_built" required><br>

        <label for="sort-by">Sort By:</label>
        <select id="sort-by" name="sort_by">
            <option value="price">Price</option>
            <option value="square_feet">Square Feet</option>
            <option value="year_built">Year Built</option>
        </select>

        <button type="submit">Search</button>
    </form>

    <div id="results"></div>

    <script>
        document.getElementById('search-form').addEventListener('submit', function(e) {
            e.preventDefault();
            const formData = new FormData(e.target);
            const searchParams = Object.fromEntries(formData);
            
            axios.post('/search', searchParams)
                .then(function (response) {
                    const resultsDiv = document.getElementById('results');
                    resultsDiv.innerHTML = '';
                    response.data.forEach(home => {
                        resultsDiv.innerHTML += `
                            <div>
                                <h2>${home.city} - ${home.neighborhood}</h2>
                                <p>Price: $${home.price}</p>
                                <p>Bedrooms: ${home.bedrooms}, Bathrooms: ${home.bathrooms}</p>
                                <p>Square Feet: ${home.square_feet}</p>
                                <p>Year Built: ${home.year_built}</p>
                            </div>
                        `;
                    });
                })
                .catch(function (error) {
                    console.error('Error:', error);
                });
        });
    </script>
</body>
</html>
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Home Finder</title>
    <script src="https://cdn.jsdelivr.net/npm/axios/dist/axios.min.js"></script>
</head>
<body>
    <h1>Home Finder</h1>
    <form id="search-form">
        <label for="min-price">Min Price:</label>
        <input type="number" id="min-price" name="min_price" required><br>

        <label for="max-price">Max Price:</label>
        <input type="number" id="max-price" name="max_price" required><br>

        <label for="min-bedrooms">Min Bedrooms:</label>
        <input type="number" id="min-bedrooms" name="min_bedrooms" required><br>

        <label for="min-bathrooms">Min Bathrooms:</label>
        <input type="number" id="min-bathrooms" name="min_bathrooms" step="0.5" required><br>

        <label for="min-square-feet">Min Square Feet:</label>
        <input type="number" id="min-square-feet" name="min_square_feet" required><br>

        <label for="city">City:</label>
        <input type="text" id="city" name="city" required><br>

        <label for="min-year-built">Min Year Built:</label>
        <input type="number" id="min-year-built" name="min_year_built" required><br>

        <label for="sort-by">Sort By:</label>
        <select id="sort-by" name="sort_by">
            <option value="price">Price</option>
            <option value="square_feet">Square Feet</option>
            <option value="year_built">Year Built</option>
        </select>

        <button type="submit">Search</button>
    </form>

    <div id="results"></div>

    <script>
        document.getElementById('search-form').addEventListener('submit', function(e) {
            e.preventDefault();
            const formData = new FormData(e.target);
            const searchParams = Object.fromEntries(formData);
            
            axios.post('/search', searchParams)
                .then(function (response) {
                    const resultsDiv = document.getElementById('results');
                    resultsDiv.innerHTML = '';
                    response.data.forEach(home => {
                        resultsDiv.innerHTML += `
                            <div>
                                <h2>${home.city} - ${home.neighborhood}</h2>
                                <p>Price: $${home.price}</p>
                                <p>Bedrooms: ${home.bedrooms}, Bathrooms: ${home.bathrooms}</p>
                                <p>Square Feet: ${home.square_feet}</p>
                                <p>Year Built: ${home.year_built}</p>
                            </div>
                        `;
                    });
                })
                .catch(function (error) {
                    console.error('Error:', error);
                });
        });
    </script>
</body>
</html>
def perform_search(params):
    conn = get_db_connection()
    cursor = conn.cursor()

    query = """
    SELECT * FROM homes
    WHERE price BETWEEN ? AND ?
    AND bedrooms >= ?
    AND bathrooms >= ?
    AND square_feet >= ?
    AND city = ?
    AND year_built >= ?
    ORDER BY {} {}
    """

    sort_column = params.get('sort_by', 'price')
    sort_order = params.get('sort_order', 'ASC')
    
    # Prevent SQL injection by validating sort_column
    allowed_columns = {'price', 'square_feet', 'year_built'}
    if sort_column not in allowed_columns:
        sort_column = 'price'
    
    query = query.format(sort_column, sort_order)

    cursor.execute(query, (
        params['min_price'],
        params['max_price'],
        params['min_bedrooms'],
        params['min_bathrooms'],
        params['min_square_feet'],
        params['city'],
        params['min_year_built']
    ))

    results = cursor.fetchall()
    conn.close()

    return [dict(row) for row in results]
    