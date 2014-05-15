/**
	*
	* Client side for recieving multiple device orientations etc.
	* 
	* @author Norm Tasfi <ntasfi@gmail.com>
	* @copyright Norm Tasfi <http://343hz.com>
	* @version 0.0.1
	* @license MIT License
	*
	*/

/*
*/

//create an event that signals new update has been recieved
//write to an object with updates for each device


var lobbySocket = (function(){
	var instance;

	function init() {
		var serverAddress = 'localhost:8080';
		var updateCount = 0;
		var socket = null;
		var failureCount = 0;
		var events = {};
		var settings = {
			clientTimeoutAmount: 5000, //5 seconds
			lobbyID: null,
			useWSS: false
		};


		var clients = {
			// "device-guid-here": {
			// 	"playerNumber": 1,
			//  "timestamp": 1234567890
			// 	measurments : {
			// 	"X": 123,
			// 	"Y": 456,
			// 	"Z": 789
			// 	}
			// },
		};

		function createSocket() {
			var conType = "ws";
			if (settings.useWSS) {
			  conType = "wss";
			}
			socket = new WebSocket(conType+"://"+serverAddress); //connect to server
			console.log("Socket created.");
			socket.onmessage = function(e) {
			  			  
			  if (updateCount == 250) {
			  	updateCount = 0;
				document.getElementById("textUpdate").value = ""
			  } else {
			  	document.getElementById("textUpdate").value = e.data + "\n" + document.getElementById("textUpdate").value
			  }

			  var recievedData = null;
			  //the server sends two blobs to the client with different sizes. This might be PingFrames or connectFrames etc.
			  if (e.data.size) { //just a check for it
			     return //thats it
			  } else {
			    recievedData = JSON.parse(e.data);  //recieve update.
			  }
			  
			  var key = recievedData['DeviceID']; //cleaner

			  if (!clients.hasOwnProperty(key)) { //create
			  	clients[key] = {};
			  	clients[key]['playerNumber'] = Object.keys(clients).length; //we just added one in so we are good. (created object!)
			  }

			  //these two will always be updated
			  clients[key]['measurments'] = recievedData['Measurments'];
			  clients[key]['timestamp'] = recievedData['Timestamp'];
			  clients[key]['clientTimeout'] = setTimeout(clientTimeout(key), settings.clientTimeoutAmount); //cleans up if they dont do anything
			  

			} //end onmessage

			socket.onopen = function(e) { //get our first message from server which tells us what to set
			    console.log("WebSocket: Connected to "+ serverAddress);
			    var entity = {
			      LobbyID: settings.lobbyID,
			      DeviceID: "lobby",
			      IsDesktop: true,
			    };
			    socket.send(JSON.stringify(entity)); //let the server know this is our first connection.
			}//end onopen
		};

		function clientTimeout(DeviceID) {
			delete clients[DeviceID]; //delete it
			if (Object.keys(clients).length > 0) { //lower the player numbers
				for(var key in clients) {
					clients[key].playerNumber = clients[key].playerNumber == 1 ? 1 : clients[key].playerNumber - 1; 
				}
			}
		};

		return { //public methods and variables

			//measurments: {/*this will be filled in as events trickle in*/},

			currentPlayers: clients,

			setLobbyID: function(newLobbyID) {
				settings.lobbyID = newLobbyID;
			},

			setServerAddress: function(newServerAddress) {
				serverAddress = newServerAddress;
			},

			start: function() {
				//create events
				
				console.log("Creating Socket...");
        		createSocket();
			},
		}; //end return 
	};// end init

	return {
		getInstance: function() {
			if(!instance) {
				instance = init();
			}
			return instance;
		}
	}; //end return

})();