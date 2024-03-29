if (document.cookie.indexOf("cookies") > -1) {
   // document.getElementById("defaultModal").style.display = "none"
}

document.addEventListener("DOMContentLoaded", function() {
    document.getElementById("search-dropdown").addEventListener("keyup", function() {displaySuggestions();});
});

document.addEventListener("keydown", function(event) {
    const suggestionsDiv = document.getElementById("suggestions-div");
    const links = suggestionsDiv.getElementsByTagName("a");
    let selectedIndex = -1;
    
    for (let i = 0; i < links.length; i++) {
        if (links[i].classList.contains("selected")) {
            selectedIndex = i;
            break;
        }
    }
        
    if (event.key === "ArrowDown") {
        event.preventDefault();
        if (selectedIndex === -1 || selectedIndex === links.length - 1) {
            links[0].classList.add("selected");
        } else {
            links[selectedIndex].classList.remove("selected");
            links[selectedIndex + 1].classList.add("selected");
        }
    } else if (event.key === "ArrowUp") {
        event.preventDefault();
        if (selectedIndex === -1 || selectedIndex === 0) {
            links[links.length - 1].classList.add("selected");
        } else {
            links[selectedIndex].classList.remove("selected");
            links[selectedIndex - 1].classList.add("selected");
        }
    }
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
            stopAll();
            if (this.currentlyPlaying) {
                console.log("stop");
                this.$refs.audio.pause();
                this.$refs.audio.currentTime = 0;
                this.currentlyPlaying = false;
            } else {
                console.log("play");
                this.$refs.audio.play();
                this.currentlyPlaying = true;
            }
        },
        stop() {
            console.log("stop");
            this.$refs.audio.pause();
            this.$refs.audio.currentTime = 0;
            this.currentlyPlaying = false;
        }
    }
}

function stopAll() {
    var divAudio = document.getElementById("tracks");
    for(var i = 0, len = divAudio.length; i < len;i++) {
        if (divAudio.currentlyPlaying) {
            playaudio();
        } 
    }
    var audios = document.getElementsByTagName('audio');
    for(var i = 0, len = audios.length; i < len;i++) {
        audios[i].pause();
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


var lastfilter = ""

function displaySuggestions() {
    var textbar = document.getElementById("search-dropdown");
    
    if (textbar.value == "") {
        document.getElementById("suggestions-div").innerHTML = "";
        return;
    }

    if (textbar.value == lastfilter) {
        return;
    }

    lastfilter = textbar.value;
    
    var apiSearch = "/api/search?q="+textbar.value;
    var classList = "py-2 px-5 text-align-left hover:bg-gray-100 dark:hover:bg-gray-600 dark:hover:text-white border-b border-gray-200 dark:border-gray-700"
    var classTitle = "text-gray-500 dark:text-gray-400 text-sm font-bold px-5 py-2 border-b border-gray-200 dark:border-gray-700"  
    // var classGroup = "py-2 px-5 text-align-right text-gray-500 dark:text-gray-400 text-sm font-bold border-b border-gray-200 dark:border-gray-700"
    
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
            groupName = data.tracks[i].group_name;
            suggestion.textContent = data.tracks[i].title_short + " - " + groupName;
            suggestion.className = classList;
            link.href = "/group/"+data.tracks[i].GroupId+"#musique";
            link.appendChild(suggestion);
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
            var link = document.createElement("a");
            groupName = data.artists[i].group_name;
            suggestion.textContent = data.artists[i].name + " - " + groupName;
            suggestion.className = classList;
            link.href = "/group/"+data.artists[i].groupid+"/";
            link.appendChild(suggestion);
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
            var link = document.createElement("a");
            suggestion.textContent = data.groups[i].name;
            suggestion.className = classList;
            link.href = "/group/"+data.groups[i].id+"/";
            link.appendChild(suggestion);
            suggestionsDiv.appendChild(link);
        }

    })
    
    .catch(error => {
        console.error('Error:', error);
    });
}

function getCurrentYear() {
    var d = new Date();
    var n = d.getFullYear();
    return n;
}
  
function updateDateRangeValue() {
    document.getElementById('date-range-value').innerHTML = document.getElementById('date-range').value;
}

function updateMembersNumberValue() {
    document.getElementById('members-number-value').innerHTML = document.getElementById('members-number').value;
}

function init() {
    if (document.getElementById('date-range')) {
        updateDateRangeValue();
        document.getElementById('date-range').addEventListener('input', updateDateRangeValue);
    }
    if (document.getElementById('members-number')) {
        updateMembersNumberValue();
        document.getElementById('members-number').addEventListener('input', updateMembersNumberValue);
    }
}