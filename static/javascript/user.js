
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

function load_user_info() {
    let uri = get_uri_prefix(window.location) + "/login/check-login-status"
    let ws = new WebSocket(uri)
    ws.onmessage = function (e) {
        let res = e.data
        if (res === "no") {
            window.location = "/login"
            return 
        }
        let data = JSON.parse(res)
        let name_obj = document.getElementById("user-info-name")
        name_obj.innerText = data.Name
        let power_obj = document.getElementById("user-info-power")
        let power = data.Power
        let power_str
        if (power === 3) {
            power_str = "创始人"
        }
        else if (power === 2) {
            power_str = "超级管理员"
        }
        else if (power === 2) {
            power_str = "管理员"
        }
        else {
            power_str = "普通用户"
        }
        power_obj.innerText = power_str
    }
}

function logout() {
    let loc = window.location
    let uri = get_uri_prefix(loc) + "/user/logout"
    let ws = new WebSocket(uri)
    ws.onclose = function () {
        window.location = "/login"
    }
}