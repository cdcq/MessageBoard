
let last_page = 1

let search_name = ""

let login_power = 1

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

function check_login() {
    let loc = window.location
    let uri = get_uri_prefix(loc) + "/login/check-login-status"
    let ws = new WebSocket(uri)
    ws.onmessage = function (e) {
        let res = e.data
        if (res === "no") {
            window.location = "/login"
            return
        }
        data = JSON.parse(res)
        login_power = data.Power
        if (data.Power < 1) {
            window.location = "/watch"
        }
    }
}

function find_user_name() {
    search_name = document.getElementById("search-user-input").value
    display_lower_power_users(1)
}

function get_user_html(id) {
    let user_html = "" +
        "    <tr><td>\n" +
        "        <div class='user-power' id='user-power-${id}'>\n" +
        "        </div>\n" +
        "        <div class='user-name' id='user-name-${id}'>\n" +
        "        </div>\n" +
        "        <label>\n" +
        "            <button class='user-select-button' onclick='change_user_power(${id})'>\n" +
        "                修改" +
        "            </button>\n" +
        "        </label>\n" +
        "        <label>\n" +
        "            <select class='user-power-select' id='user-power-select-${id}'>\n" +
        "                <option value='普通用户'>普通用户</option>\n" +
        "                <option value='管理员'>管理员</option>\n" +
        "                <option value='超级管理员'>超级管理员</option>\n" +
        "            </select>\n" +
        "        </label>\n" +
        "    </td></tr>\n" +
        "    <tr><td><hr></td></tr>\n"
    user_html = user_html.replaceAll("${id}", id)
    return user_html
}

function change_user_power(uid) {
    let ws = new WebSocket(get_uri_prefix(window.location) + "/manage/change_user_power")
    ws.onopen = function () {
        let powerS = document.getElementById("user-power-select-" + uid).value
        let power = 0
        if (powerS === "管理员") {
            power = 1
        }
        else if (powerS === "超级管理员") {
            power = 2
        }
        else {
            power = 0
        }
        let data = {
            Uid: uid,
            Power: power
        }
        ws.send(JSON.stringify(data))
    }
    ws.onclose = function () {
        display_lower_power_users(last_page)
    }
}

function display_lower_power_users(page) {
    let ws = new WebSocket(get_uri_prefix(window.location) + "/manage/get_lower_power_users")
    ws.onopen = function () {
        let data = {
            Name: search_name,
            Page: page
        }
        ws.send(JSON.stringify(data))
    }
    ws.onmessage = function (e) {
        last_page = page
        let data = JSON.parse(e.data)
        let rows = data.Rows

        if (data.PageCnt === 0) {
            document.getElementById("pages-div").innerHTML = ""
        }
        else {
            document.getElementById("pages-div").innerHTML = get_page_inner_html(
                page, data.PageCnt, "display_lower_power_users")
            if (page > data.PageCnt) {
                display_examine_messages(data.PageCnt)
                return
            }
            document.getElementById("watch-page-choice-" + page).style.background = "lightgray"
        }

        if (rows === null) {
            document.getElementById("users-table").innerHTML = "" +
                "<tr><td><div style='font-size: 50px; margin: 100px auto 0 auto;'>啥也没有</div></td></tr>"
            return
        }

        let users_inner_html = ""
        for (let i = 0; i < rows.length; i++) {
            users_inner_html += get_user_html(rows[i].Uid)
        }
        document.getElementById("users-table").innerHTML = users_inner_html
        let power_select_html = ""
        if (login_power >= 1) {
            power_select_html = "" +
                "                <option value='普通用户'>普通用户</option>\n" +
                "                <option value='管理员'>管理员</option>\n"
        }
        if (login_power >= 2) {
            power_select_html += "" +
                "                <option value='超级管理员'>超级管理员</option>\n"
        }
        for (let i = 0; i < rows.length; i++) {
            document.getElementById("user-name-" + rows[i].Uid).innerText = rows[i].Name
            document.getElementById("user-power-" + rows[i].Uid).innerText = rows[i].Power
            document.getElementById("user-power-select-" + rows[i].Uid).innerHTML = power_select_html
        }
    }
}