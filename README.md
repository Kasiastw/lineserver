## lineserver
 
To start the aplication, build and run from the root: 
* docker build --tag lineserver .
* docker run --publish 8080:8080 lineserver

When the application starts up, access the endpoints at the URL:
* http://localhost:8080/files/<linenumber> which serves the line using line number
* http://localhost:8080/files/the which serves the lines that contains searching phrase run curl 
