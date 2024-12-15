SETUP:

SQL Database:

Download the sql.sql file, and open it with notepad, or your code editor of choice. In the first line, where it says "YourDatabaseName", replace that with your SQL database name, then save, and run the file.

Golang Backend:

You can download the latest release file, and run that, or you could download and modify the source code yourself. 

When running main.exe, you'll be prompted for a SQL connection URL. If you don't know it, use this template:
<username>:<password>@<host>:<port>/<databasename>.

    Replace <username> with your SQL username (usually "root" for local hosting).
    Replace <password> with your SQL password (if applicable).
    Replace <host> with your host IP (use 127.0.0.1 for local hosting).
    <port> is the SQL server port (default is 3306).
    Replace <databasename> with the name of your database (the same as in the "sql.sql" file).


For users modifying the source code yourself, make sure you have a working version of golang, and MYSQL installed.
You will also need these libraries installed:
	
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	_ "github.com/go-sql-driver/mysql"


HTML, CSS, JS Frontend:

Changing the names of classes, and ids (specifically in the form element) will likely break some styling as well as javascript, I recommend you don't change them unless you know what you are doing.

If you would like to use the golang backend with the login and register pages, you will need a software to route the requests to your machine. I would recommend ngrok.

If you decide to use ngrok, after downloading, launch the software, and type "ngrok http 8080" this will start a tunnel to port 8080 on your machine.
Once you have ran the command, copy the forwarding url.

Then download the html, css, and js pages. For each of the js pages, paste your ngrok url into the "YourUrl" variable, save the script, and then the setup is done.
