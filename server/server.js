//Import dependencies
const path = require('path');
const http = require('http');
const express = require('express');
const socketIO = require('socket.io');
const bodyParser = require('body-parser');
const cors = require('cors');
var pretty = require('express-prettify');

// const publicPath = path.join(__dirname, '../public');
const app = express();

const corsOptions = {
    origin: ["http://localhost:4242", "http://localhost:4200", "http://localhost:8000", "https://quizz.eedama.org"], // 'http://localhost:4242'
    default: "http://localhost:8000"
}

console.log("corsOptions: ", corsOptions)

app.use(pretty({ query: 'pretty' }));

// allow OPTIONS on all resources
// app.options('*', cors(corsOptions));
app.all('*', function(req, res, next) {
    //if (req.header('origin') !== undefined ) {
   //     var origin = corsOptions.origin.indexOf(req.header('origin').toLowerCase()) > -1 ? req.headers.origin : corsOptions.default;
   // } else {
        var origin = corsOptions.default;
    //}
    res.header("Access-Control-Allow-Origin", origin);
    res.header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept");
    next();
});

const router = express.Router();

const server = http.createServer(app);
// change the server parameter for another port or domain

const io = socketIO(server);

var games = new LiveGames();
var players = new Players();

// configuramos la app para que use bodyParser(), esto nos dejara usar la informacion de los POST
app.use(bodyParser.urlencoded({ extended: true }));
app.use(bodyParser.json());


// Le decimos a la aplicación que utilize las rutas que agregamos
app.use('/api', router);
// app.use(express.static(publicPath));

//Mongodb setup
var MongoClient = require('mongodb').MongoClient;
var mongoose = require('mongoose');
var url = "mongodb://mongo:27017/";


router.get('', function(req, res) {
    res.json({ games });
});

router.get('/games/:gameId', function (req, res) {
    MongoClient.connect(url,'useUnifiedTopology: true', function(err, db) {
        if (err) throw err;
        var dbo = db.db("AskDB");
        var query = { id: parseInt(req.params["gameId"]) };
        dbo.collection('Questions').find(query).toArray(function(err, result) {
            if (err) throw err;
            db.close();
            res.json({ result });
        });
    });
    // console.log("could not fetch game");
})

router.get('/games', function (req, res) {
    MongoClient.connect(url,'useUnifiedTopology: true', function(err, db) {
        if (err) throw err;
        var dbo = db.db("AskDB");
        // var query = { id: parseInt(req.params["gameId"]) };
        dbo.collection('Questions').find("*").toArray(function(err, result) {
            if (err) throw err;
            db.close();
            res.json({ result });
        });
    });
    // console.log("could not fetch games");
})

//Starting server on port 3000
server.listen(3000, () => {
    console.log("Server started on port 3000..");
});
