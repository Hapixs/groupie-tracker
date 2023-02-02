if (document.cookie.indexOf("cookies") > -1) {
   // document.getElementById("defaultModal").style.display = "none"
}

// document.getElementById("search-dropdown").addEventListener("change", function() {suggest();});
document.getElementById("search-dropdown").addEventListener("keyup", function() {suggest();});

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
    // play audio and only play one at a time
    return {
        currentlyPlaying: false,
        //play and stop the audio
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
        //stop the audio
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

function suggest() {
    console.log("suggest");
    textbar = document.getElementById("search-dropdown");
    apiSearch = url+"/api/search?q="+textbar.value;
    fetch(apiSearch)
    // add the proposition to the dropdown
    .then(response => response.json())
    .then(data => {
        var options = document.getElementById("search-dropdown").options;
        for (var i = 0; i < options.length; i++) {
            options[i].remove();
        }
        for (var i = 0; i < data.length; i++) {
            var option = document.createElement("option");
            option.text = data[i].name;
            option.value = data[i].id;
            document.getElementById("search-dropdown").add(option);
        }
    })
    .catch(error => {
        console.error('Error:', error);
    }
    );
}