# MCManager

MCManager is an easy-to-use Minecraft server manager web application.

## Tech used:
### Go
- The Gin framework was used for the REST api server.

### SQLite
- To communicate betweeen Go and SQLite I used the Go [sql package](https://pkg.go.dev/database/sql).

### TypeScript - Next.js
- Next.js was the choice for the frontend framework.
- No other 3rd party library, package or framework was used besides Next.js.

### Features:

- Run the Minecraft server as a background process.
- Support for Minecraft servers running in Docker containers.
- Turn Minecraft server on and off.
- Check online players.
- Use Rcon.
- Create and download a backup.
- Configurable backup.
- Edit and remove files from config and world folders.
- Edit server.properties file.
- View log files.

### How to use:

Download the zip file from releases. Once unzipped, simply run the executable inside "MCManager" directory.

```
./MCManager
```

The default port is 5555, but it can be changed using the -p flag. For example run:

```
./MCManager -p=5001
```

Development:

For backend cd into the api directory and run the main.go file using the -dev flag. It'll proxy the requests to the front-end running in dev mode in port 3002.
For frontend cd into the src directory and use npm run dev, it defaults to port 3002. If you want to change the port, it needs to be changed in package.json and in the main.go file.

### TODO:
