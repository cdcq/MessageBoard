
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

function check_login_status() {
    let loc = window.location
    let uri = get_uri_prefix(loc) + "/login/check-login-status"
    let ws = new WebSocket(uri)
    ws.onmessage = function (e) {
        let res = e.data
        if (res === "no") {
            return 
        }
        let data = JSON.parse(res)
        document.getElementById("login-or-user-text").innerText = data.Name
        document.getElementById("login-or-user-a").href = "/user"
        if (data.Power >= 1) {
            let exam_obj = document.getElementById("header-examine")
            exam_obj.style.display = "inline"
            let manage_obj = document.getElementById("header-management")
            manage_obj.style.display = "inline"
        }
        let path = loc.pathname
        if (path === "/login") {
            window.location = "/watch"
        }
    }
}