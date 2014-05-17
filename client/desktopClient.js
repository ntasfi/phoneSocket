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
			clientTimeoutAmount: 3000, //5 seconds
			lobbyID: null,
			useWSS: false
		};
		var clients = {};

		function createSocket() {
			var conType = "ws";
			if (settings.useWSS) {
			  conType = "wss";
			}
			socket = new WebSocket(conType+"://"+serverAddress); //connect to server
			console.log("Socket created.");

			socket.onmessage = function(e) {

			  var recievedData = null;
			  //the server sends two blobs to the client with different sizes. This might be PingFrames or connectFrames etc.
			  if (e.data.hasOwnProperty("size")) { //just a check for it
			     return //thats it
			  } else {
			    recievedData = JSON.parse(e.data);  //recieve update.
			  }

			  if (updateCount == 250) {
			  	updateCount = 0;
				document.getElementById("textUpdate").value = ""
			  } else {
			  	document.getElementById("textUpdate").value = JSON.stringify(recievedData['Measurments']) + "\n" + document.getElementById("textUpdate").value
			  }
			  
			  var key = recievedData['DeviceID']; //cleaner

			  if (!clients.hasOwnProperty(key)) { //create
          clients[key] = {};
          clients[key]["playerNumber"] = Object.keys(clients).length; //we just added one in so we are good. (created object!)			  
			  }

        clients[key]['measurments'] = recievedData['Measurments'];
        clients[key]['timestamp'] = recievedData['Timestamp'];

        //if we previously set a timer...
        if (clients[key].hasOwnProperty('clientTimeoutObj')) {
          clearTimeout(clients[key]['clientTimeoutObj']); //clear the timeout timer.  
        }

        //timer
        clients[key]['clientTimeoutObj'] = setTimeout(clientTimeout, settings.clientTimeoutAmount, key); //this doesnt work in IE.
        
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
      console.log("Player "+ clients[DeviceID]['playerNumber'] +" has timed out.")
			delete clients[DeviceID]; //delete it
		};

		return { //public methods and variables

			clients: clients,

      //returns a list of all deviceIDs in an array
      currentDeviceIDs: function() {
        var temp = [];
        for(var key in clients) {
          temp.push(key);
        }
        return temp
      },

      //returns the number of players present in the lobby
      numberOfPlayers: function() {
        return Object.keys(clients).length
      },

      //gets the measurments of a device by its playerNumber. Returns the object.
      measurmentsByPlayerNumber: function(playerNumber) {
        var temp = -1;
        for(var key in clients) {
          if (clients[key]['playerNumber'] == playerNumber) {
            temp = clients[key]['measurments'];
          } 
        }
        return temp
      },

      //gets the measurments of a device by its deviceID. Returns the object.
      measurmentsByDeviceID: function(deviceID) {
        var temp = -1;
        for(var key in clients) {
          if (key == deviceID) {
            temp = clients[key]['measurments'];
          } 
        }
        return temp
      },

			setLobbyID: function(newLobbyID) {
				settings.lobbyID = newLobbyID;
			},

			setServerAddress: function(newServerAddress) {
				serverAddress = newServerAddress;
			},

			start: function() {
				//create events?
				
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