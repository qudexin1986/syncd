// Copyright 2018 tinystack Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package user

import (
    "errors"
    "fmt"
    "time"

    "github.com/tinystack/syncd"
    "github.com/tinystack/goutil/gostring"
    "github.com/tinystack/goutil/goaes"
)

type Login struct {
    Name        string
    Pass        string
    token       string
    userDetail  *User
}

func (u *Login) Login() error {
    user := &User{
        Name: u.Name,
    }
    if err := user.GetByName(); err != nil {
        return err
    }
    if user.ID == 0 {
        return errors.New("user not exists")
    }
    if user.LockStatus == 0 {
        return errors.New("user is locked")
    }
    password := gostring.StrMd5(gostring.JoinStrings(u.Pass, user.Salt))
    if password !=user.Password {
        return errors.New("password incorrect")
    }

    //create token
    loginKey := gostring.StrRandom(40)
    loginRaw := fmt.Sprintf("%d\t%s", user.ID, loginKey)
    var (
        err error
        tokenBytes []byte
    )
    tokenBytes, err = goaes.Encrypt(syncd.CipherKey, []byte(loginRaw))
    if err != nil {
        return err
    }

    // u.token = gostring.Base64UrlEncode(tokenBytes)
    u.token = gostring.Base64Encode(tokenBytes)

    token := &Token{
        UserId: user.ID,
        Token: loginKey,
        ExpireTime: int(time.Now().Unix()) + 3600 * 30,
    }
    if err := token.CreateOrUpdate(); err != nil {
        return err
    }

    u.userDetail = user

    return nil
}

func (u *Login) GetToken() string {
    return u.token
}

func (u *Login) GetUserDetail() *User {
    return u.userDetail
}
