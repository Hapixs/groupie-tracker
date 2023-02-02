function initMap() {
    let map = document.getElementById("map");
    let la = map.getAttribute("data-lat");
    let lo = map.getAttribute("data-lng");
    var myLatLng = new google.maps.LatLng(la, lo);
    var options = {     
        zoom: 8,
        center: myLatLng,
        panControl:false,
        zoomControl:false,
        mapTypeControl:false,
        scaleControl:false,
        streetViewControl:false,
        overviewMapControl:false,
        rotateControl:false,    
        mapTypeId: google.maps.MapTypeId.ROADMAP
    }
    map = new google.maps.Map(map, options);
    var marker = new google.maps.Marker({
        position: myLatLng,
        map: map,
    });
}

function ChangeMapCoordinates(la, lo) {
    map = document.getElementById("map")
    map.setAttribute("data-lat", la)
    map.setAttribute("data-lng", lo)
    console.info(la)
    console.info(lo)
    initMap()
}