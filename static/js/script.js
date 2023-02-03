if (document.cookie.indexOf("cookies") > -1) {
   // document.getElementById("defaultModal").style.display = "none"
}

document.addEventListener("DOMContentLoaded", function() {
    document.getElementById("search-dropdown").addEventListener("keyup", function() {displaySuggestions();});
});

function CloseMessage (id) {
    document.getElementById(id).style.display = "none";
}

function AcceptCookies() {
    document.cookie = "cookies=true"
    document.getElementById('defaultModal').style.display = "none"
}

function RefuseCookies() {
    // redirect to not accept cookies page
    window.location.href = "/nocookies";
}

function ChangeTheme() {
    // change tailwind dark mode
    document.documentElement.classList.toggle('dark');
    // change button text
    var button = document.getElementById("theme-button");
    if (button.innerHTML == "Dark Mode") {
        button.innerHTML = "Light Mode";
    }
    else {
        button.innerHTML = "Dark Mode";
    }
}

function playaudio() {
    return {
        currentlyPlaying: false,
        playAndStop() {
            if (this.currentlyPlaying) {
                this.$refs.audio.pause();
                this.$refs.audio.currentTime = 0;
                this.currentlyPlaying = false;
            } else {
                this.$refs.audio.play();
                this.currentlyPlaying = true;
            }
        },
        stop() {
            this.$refs.audio.pause();
            this.$refs.audio.currentTime = 0;
            this.currentlyPlaying = false;
        }
    }
}

function getAverageRGB(imgEl) {

    var blockSize = 5, // only visit every 5 pixels
        defaultRGB = {r:0,g:0,b:0}, // for non-supporting envs
        canvas = document.createElement('canvas'),
        context = canvas.getContext && canvas.getContext('2d'),
        data, width, height,
        i = -4,
        length,
        rgb = {r:0,g:0,b:0},
        count = 0;

    if (!context) {
        return defaultRGB;
    }

    height = canvas.height = imgEl.naturalHeight || imgEl.offsetHeight || imgEl.height;
    width = canvas.width = imgEl.naturalWidth || imgEl.offsetWidth || imgEl.width;

    context.drawImage(imgEl, 0, 0);

    try {
        data = context.getImageData(0, 0, width, height);
    } catch(e) {
        /* security error, img on diff domain */
        return defaultRGB;
    }

    length = data.data.length;

    while ( (i += blockSize * 4) < length ) {
        ++count;
        rgb.r += data.data[i];
        rgb.g += data.data[i+1];
        rgb.b += data.data[i+2];
    }

    // ~~ used to floor values
    rgb.r = ~~(rgb.r/count);
    rgb.g = ~~(rgb.g/count);
    rgb.b = ~~(rgb.b/count);

    return rgb;

}

function componentToHex(c) {
    var hex = c.toString(16);
    return hex.length == 1 ? "0" + hex : hex;
}
  
function rgbToHex(r, g, b) {
    return "#" + componentToHex(r) + componentToHex(g) + componentToHex(b);
}

function updateMemberBackWithColor() {
    var img =  document.getElementById("groupImage");
    var data = getAverageRGB(img);
    var hex = rgbToHex(data.r, data.g, data.b);
    document.getElementById("memberList").setAttribute('class', document.getElementById("memberList").getAttribute('class')+" bg-["+hex+"]");
}

function getGroupNameById(id) {
    // console.log("getGroupNameById");
    var apiSearch = "/api/group?id="+id;
    var groupName = "";
    fetch(apiSearch)
    .then(response => response.json())
    .then(data => {
        groupName = data.name;
    });
    return groupName;
}

function displaySuggestions() {
    // console.log("displaySuggestions");
    var textbar = document.getElementById("search-dropdown");
    
    if (textbar.value == "") {
        document.getElementById("suggestions-div").innerHTML = "";
        return;
    }
    
    var apiSearch = "/api/search?q="+textbar.value;
    var classList = "py-2 px-5 text-align-left hover:bg-gray-100 dark:hover:bg-gray-600 dark:hover:text-white border-b border-gray-200 dark:border-gray-700"
    var classTitle = "text-gray-500 dark:text-gray-400 text-sm font-bold px-5 py-2 border-b border-gray-200 dark:border-gray-700"  
    var classGroup = "py-2 px-5 text-align-right text-gray-500 dark:text-gray-400 text-sm font-bold border-b border-gray-200 dark:border-gray-700"
    
    fetch(apiSearch)
    
    .then(response => response.json())
    .then(data => {
        var suggestionsDiv = document.getElementById("suggestions-div");
        suggestionsDiv.innerHTML = "";
        
        if (data.tracks.length > 0) {
            var trackTitles = document.createElement("div");
            trackTitles.textContent = "Musiques";
            trackTitles.className = classTitle;
            suggestionsDiv.appendChild(trackTitles);
        }

        for (var i = 0; i < data.tracks.length; i++) {
            var suggestion = document.createElement("div");
            var link = document.createElement("a");
            // var group = document.createElement("p");
            suggestion.textContent = data.tracks[i].title_short;
            suggestion.className = classList;
            // group.textContent = getGroupNameById(data.tracks[i].GroupId);
            // group.className = classGroup;
            link.href = "/group/"+data.tracks[i].GroupId+"#musique";
            link.appendChild(suggestion);
            // link.appendChild(group);
            suggestionsDiv.appendChild(link);
        }

        if (data.artists.length > 0) {
            var artistTitles = document.createElement("div");
            artistTitles.textContent = "Artistes";
            artistTitles.className = classTitle;
            suggestionsDiv.appendChild(artistTitles);
        }

        for (var i = 0; i < data.artists.length; i++) {
            var suggestion = document.createElement("div");
            var group = document.createElement("p");
            var link = document.createElement("a");
            suggestion.textContent = data.artists[i].name;
            suggestion.className = classList;
            group.textContent = data.artists[i].group_name;
            group.className = classGroup;
            link.href = "/group/"+data.artists[i].groupid+"#members";
            link.appendChild(suggestion);
            link.appendChild(group);
            suggestionsDiv.appendChild(link);
        }

        if (data.groups.length > 0) {
            var groupTitles = document.createElement("div");
            groupTitles.textContent = "Groupes";
            groupTitles.className = classTitle;
            suggestionsDiv.appendChild(groupTitles);
        }

        for (var i = 0; i < data.groups.length; i++) {
            var suggestion = document.createElement("div");
            // var group = document.createElement("p");
            var link = document.createElement("a");
            suggestion.textContent = data.groups[i].Name;
            suggestion.className = classList;
            // group.textContent = getGroupNameById(data.tracks[i].GroupId);
            // group.className = classGroup;
            link.href = "/group/"+data.groups[i].Id+"/";
            link.appendChild(suggestion);
            // link.appendChild(group);
            suggestionsDiv.appendChild(link);
        }

    })
    
    .catch(error => {
        console.error('Error:', error);
    });
}
