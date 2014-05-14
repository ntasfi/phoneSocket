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
			lobbyID: null,
			useWSS: false
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
			  // var recievedData = null;
			  // //the server sends two blobs to the client with different sizes. This might be PingFrames or connectFrames etc.
			  // if (e.data.size) { //just a check for it
			  //    recievedData = e.data;
			  // } else {
			  //   recievedData = JSON.parse(e.data);  
			  // }
			  
			  console.log(e.data);
			  //recieve update.
			  //check if key exists in measurments for each client.
			  	//if it does then update whatever was sent over
			  		//done
			  	//if not then add it and update with whatever was sent over

			  //trigger a message	

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

		return { //public methods and variables

			//measurments: {/*this will be filled in as events trickle in*/},

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