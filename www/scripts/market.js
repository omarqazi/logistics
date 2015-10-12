var trackedUser = "a5d5ce3c-9b83-4457-9b1a-4159f793cbe2";
var websock;
var map = null;
var marker = null;
var moved = false;
$(window).load(function() {
	var webSocketUrl = "ws://localhost:8080/locations/";
	var fullUrl = webSocketUrl + trackedUser;
	websock = new WebSocket(fullUrl);
	websock.onopen= function(evt) {
		console.log("opened web socket");
		console.log(evt);
	};
	
	websock.onclose = function(evt) {
		console.log("closed web socket");
		console.log(evt);
	};
	
	websock.onmessage = function(evt) {
		var userPayload = JSON.parse(evt.data);
		var myLatLng = {lat: userPayload.Latitude, lng: userPayload.Longitude};
		
		if (marker !== null) {
			marker.setPosition(myLatLng);
		} else {
	  	  marker = new google.maps.Marker({
	  	    position: myLatLng,
	  	    map: map,
	  	    title: 'Hello World!'
	  	  });	
		}
		
		if (!moved) {
		  	map.setCenter(myLatLng);
			moved = true;	
		}
	};
	
	websock.onerror = function(evt) {
		console.log("websocket ERROR");
		console.log(evt);	
	};
});

function loadMap() {
    map = new google.maps.Map(document.getElementById('map'), {
       center: {lat: -34.397, lng: 150.644},
		zoom: 12
     });
}