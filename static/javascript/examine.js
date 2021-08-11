
let last_page = 1

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
        if (data.Power < 1) {
            window.location = "/watch"
        }
    }
}

function get_message_html(id) {
    let message_html = "" +
        "    <tr><td>\n" +
        "        <div class='watch-name' id='message-name-${id}'>\n" +
        "        </div>\n" +
        "        <div class='watch-time' id='message-time-${id}'>\n" +
        "        </div>\n" +
        "    </td></tr>\n" +
        "    <tr><td>\n" +
        "        <div class='examine-content' id='message-content-${id}'>\n" +
        "        </div>\n" +
        "        <button class='examine-access-button' onclick='access_message(${id})'>\n" +
        "            通过\n" +
        "        </button>\n" +
        "    </td></tr>\n"
    message_html = message_html.replaceAll("${id}", id)
    return message_html
}

function display_examine_messages(page) {
    let loc = window.location
    let uri = get_uri_prefix(loc) + "/examine/get_examine_messages"
    let ws = new WebSocket(uri)
    ws.onopen = function () {
        ws.send(page)
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
                page, data.PageCnt, "display_examine_messages")
            if (page > data.PageCnt) {
                display_examine_messages(data.PageCnt)
                return
            }
            document.getElementById("watch-page-choice-" + page).style.background = "lightgray"
        }

        if (rows === null) {
            document.getElementById("messages-table").innerHTML = "" +
                "<tr><td><div style='font-size: 50px; margin: 100px auto 0 auto;'>啥也没有</div></td></tr>"
            return
        }

        let message_inner_html = ""
        for (let i = 0; i < rows.length; i++) {
            message_inner_html += get_message_html(rows[i].Mid)
        }
        document.getElementById("messages-table").innerHTML = message_inner_html
        for (let i = 0; i < rows.length; i++) {
            document.getElementById("message-name-" + rows[i].Mid).innerText = rows[i].Name
            document.getElementById("message-time-" + rows[i].Mid).innerText = rows[i].Time
            document.getElementById("message-content-" + rows[i].Mid).innerText = rows[i].Content
        }
    }
}

function access_message(id) {
    let uri = get_uri_prefix(window.location) + "/examine/access_message"
    let ws = new WebSocket(uri)
    ws.onopen = function () {
        ws.send(id)
    }
    ws.onclose = function () {
        display_examine_messages(last_page)
    }
}