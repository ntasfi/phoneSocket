/**
	*
	* Client side for sending the device orientation etc.
	* 
	* @author Norm Tasfi <ntasfi@gmail.com>
	* @copyright Norm Tasfi <http://343hz.com>
	* @version 0.0.1
	* @license MIT License
	*
	*/

/*
todo:
  avoid using Object.prototype.method = function(){}; - DONE
  avoid global scope -> check nothing ie being hoisted into 
  better error handling -> throw new Error(message)
  refactor - DONE x1
  make it cleaner
*/

//first one is the value of the measurement and the second is our basevalue aka calibration

var deviceSocket = (function() {
  var instance;

  function init() { //private methods and variables
    
    var serverAddress = 'localhost:8080';
    var clientSendIntervalObj = null; //holds the interval timer object
    var socket = null; //connect to server
    var deviceID = guid(); //make a unique id for this device
    var failureCount = 0;
    var settings = {
      hasBeenConfigured: false,
      lobbyID: null,
      useWSS: false,
      frequency: 500,
      updateValues: {
        X: false,
        Y: false,
        Z: false,
        Alpha: false, //tilted around the z-axis
        Beta: false, //titled front-to-back
        Gamma: false //titled side-to-side
      }
    };

    var calibration = {
        X: 0,
        Y: 0,
        Z: 0,
        Alpha: 0,
        Beta: 0,
        Gamma: 0
    };

    var measurements = {
        X: null,
        Y: null,
        Z: null,
        Alpha: null,
        Beta: null,
        Gamma: null
    };

    //source: http://note19.com/2007/05/27/javascript-guid-generator/
    function guid() {
      function s4() { return Math.floor((1 + Math.random()) * 0x10000).toString(16).substring(1); }
        return s4() + s4() + '-' + s4() + '-' + s4() + '-' + s4() + '-' + s4() + s4() + s4();
    };


    /*
    * Create a socket connection to the server address.
    *
    * Returns true or false.
    */
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
        if (e.data.size) { //just a check for it
           recievedData = e.data;
        } else {
          recievedData = JSON.parse(e.data);  
        }
        

        if (recievedData['Error']) { //if we got a 
          console.log(recievedData['Error']);
          alert(recievedData['Error']);
        } else if (recievedData['ping']) {
          console.log("Recieved a ping from server.")
        } else if (recievedData['Settings']) {
          for (val in recievedData['Settings']) {
            var v = val.charAt(0).toUpperCase() + val.substring(1);
            settings.updateValues[v] = true;
          }
          settings.hasBeenConfigured = true;
          console.log("Client has been configured.");
        } else {
          console.log("Recieved unknown data...");
        }
      } //end onmessage

      socket.onopen = function(e) { //get our first message from server which tells us what to set
        console.log("WebSocket: Connected to "+ serverAddress);
        var entity = {
          LobbyID: settings.lobbyID,
          DeviceID: deviceID,
          IsDesktop: false,
        };
        socket.send(JSON.stringify(entity)); //let the server know this is our first connection.
      }
    };

    /*
    * 
    * Sends updates to the server in JSON form, based on what was set as important
    * 
    */
    function sendUpdate() {
      if (failureCount > 50) { //threshold to cut out losses
        clearInterval(clientSendIntervalObj);
        console.log("Cleared send interval. Reached max failure count.")
      }
      if (settings.hasBeenConfigured == false) {
        console.log("Error: Have not configured client.");
        failureCount++;
        return -1;
      }

      if (settings.updateValues.X == null) { //check if our device updated from the server
        console.log("Error: Update values have not been set.");
        failureCount++;
        return -1;
      }

      if (socket == null) { //has the socket been created?
        console.log("Error Socket: Socket has not been set.");
        failureCount++;
        return -1;
      }

      if (socket.readyState != WebSocket.OPEN) { //is the socket open?
        console.log("Warning Socket: Socket is not open.");
        failureCount++;
        return -1;
      }

      var temp = {};

      for (val in settings.updateValues) { //for each value in updateValues
        if (settings.updateValues[val] == true) { //if its not equal to null, aka its important
          if (measurements[val] == null) {
            temp[val] = 0;  //Mobile Chrome isnt ready as fast as Safari. So we need to just set it to zero.
          } else {
            temp[val] = measurements[val];
          }
        }
      }

      failureCount = 0;
      //encode measurements here and send.
      socket.send(JSON.stringify(temp));
    }; //end sendUpdate

    /*
    * Attemps to bind to the devices sensors via eventListeners
    *
    * Returns true if successful and false otherwise.
    */
    function bindPolling() {

      if (checkBrowserSupport() == false) {
        return false;
      }

      if (window && window.addEventListener) {
        window.addEventListener('devicemotion', function(e) {
          measurements.X = e.acceleration.x - calibration.X;
          measurements.Y = e.acceleration.y - calibration.Y;
          measurements.Z = e.acceleration.z - calibration.Z;
        }, false);
        window.addEventListener('deviceorientation', function(e) {
          measurements.Alpha = e.alpha - calibration.Alpha;
          measurements.Beta = e.beta - calibration.Beta;
          measurements.Gamma = e.gamma - calibration.Gamma;
        }, false);
        return true;
      } else {
        return false;
      }
    }; //end bindPolling


    /*
    *
    * Gets the measurments from the sensors at this exact moment and sets the calibration variables
    *
    */
    function calibrateDevice() {
        for(var i in calibration) { //for each measurement
          if (measurements[i] != null) { //if the reading isnt null
            calibration[i] = measurements[i]; //set our calibration to the new reading
          } else {
            calibration[i] = 0; //otherwise, aka our measurement reading IS null then set measurments to 0.
          }
        }
    }; //end calibrateDevice


    function checkBrowserSupport() {
      var err = doesDeviceSupport('devicemotion');
      if (err != true) {
        console.log('Device Error: Does not support DeviceMotionEvent.');
        return false;
      }
      var err = doesDeviceSupport('deviceorientation');
      if (err != true) {
        console.log('Device Error: Does not support DeviceOrientationEvent.');
        return false;
      } 
      return true;     
    };

    /*
    *
    * Checks to see if the eventName is supported by the browser.
    *
    * Returns true or false.
    */
    function doesDeviceSupport(eventName) {
      var d = {
        eName: eventName,
        obj: null,
      }

      function handleEvent(e) { //they BOTH have beta, gamma and alpha and if they are all null then its the same thing
        console.log(e);
        if (e === undefined) {
          return false;
        }
        if (e.gamma == null && e.beta == null && e.alpha == null) {
          return false;
        }
      }

      switch(d.eName) {
        case "deviceorientation":
          d.obj = window.DeviceOrientationEvent;
          break;
        case "devicemotion":
          d.obj = window.DeviceMotionEvent;
          break;
        default:
          console.log('Error: No value given for eventName.');
          return -1; //nothing given
      }

      if (d.obj) { //check for our object if it exists
        // the first ipad has a little gotcha where it says it supports it
        // and fires once but everything is null.
        window.addEventListener(d.eName, handleEvent, false);
        window.removeEventListener(d.eName, handleEvent, false);
        return true;
      } else {
        return false;
      }
    }; //end supportsDeviceOrientation

    return {
      //public methods and variables

      setLobbyID: function(newLobbyID) {
        settings.lobbyID = newLobbyID;
      },

      setServerAddress: function(newAddress) {
        serverAddress = newAddress;
      },

      setFrequency: function(newFreq) {
        settings.frequency = newFreq;
      },

      checkBrowserSupport: function() {
        return checkBrowserSupport();
      },

      measurements: function() {
        return measurements;
      },

      calibration: function() {
        return calibration;
      },

      recalibrateDevice: function() {
        calibrateDevice();
      },

      start: function() {
        console.log("Creating Socket...");
        createSocket();

        console.log("Binding polling.");
        var err = bindPolling();
        if (err != true) {
          console.log("Error: Could not bind.");
        }

        console.log("Calibrating device.");
        setTimeout(calibrateDevice, 100); //wait 100ms (safe bet) and calibrate the device. Give the events a chance to fire.
        
        console.log("Starting updates to server.");
        clientSendIntervalObj = setInterval(sendUpdate, settings.frequency);
      }
    }; //end return inside init

  }; //end init

  return {
    getInstance: function() {
      if (!instance) {
        instance = init();
      }
      return instance;
    }
  }; //end return
})();