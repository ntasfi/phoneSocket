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
  better error handling
  refactor
  make it cleaner
*/

//first one is the value of the measurement and the second is our basevalue aka calibration

var phoneSocket = function(serverAddress, frequency) {
  this.serverAddress = serverAddress || 'localhost:9999';
  this.socket = new WebSocket("ws://"+this.serverAddress); //connect to server
  this.settings = {
    frequency: frequency || 500,
    updateValues: {
      x: null,
      y: null,
      z: null,
      alpha: null, //tilted around the z-axis
      beta: null, //titled front-to-back
      gamma: null, //titled side-to-side
    }
  };
  this.measurements = {
      x: [null, 0],
      y: [null, 0],
      alpha: [null, 0],
      beta: [null, 0],
      gamma: [null, 0],
  };


  this.socket.onmessage = function(e) { //get our first message from server which tells us what to set
    console.log("WebSocket: Connected.");
    console.log('websock: '+ e.data);

    //set our values what to send over...
    var setValues = JSON.parse(e.data);
    for (val in setValues) {
      this.settings.updateValues[val] = setValues[val];
    }
  };

  var err = this.bindPolling(); //bind the polling events
  if (err != true) {
    console.log('Error: Could not bind polling.');
  }

  this.calibrateDevice(); //calibrate device

  //this bypasses the window being set as a scope
  setInterval(
    function(e){
      e.sendUpdate.call(e)
    }
  , this.frequency, this); //start sending!
};

//handles sending to the server
phoneSocket.prototype.sendUpdate = function() {
  if (this.settings.updateValues.x == null) { //check if our device updated from the server
    console.log("Error: Update values have not been set.");
    return -1;
  }

  if (this.socket == null) { //has the socket been created?
    console.log("Error Socket: Socket has not been set.");
    return -1;
  }

  if (this.socket.readyState != WebSocket.OPEN) { //is the socket open?
    console.log("Error Socket: Socket not ready.");
    return -1;
  }

  var temp = {};

  for (val in this.settings.updateValues) { //for each value in updateValues
    if (this.settings.updateValues[val] == true) { //if its not equal to null, aka its important
      temp[val] = this.measurements[0][val];
    }
  }
  //encode measurements here and send.
  this.socket.send(JSON.stringify(temp));
} //end sendUpdate

phoneSocket.prototype.bindPolling = function() {
  var err = this.doesDeviceSupport('devicemotion');
  if (err != true) {
    console.log('Device Error: Does not support DeviceMotionEvent.');
    return -1;
  }

  var err = this.doesDeviceSupport('deviceorientation');
  if (err != true) {
    console.log('Device Error: Does not support DeviceOrientationEvent.');
    return -1;
  }

  if (window && window.addEventListener) {
    window.addEventListener('devicemotion', function(e) {
      this.measurements.x[0] = e.x - this.measurements.x[1];
      this.measurements.y[0] = e.y - this.measurements.y[1];
      this.measurements.z[0] = e.z - this.measurements.z[1];
    }, false);
    window.addEventListener('deviceorientation', function(e) {
      this.measurements.alpha[0] = e.alpha - this.measurements.alpha[1];
      this.measurements.beta[0] = e.beta - this.measurements.beta[1];
      this.measurements.gamma[0] = e.gamma - this.measurements.gamma[1];
    }, false);
    return true;
  } else {
    return -1;
  }
}

//so with mobile-webkit it has alpha set to whatever the first value is read in as so we need to calibrate it!
phoneSocket.prototype.calibrateDevice = function() {
    for(var i in this.measurements) { //for each measurement
      if (this.measurements[i][0] != null) { //if the reading isnt null
        this.measurements[i][1] = this.measurements[i][0] //set our calibration to the new reading
      } else {
        this.measurements[i][1] = 0 //otherwise, aka our measurement reading IS null then set measurments to 0.
      }
    }
} //end calibrateDevice

//check if the device supports the events
phoneSocket.prototype.doesDeviceSupport = function(eventName) {
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