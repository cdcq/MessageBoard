
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

function get_message_html(id) {
    let message_html = "" +
        "    <tr><td>\n" +
        "        <div class='watch-name' id='message-name-${id}'>\n" +
        "        </div>\n" +
        "        <div class='watch-time' id='message-time-${id}'>\n" +
        "        </div>\n" +
        "    </td></tr>\n" +
        "    <tr><td>\n" +
        "        <div class='watch-content' id='message-content-${id}'>\n" +
        "        </div>\n" +
        "    </td></tr>\n"
    message_html = message_html.replaceAll("${id}", id)
    return message_html
}

function get_page_html(page_num, page_symbol, page_func) {
    let page_html = "" +
        "    <a class='watch-page-choice-a' id='watch-page-choice-${page_num}'\n" +
        "       href='javascript:void(0)' onclick='${page_func}(${page_num})'>${page_symbol}</a>\n"
    page_html = page_html.replaceAll("${page_func}", page_func)
    page_html = page_html.replaceAll("${page_num}", page_num)
    page_html = page_html.replaceAll("${page_symbol}", page_symbol)
    return page_html
}

function get_page_inner_html(page, page_cnt, page_func) {
    if (page_cnt === 0) {
        return "<a class='watch-page-choice-a'>0</a>\n"
    }
    let page_inner_html = ""
    if (page > 1) {
        page_inner_html += get_page_html(page - 1, "<", page_func)
    }
    if (page - 2 > 1) {
        page_inner_html += get_page_html(1, "1", page_func)
    }
    if (page - 2 > 2) {
        page_inner_html += "<a class='watch-page-choice-a'>...</a>\n"
    }
    let left_page = page - 2
    if (left_page <= 0) {
        left_page = 1
    }
    let right_page = page + 2
    if (right_page > page_cnt) {
        right_page = page_cnt
    }
    for (let i = left_page; i <= right_page; i++) {
        page_inner_html += get_page_html(i, i, page_func)
    }
    if (page + 2 < page_cnt - 1) {
        page_inner_html += "<a class='watch-page-choice-a'>...</a>\n"
    }
    if (page + 2 < page_cnt) {
        page_inner_html += get_page_html(page_cnt, page_cnt, page_func)
    }
    if (page < page_cnt) {
        page_inner_html += get_page_html(page + 1, ">", page_func)
    }
    return page_inner_html
}

function display_messages(page) {
    let loc = window.location
    let uri = get_uri_prefix(loc) + "/watch/get_messages"
    let ws = new WebSocket(uri)
    ws.onopen = function () {
        ws.send(page)
    }
    ws.onmessage = function (e) {
        let data = JSON.parse(e.data)
        let rows = data.Rows

        if (data.PageCnt === 0) {
            document.getElementById("pages-div").innerHTML = ""
        }
        else {
            document.getElementById("pages-div").innerHTML = get_page_inner_html(
                page, data.PageCnt, "display_messages")
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
            document.getElementById("message-name-" + rows[i].Mid).innerText =
                "「" + rows[i].Power + "」" + rows[i].Name
            document.getElementById("message-time-" + rows[i].Mid).innerText = rows[i].Time
            document.getElementById("message-content-" + rows[i].Mid).innerText = rows[i].Content
        }

    }
}