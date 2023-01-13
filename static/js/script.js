if (document.cookie.indexOf("cookies") > -1) {
    document.getElementById("defaultModal").style.display = "none"
}

function CloseMessage (id) {
    document.getElementById(id).style.display = "none";
}

function AcceptCookies() {
    document.cookie = "cookies=true"
    document.getElementById('defaultModal').style.display = "none"
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
        }
    }
}
