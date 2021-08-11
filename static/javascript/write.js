
function get_uri_prefix(loc) {
    let uri
    if (loc.protocol === "https:") {
        uri = "wss:"
    }
    else {
        uri = "ws:"
    }
    return uri + "//" + loc.host
}

function goto_login() {
    let loc = window.location
    let uri = get_uri_prefix(loc) + "/login/check-login-status"
    let ws = new WebSocket(uri)
    ws.onmessage = function (e) {
        let res = e.data
        if (res === "no") {
            window.location = "/login"
        }
    }
}