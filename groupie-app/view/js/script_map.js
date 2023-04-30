var cities = [[]];

function AddCity(lat, lon){
    latInt = lat.replace('\'','');
    lonInt = lon.replace('\'','');
    cities.push([latInt, lonInt]);
}   

function InitMap(){
    var map = L.map('map').setView([0, 0], 1);

    L.tileLayer('http://{s}.tile.openstreetmap.fr/osmfr/{z}/{x}/{y}.png', {
        attribution: 'données © <a href="http://osm.org/copyright">OpenStreetMap</a> contributors',
        minZoom: 1,
        maxZoom: 5
    }).addTo(map);

    for (var i = 1; i < cities.length; i++) {
        var marker = L.marker(cities[i]).addTo(map);
    }
}
