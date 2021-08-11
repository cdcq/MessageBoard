package cookies

import (
	"crypto/md5"
	"database/sql"
	"encoding/binary"
	"encoding/hex"
	"message_board/utils/queries"
	"net/http"
	"time"
)

func GetLoginCookie(uid int) string {
	unixTime := time.Now().Unix()
	var bytes1 [8] byte
	var bytes2 [4] byte
	binary.BigEndian.PutUint64(bytes1[:], uint64(unixTime))
	binary.BigEndian.PutUint32(bytes2[:], uint32(uid))
	loginCookieMd5 := md5.Sum(append(bytes1[:], bytes2[:]...))
	loginCookie := hex.EncodeToString(loginCookieMd5[:])
	return loginCookie
}

func GetLoginUid(cookieCollection []*http.Cookie, db *sql.DB) (int, error) {
	for _, cookie := range cookieCollection {
		if cookie.Name == "login" {
			loginCookie := cookie.Value
			uid, err := queries.FindLoginCookie(loginCookie, db)
			if err != nil {
				return -1, err
			}
			return uid, nil
		}
	}
	return -1, nil
}
