function openForm(bandId) {
    var all = document.getElementsByClassName('band_info_popup');
    for (var i = 0; i < all.length; i++) {
        all[i].style.display = "none"; 
    }
    document.getElementById(bandId).style.display = "block";
}

function closeForm(bandId) {
    document.getElementById(bandId).style.display = "none";
}

document.onclick = function (e) {
    if (e.target.className != '' && e.target.className !== 'members' && e.target.className !== 'bandname' && e.target.className !== 'cover') {
        var all = document.getElementsByClassName('band_info_popup');
        for (var i = 0; i < all.length; i++) {
            all[i].style.display = "none"; 
        }
    }
}