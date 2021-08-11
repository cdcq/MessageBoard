
let cnt = 0

function add_alarm(obj, message) {
    obj.innerText = message
    if (obj.style.display === "inline") {
        return 
    }
    obj.style.display = "inline"
    cnt++
    if (cnt === 1) {
        document.getElementById("register-button").disabled = true
    }
}

function remove_alarm(obj) {
    if (obj.style.display === "none" || obj.style.display === "") {
        return 
    }
    obj.style.display = "none"
    cnt--
    if (cnt === 0) {
        document.getElementById("register-button").disabled = false
    }
}

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

function check_name() {
    let name = document.getElementById("register-name").value
    let obj = document.getElementById("register-name-alarm")
    let loc = window.location
    let uri = get_uri_prefix(loc) + "/register/check-user-name"
    let ws = new WebSocket(uri)
    ws.onopen = function () {
        ws.send(name)
    }
    ws.onmessage = function (e) {
        let res = e.data
        if (res !== 'ok') {
            add_alarm(obj, res)
        }
        else {
            remove_alarm(obj)
        }
    }
}

function check_pass1() {
    check_pass2()
    let pass1 = document.getElementById("register-pass1").value
    let obj = document.getElementById("register-pass1-alarm")
    if (pass1.length === 0) {
        remove_alarm(obj)
    }
    else if (pass1.length < 6 || pass1.length > 20) {
        add_alarm(obj, "密码长度为 6 - 20 个字符")
    }
    else {
        remove_alarm(obj)
    }
}

function check_pass2() {
    let pass1 = document.getElementById("register-pass1").value
    let pass2 = document.getElementById("register-pass2").value
    let obj = document.getElementById("register-pass2-alarm")
    if (pass2.length === 0) {
        remove_alarm(obj)
    }
    else if (pass1 !== pass2) {
        add_alarm(obj, "两次输入不一致")
    }
    else {
        remove_alarm(obj)
    }
}
