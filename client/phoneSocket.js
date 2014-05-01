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

var phoneSocket = (function() {
  var instance;

  function init() { //private methods and variables
    
    var serverAddress = 'localhost:8080';
    var socket = null; //connect to server

    var settings = {
      useWSS: false,
      frequency: 500,
      updateValues: {
        x: null,
        y: null,
        z: null,
        alpha: null, //tilted around the z-axis
        beta: null, //titled front-to-back
        gamma: null //titled side-to-side
      }
    };

    var calibration = {
        x: 0,
        y: 0,
        alpha: 0,
        beta: 0,
        gamma: 0
    };

    /*
    * Create a socket connection to the server address.
    *
    * Returns true or false.
    */
    function createSocket() {
      var conType = "ws";

      if (useWSS) {
        conType = "wss";
      }

      socket = new WebSocket(conType+"://"+serverAddress); //connect to server
      socket.onmessage = function(e) { //get our first message from server which tells us what to set
        console.log("WebSocket: Connected.");
        console.log('websock: '+ e.data);

        //set our values what to send over...
        var setValues = JSON.parse(e.data);
        for (val in setValues) {
          settings.updateValues[val] = setValues[val];
        }

      }
    };

    /*
    * 
    * Sends updates to the server in JSON form, based on what was set as important
    * 
    */
    function sendUpdate() {
      if (settings.updateValues.x == null) { //check if our device updated from the server
        console.log("Error: Update values have not been set.");
        return -1;
      }

      if (socket == null) { //has the socket been created?
        console.log("Error Socket: Socket has not been set.");
        return -1;
      }

      if (socket.readyState != WebSocket.OPEN) { //is the socket open?
        console.log("Error Socket: Socket not ready.");
        return -1;
      }

      var temp = {};

      for (val in settings.updateValues) { //for each value in updateValues
        if (settings.updateValues[val] == true) { //if its not equal to null, aka its important
          temp[val] = measurements[val];
        }
      }

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
          measurements.x = e.x - calibration.x;
          measurements.y = e.y - calibration.y;
          measurements.z = e.z - calibration.z;
        }, false);
        window.addEventListener('deviceorientation', function(e) {
          measurements.alpha = e.alpha - calibration.alpha;
          measurements.beta = e.beta - calibration.beta;
          measurements.gamma = e.gamma - calibration.gamma;
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
          if (measurements != null) { //if the reading isnt null
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

      setServerAddress: function(newAddress) {
        serverAddress = newAddress;
      },

      setFrequency: function(newFreq) {
        settings.frequency = newFreq;
      },

      isBrowserSupported: function() {
        return checkBrowserSupport();
      },

      start: function() {
        createSocket();

        var err = bindPolling();
        if (err != true) {
          console.log("Error: Could not bind.");
        }

        calibrateDevice();

        setInterval(sendUpdate, frequency);
      },

      measurements: {
        x: null,
        y: null,
        alpha: null,
        beta: null,
        gamma: null
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