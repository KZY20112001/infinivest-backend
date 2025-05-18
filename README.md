# Running the backend

1. First, run the database and the redis instance via docker compose: `docker compose up -d`

2. Create .env.local file and add the configurations. An example is given here:

   ```
        PORT=8080

        POSTGRES_HOST=localhost
        POSTGRES_USER=admin
        POSTGRES_PASSWORD=admin
        POSTGRES_DB=postgres_db
        POSTGRES_PORT=5432

        REDIS_HOST=localhost
        REDIS_PORT=6379

        JWT_SECRET=4bfce2fe3346dd7eea7094a28a5e190c

        # for profile image uploading

        AWS_ACCESS_KEY='aws access key here'
        AWS_SECRET_KEY='aws secret here'

        FLASK_MICROSERVICE_URL=http://localhost:5000

        # for GoMail

        EMAIL_FROM="gmail here"
        EMAIL_PASS="gmail access code here"
   ```

3. Install packages via `go mod tidy`

4. Run the server via CompileDaemon: `CompileDaemon -build="go build -o ./build/infinivest.exe ./cmd/infinivest" -command="./build/infinivest.exe"`
